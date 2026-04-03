package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/finetune"
)

// CreateFineTuningJob creates a new fine-tuning job.
func (c *Client) CreateFineTuningJob(ctx context.Context, req *finetune.JobIn) (*finetune.JobOut, error) {
	var resp finetune.JobOut
	if err := c.doJSON(ctx, "POST", "/v1/fine_tuning/jobs", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListFineTuningJobs lists fine-tuning jobs with optional filters.
func (c *Client) ListFineTuningJobs(ctx context.Context, params *finetune.ListParams) (*finetune.JobsOut, error) {
	path := "/v1/fine_tuning/jobs"
	if params != nil {
		q := url.Values{}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Model != nil {
			q.Set("model", *params.Model)
		}
		if params.CreatedAfter != nil {
			q.Set("created_after", *params.CreatedAfter)
		}
		if params.CreatedBefore != nil {
			q.Set("created_before", *params.CreatedBefore)
		}
		if params.CreatedByMe != nil {
			q.Set("created_by_me", strconv.FormatBool(*params.CreatedByMe))
		}
		if params.Status != nil {
			q.Set("status", *params.Status)
		}
		if params.Suffix != nil {
			q.Set("suffix", *params.Suffix)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp finetune.JobsOut
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetFineTuningJob retrieves a fine-tuning job by ID.
func (c *Client) GetFineTuningJob(ctx context.Context, jobID string) (*finetune.JobOut, error) {
	var resp finetune.JobOut
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/fine_tuning/jobs/%s", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelFineTuningJob requests cancellation of a fine-tuning job.
func (c *Client) CancelFineTuningJob(ctx context.Context, jobID string) (*finetune.JobOut, error) {
	var resp finetune.JobOut
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/fine_tuning/jobs/%s/cancel", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// StartFineTuningJob starts a validated fine-tuning job.
func (c *Client) StartFineTuningJob(ctx context.Context, jobID string) (*finetune.JobOut, error) {
	var resp finetune.JobOut
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/fine_tuning/jobs/%s/start", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateFineTunedModel updates a fine-tuned model's name or description.
func (c *Client) UpdateFineTunedModel(ctx context.Context, modelID string, req *finetune.UpdateFTModelIn) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, "PATCH", fmt.Sprintf("/v1/fine_tuning/models/%s", modelID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// ArchiveFineTunedModel archives a fine-tuned model.
func (c *Client) ArchiveFineTunedModel(ctx context.Context, modelID string) (*finetune.ArchiveFTModelOut, error) {
	var resp finetune.ArchiveFTModelOut
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/fine_tuning/models/%s/archive", modelID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnarchiveFineTunedModel unarchives a fine-tuned model.
func (c *Client) UnarchiveFineTunedModel(ctx context.Context, modelID string) (*finetune.UnarchiveFTModelOut, error) {
	var resp finetune.UnarchiveFTModelOut
	if err := c.doJSON(ctx, "DELETE", fmt.Sprintf("/v1/fine_tuning/models/%s/archive", modelID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
