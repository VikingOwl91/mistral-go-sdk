package mistral

import (
	"context"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// GetWorkflowWorkerInfo returns the scheduler URL, namespace, and TLS setting
// the API expects custom workers to connect with.
//
// Most callers using managed deployments do not need this — see
// Registration.DeploymentID. It is exposed for users running custom workers.
func (c *Client) GetWorkflowWorkerInfo(ctx context.Context) (*workflow.WorkerInfo, error) {
	var resp workflow.WorkerInfo
	if err := c.doJSON(ctx, "GET", "/v1/workflows/workers/whoami", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
