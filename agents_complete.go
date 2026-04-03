package mistral

import (
	"context"

	"github.com/VikingOwl91/mistral-go-sdk/agents"
	"github.com/VikingOwl91/mistral-go-sdk/chat"
)

// AgentsComplete sends an agents completion request.
func (c *Client) AgentsComplete(ctx context.Context, req *agents.CompletionRequest) (*chat.CompletionResponse, error) {
	var resp chat.CompletionResponse
	if err := c.doJSON(ctx, "POST", "/v1/agents/completions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AgentsCompleteStream sends an agents request and returns a stream of chunks.
func (c *Client) AgentsCompleteStream(ctx context.Context, req *agents.CompletionRequest) (*Stream[chat.CompletionChunk], error) {
	req.EnableStream()
	resp, err := c.doStream(ctx, "POST", "/v1/agents/completions", req)
	if err != nil {
		return nil, err
	}
	return newStream[chat.CompletionChunk](resp.Body), nil
}
