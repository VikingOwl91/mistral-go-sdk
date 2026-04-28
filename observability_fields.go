package mistral

import (
	"context"
	"fmt"
	"net/url"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

// GetChatCompletionFields returns the searchable field definitions and groups
// for chat-completion observability events.
func (c *Client) GetChatCompletionFields(ctx context.Context) (*observability.ChatCompletionFields, error) {
	var resp observability.ChatCompletionFields
	if err := c.doJSON(ctx, "GET", "/v1/observability/chat-completion-fields", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetChatCompletionFieldOptions returns the distinct values seen for the given
// field, filtered by the requested operator.
func (c *Client) GetChatCompletionFieldOptions(ctx context.Context, fieldName string, operator observability.FieldOperator) (*observability.ChatCompletionFieldOptions, error) {
	q := url.Values{}
	q.Set("operator", string(operator))
	path := fmt.Sprintf("/v1/observability/chat-completion-fields/%s/options?%s", url.PathEscape(fieldName), q.Encode())
	var resp observability.ChatCompletionFieldOptions
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetChatCompletionFieldOptionsCounts returns per-value event counts for the
// given field, optionally filtered by the supplied filter payload.
func (c *Client) GetChatCompletionFieldOptionsCounts(ctx context.Context, fieldName string, req *observability.FieldOptionCountsRequest) (*observability.FieldOptionCounts, error) {
	if req == nil {
		req = &observability.FieldOptionCountsRequest{}
	}
	path := fmt.Sprintf("/v1/observability/chat-completion-fields/%s/options-counts", url.PathEscape(fieldName))
	var resp observability.FieldOptionCounts
	if err := c.doJSON(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
