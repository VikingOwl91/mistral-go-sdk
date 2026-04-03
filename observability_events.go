package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

// SearchChatCompletionEvents searches for chat completion events.
func (c *Client) SearchChatCompletionEvents(ctx context.Context, req *observability.SearchEventsRequest) (*observability.SearchEventsResponse, error) {
	var resp observability.SearchEventsResponse
	if err := c.doJSON(ctx, "POST", "/v1/observability/chat-completion-events/search", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchChatCompletionEventIDs searches for chat completion event IDs.
func (c *Client) SearchChatCompletionEventIDs(ctx context.Context, req *observability.SearchEventIDsRequest) (*observability.SearchEventIDsResponse, error) {
	var resp observability.SearchEventIDsResponse
	if err := c.doJSON(ctx, "POST", "/v1/observability/chat-completion-events/search-ids", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetChatCompletionEvent retrieves a chat completion event by ID.
func (c *Client) GetChatCompletionEvent(ctx context.Context, eventID string) (*observability.ChatCompletionEvent, error) {
	var resp observability.ChatCompletionEvent
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/chat-completion-events/%s", eventID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetSimilarChatCompletionEvents retrieves events similar to a given event.
func (c *Client) GetSimilarChatCompletionEvents(ctx context.Context, eventID string, params *observability.PaginationParams) (*observability.SimilarEventsResponse, error) {
	path := fmt.Sprintf("/v1/observability/chat-completion-events/%s/similar-events", eventID)
	if params != nil {
		q := url.Values{}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp observability.SimilarEventsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// JudgeChatCompletionEvent judges a chat completion event.
func (c *Client) JudgeChatCompletionEvent(ctx context.Context, eventID string, req *observability.JudgeEventRequest) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/chat-completion-events/%s/live-judging", eventID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
