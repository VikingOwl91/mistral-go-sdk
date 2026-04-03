package mistral

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/file"
)

// UploadFile uploads a file for use with fine-tuning, batch, or OCR.
func (c *Client) UploadFile(ctx context.Context, filename string, r io.Reader, params *file.UploadParams) (*file.File, error) {
	fields := map[string]string{}
	if params != nil {
		if params.Purpose != "" {
			fields["purpose"] = string(params.Purpose)
		}
		if params.Expiry != nil {
			fields["expiry"] = strconv.Itoa(*params.Expiry)
		}
		if params.Visibility != nil {
			fields["visibility"] = string(*params.Visibility)
		}
	}
	var resp file.File
	if err := c.doMultipart(ctx, "/v1/files", filename, r, fields, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFiles returns a list of uploaded files.
func (c *Client) ListFiles(ctx context.Context, params *file.ListParams) (*file.ListResponse, error) {
	path := "/v1/files"
	if params != nil {
		q := url.Values{}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Purpose != nil {
			q.Set("purpose", string(*params.Purpose))
		}
		if params.Search != nil {
			q.Set("search", *params.Search)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp file.ListResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFile retrieves file metadata by ID.
func (c *Client) GetFile(ctx context.Context, fileID string) (*file.File, error) {
	var resp file.File
	if err := c.doJSON(ctx, "GET", "/v1/files/"+fileID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteFile deletes a file by ID.
func (c *Client) DeleteFile(ctx context.Context, fileID string) (*file.DeleteResponse, error) {
	var resp file.DeleteResponse
	if err := c.doJSON(ctx, "DELETE", "/v1/files/"+fileID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFileContent downloads the raw content of a file.
// The caller must close the returned ReadCloser.
func (c *Client) GetFileContent(ctx context.Context, fileID string) (io.ReadCloser, error) {
	resp, err := c.do(ctx, "GET", "/v1/files/"+fileID+"/content", nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseAPIError(resp)
	}
	return resp.Body, nil
}

// GetFileURL returns a signed URL for downloading a file.
// Expiry is in hours (default 24 if 0).
func (c *Client) GetFileURL(ctx context.Context, fileID string, expiryHours int) (*file.SignedURL, error) {
	path := fmt.Sprintf("/v1/files/%s/url", fileID)
	if expiryHours > 0 {
		path += fmt.Sprintf("?expiry=%d", expiryHours)
	}
	var resp file.SignedURL
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// doRawGet performs a GET request and returns the raw response.
// Used for endpoints that return non-JSON data.
func (c *Client) doRawGet(ctx context.Context, path string) (*http.Response, error) {
	return c.do(ctx, "GET", path, nil)
}
