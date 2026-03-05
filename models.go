package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/model"
)

// ListModels returns a list of available models.
func (c *Client) ListModels(ctx context.Context) (*model.ModelList, error) {
	var resp model.ModelList
	if err := c.doJSON(ctx, "GET", "/v1/models", nil, &resp); err != nil {
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
