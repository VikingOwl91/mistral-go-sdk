package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// ListWorkflowRuns lists workflow runs.
func (c *Client) ListWorkflowRuns(ctx context.Context, params *workflow.RunListParams) (*workflow.ListRunsResponse, error) {
	path := "/v1/workflows/runs"
	if params != nil {
		q := url.Values{}
		if params.WorkflowIdentifier != nil {
			q.Set("workflow_identifier", *params.WorkflowIdentifier)
		}
		if params.Search != nil {
			q.Set("search", *params.Search)
		}
		if params.Status != nil {
			q.Set("status", *params.Status)
		}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.NextPageToken != nil {
			q.Set("next_page_token", *params.NextPageToken)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.ListRunsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowRun retrieves a workflow run by ID.
func (c *Client) GetWorkflowRun(ctx context.Context, runID string) (*workflow.Run, error) {
	var resp workflow.Run
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/runs/%s", runID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowRunHistory retrieves the history of a workflow run.
func (c *Client) GetWorkflowRunHistory(ctx context.Context, runID string) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/runs/%s/history", runID), nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
