package mistral

import (
	"context"

	"github.com/VikingOwl91/mistral-go-sdk/chat"
)

// ChatComplete sends a chat completion request and returns the full response.
func (c *Client) ChatComplete(ctx context.Context, req *chat.CompletionRequest) (*chat.CompletionResponse, error) {
	var resp chat.CompletionResponse
	if err := c.doJSON(ctx, "POST", "/v1/chat/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ChatCompleteStream sends a chat completion request and returns a stream
// of completion chunks. The caller must call Close() on the returned stream.
func (c *Client) ChatCompleteStream(ctx context.Context, req *chat.CompletionRequest) (*Stream[chat.CompletionChunk], error) {
	req.EnableStream()
	resp, err := c.doStream(ctx, "POST", "/v1/chat/completions", req)
	if err != nil {
		return nil, err
	}
	return newStream[chat.CompletionChunk](resp.Body), nil
}
