package mistral

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// doRetry executes an HTTP request with retry logic.
// buildReq is called on each attempt to create a fresh request.
func (c *Client) doRetry(ctx context.Context, buildReq func() (*http.Request, error)) (*http.Response, error) {
	maxAttempts := 1 + c.maxRetries
	var lastErr error
	var lastResp *http.Response

	for attempt := range maxAttempts {
		if attempt > 0 {
			delay := c.backoff(attempt)
			if lastResp != nil {
				if ra := retryAfterDelay(lastResp); ra > delay {
					delay = ra
				}
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		req, err := buildReq()
		if err != nil {
			return nil, fmt.Errorf("mistral: create request: %w", err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("mistral: send request: %w", err)
			if attempt < maxAttempts-1 {
				continue
			}
			return nil, lastErr
		}

		if !shouldRetry(resp.StatusCode) || attempt >= maxAttempts-1 {
			return resp, nil
		}

		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		lastResp = resp
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return lastResp, nil
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, fmt.Errorf("mistral: buffer request body: %w", err)
		}
	}

	return c.doRetry(ctx, func() (*http.Request, error) {
		var br io.Reader
		if bodyBytes != nil {
			br = bytes.NewReader(bodyBytes)
		}
		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, br)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		if bodyBytes != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		return req, nil
	})
}

func (c *Client) doJSON(ctx context.Context, method, path string, reqBody, respBody any) error {
	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("mistral: marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("mistral: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doStream(ctx context.Context, method, path string, reqBody any) (*http.Response, error) {
	var body io.Reader
	if reqBody != nil {
		data, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("mistral: marshal request: %w", err)
		}
		body = bytes.NewReader(data)
	}
	resp, err := c.do(ctx, method, path, body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseAPIError(resp)
	}
	return resp, nil
}

func (c *Client) doMultipart(ctx context.Context, path string, filename string, file io.Reader, fields map[string]string, respBody any) error {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return fmt.Errorf("mistral: create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("mistral: copy file data: %w", err)
	}
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return fmt.Errorf("mistral: write field %s: %w", k, err)
		}
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mistral: close multipart: %w", err)
	}

	bodyBytes := buf.Bytes()
	ct := w.FormDataContentType()

	resp, err := c.doRetry(ctx, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", ct)
		return req, nil
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	if respBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("mistral: decode response: %w", err)
		}
	}
	return nil
}

func (c *Client) doMultipartStream(ctx context.Context, path string, filename string, file io.Reader, fields map[string]string) (*http.Response, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("mistral: create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("mistral: copy file data: %w", err)
	}
	for k, v := range fields {
		if err := w.WriteField(k, v); err != nil {
			return nil, fmt.Errorf("mistral: write field %s: %w", k, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("mistral: close multipart: %w", err)
	}

	bodyBytes := buf.Bytes()
	ct := w.FormDataContentType()

	resp, err := c.doRetry(ctx, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Content-Type", ct)
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseAPIError(resp)
	}
	return resp, nil
}

// backoff computes the retry delay with exponential backoff and jitter.
func (c *Client) backoff(attempt int) time.Duration {
	if c.retryDelay <= 0 {
		return 0
	}
	delay := c.retryDelay * (1 << uint(attempt-1))
	jitter := 0.5 + rand.Float64() // 0.5–1.5x
	return time.Duration(float64(delay) * jitter)
}

// shouldRetry returns true if the status code is retryable.
func shouldRetry(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests || statusCode >= 500
}

// retryAfterDelay parses the Retry-After header.
func retryAfterDelay(resp *http.Response) time.Duration {
	header := resp.Header.Get("Retry-After")
	if header == "" {
		return 0
	}
	if secs, err := strconv.Atoi(header); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(header); err == nil {
		if d := time.Until(t); d > 0 {
			return d
		}
	}
	return 0
}

func parseAPIError(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("failed to read error response: %v", err),
		}
	}
	var envelope struct {
		Detail  string `json:"detail"`
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil {
		msg := envelope.Message
		if msg == "" {
			msg = envelope.Detail
		}
		if msg == "" {
			msg = string(body)
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Type:       envelope.Type,
			Message:    msg,
			Param:      envelope.Param,
			Code:       envelope.Code,
		}
	}
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
}
