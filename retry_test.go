package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/VikingOwl91/mistral-go-sdk/chat"
)

func TestRetry_429ThenSuccess(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n <= 2 {
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ok", "object": "chat.completion",
			"model": "m", "created": 0,
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "success"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(3, 1*time.Millisecond),
	)
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content.String() != "success" {
		t.Errorf("got %q", resp.Choices[0].Message.Content.String())
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts.Load())
	}
}

func TestRetry_500ThenSuccess(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"message": "server error"})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ok", "object": "chat.completion",
			"model": "m", "created": 0,
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "ok"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(2, 1*time.Millisecond),
	)
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content.String() != "ok" {
		t.Errorf("got %q", resp.Choices[0].Message.Content.String())
	}
	if attempts.Load() != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts.Load())
	}
}

func TestRetry_NoRetryOn400(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"message": "bad request"})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(3, 1*time.Millisecond),
	)
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Errorf("expected 1 attempt (no retry on 400), got %d", attempts.Load())
	}
}

func TestRetry_ExhaustedRetries(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(2, 1*time.Millisecond),
	)
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if !IsRateLimit(err) {
		t.Errorf("expected rate limit error, got: %v", err)
	}
	if attempts.Load() != 3 {
		t.Errorf("expected 3 attempts (1 + 2 retries), got %d", attempts.Load())
	}
}

func TestRetry_NoRetryByDefault(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts.Load() != 1 {
		t.Errorf("expected 1 attempt (no retries configured), got %d", attempts.Load())
	}
}

func TestRetry_RetryAfterHeader(t *testing.T) {
	var attempts atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := attempts.Add(1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ok", "object": "chat.completion",
			"model": "m", "created": 0,
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "ok"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(1, 1*time.Millisecond),
	)
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "ok" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestRetry_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
	}))
	defer server.Close()

	client := NewClient("key",
		WithBaseURL(server.URL),
		WithRetry(10, 10*time.Second),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := client.ChatComplete(ctx, &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBackoff(t *testing.T) {
	c := &Client{retryDelay: 100 * time.Millisecond}

	for i := 1; i <= 5; i++ {
		d := c.backoff(i)
		base := 100 * time.Millisecond * (1 << uint(i-1))
		minD := time.Duration(float64(base) * 0.5)
		maxD := time.Duration(float64(base) * 1.5)
		if d < minD || d > maxD {
			t.Errorf("attempt %d: backoff %v not in [%v, %v]", i, d, minD, maxD)
		}
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{200, false},
		{400, false},
		{401, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
	}
	for _, tt := range tests {
		if got := shouldRetry(tt.code); got != tt.want {
			t.Errorf("shouldRetry(%d) = %v, want %v", tt.code, got, tt.want)
		}
	}
}
