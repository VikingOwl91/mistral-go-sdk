package mistral

import (
	"context"

	"github.com/VikingOwl91/mistral-go-sdk/classification"
	"github.com/VikingOwl91/mistral-go-sdk/moderation"
)

// Moderate sends a text moderation request.
func (c *Client) Moderate(ctx context.Context, req *moderation.Request) (*moderation.Response, error) {
	var resp moderation.Response
	if err := c.doJSON(ctx, "POST", "/v1/moderations", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ModerateChat sends a chat moderation request.
func (c *Client) ModerateChat(ctx context.Context, req *moderation.ChatRequest) (*moderation.Response, error) {
	var resp moderation.Response
	if err := c.doJSON(ctx, "POST", "/v1/chat/moderations", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Classify sends a text classification request.
func (c *Client) Classify(ctx context.Context, req *classification.Request) (*classification.Response, error) {
	var resp classification.Response
	if err := c.doJSON(ctx, "POST", "/v1/classifications", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ClassifyChat sends a chat classification request.
func (c *Client) ClassifyChat(ctx context.Context, req *classification.ChatRequest) (*classification.Response, error) {
	var resp classification.Response
	if err := c.doJSON(ctx, "POST", "/v1/chat/classifications", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
