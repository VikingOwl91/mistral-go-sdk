package mistral

import (
	"context"
	"net/url"

	"somegit.dev/vikingowl/mistral-go-sdk/model"
)

// ListModels returns a list of available models.
func (c *Client) ListModels(ctx context.Context, params *model.ListParams) (*model.ModelList, error) {
	path := "/v1/models"
	if params != nil {
		q := url.Values{}
		if params.Provider != nil {
			q.Set("provider", *params.Provider)
		}
		if params.Model != nil {
			q.Set("model", *params.Model)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp model.ModelList
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetModel retrieves a specific model by ID.
func (c *Client) GetModel(ctx context.Context, modelID string) (*model.ModelCard, error) {
	var resp model.ModelCard
	if err := c.doJSON(ctx, "GET", "/v1/models/"+modelID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteModel deletes a fine-tuned model by ID.
func (c *Client) DeleteModel(ctx context.Context, modelID string) (*model.DeleteModelOut, error) {
	var resp model.DeleteModelOut
	if err := c.doJSON(ctx, "DELETE", "/v1/models/"+modelID, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
