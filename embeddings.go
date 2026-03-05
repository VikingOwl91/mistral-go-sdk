package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/embedding"
)

// CreateEmbeddings sends an embedding request and returns the response.
func (c *Client) CreateEmbeddings(ctx context.Context, req *embedding.Request) (*embedding.Response, error) {
	var resp embedding.Response
	if err := c.doJSON(ctx, "POST", "/v1/embeddings", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
