package mistral

import (
	"context"
	"fmt"
	"net/url"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
)

// GetWorkflowMetrics retrieves performance metrics for a workflow.
func (c *Client) GetWorkflowMetrics(ctx context.Context, workflowName string, params *workflow.MetricsParams) (*workflow.Metrics, error) {
	path := fmt.Sprintf("/v1/workflows/%s/metrics", workflowName)
	if params != nil {
		q := url.Values{}
		if params.StartTime != nil {
			q.Set("start_time", *params.StartTime)
		}
		if params.EndTime != nil {
			q.Set("end_time", *params.EndTime)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.Metrics
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
