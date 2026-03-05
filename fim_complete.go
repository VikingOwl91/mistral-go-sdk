package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
	"somegit.dev/vikingowl/mistral-go-sdk/fim"
)

// FIMComplete sends a Fill-In-the-Middle completion request.
func (c *Client) FIMComplete(ctx context.Context, req *fim.CompletionRequest) (*chat.CompletionResponse, error) {
	var resp chat.CompletionResponse
	if err := c.doJSON(ctx, "POST", "/v1/fim/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// FIMCompleteStream sends a FIM request and returns a stream of chunks.
func (c *Client) FIMCompleteStream(ctx context.Context, req *fim.CompletionRequest) (*Stream[chat.CompletionChunk], error) {
	req.EnableStream()
	resp, err := c.doStream(ctx, "POST", "/v1/fim/completions", req)
	if err != nil {
		return nil, err
	}
	return newStream[chat.CompletionChunk](resp.Body), nil
}
