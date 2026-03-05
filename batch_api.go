package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"somegit.dev/vikingowl/mistral-go-sdk/batch"
)

// CreateBatchJob creates a new batch inference job.
func (c *Client) CreateBatchJob(ctx context.Context, req *batch.JobIn) (*batch.JobOut, error) {
	var resp batch.JobOut
	if err := c.doJSON(ctx, "POST", "/v1/batch/jobs", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListBatchJobs lists batch jobs with optional filters.
func (c *Client) ListBatchJobs(ctx context.Context, params *batch.ListParams) (*batch.JobsOut, error) {
	path := "/v1/batch/jobs"
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
		if params.AgentID != nil {
			q.Set("agent_id", *params.AgentID)
		}
		if params.CreatedAfter != nil {
			q.Set("created_after", *params.CreatedAfter)
		}
		if params.CreatedByMe != nil {
			q.Set("created_by_me", strconv.FormatBool(*params.CreatedByMe))
		}
		if len(params.Status) > 0 {
			q.Set("status", strings.Join(params.Status, ","))
		}
		if params.OrderBy != nil {
			q.Set("order_by", *params.OrderBy)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp batch.JobsOut
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBatchJob retrieves a batch job by ID.
func (c *Client) GetBatchJob(ctx context.Context, jobID string) (*batch.JobOut, error) {
	var resp batch.JobOut
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/batch/jobs/%s", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelBatchJob requests cancellation of a batch job.
func (c *Client) CancelBatchJob(ctx context.Context, jobID string) (*batch.JobOut, error) {
	var resp batch.JobOut
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/batch/jobs/%s/cancel", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
