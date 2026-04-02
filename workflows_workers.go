package mistral

import (
	"context"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
)

// GetWorkflowWorkerInfo retrieves information about the current worker.
func (c *Client) GetWorkflowWorkerInfo(ctx context.Context) (*workflow.WorkerInfo, error) {
	var resp workflow.WorkerInfo
	if err := c.doJSON(ctx, "GET", "/v1/workflows/workers/whoami", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
