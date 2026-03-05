package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

// ChatComplete sends a chat completion request and returns the full response.
func (c *Client) ChatComplete(ctx context.Context, req *chat.CompletionRequest) (*chat.CompletionResponse, error) {
	var resp chat.CompletionResponse
	if err := c.doJSON(ctx, "POST", "/v1/chat/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
