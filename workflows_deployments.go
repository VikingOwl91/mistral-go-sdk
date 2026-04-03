package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// ListWorkflowDeployments lists workflow deployments.
func (c *Client) ListWorkflowDeployments(ctx context.Context, params *workflow.DeploymentListParams) (*workflow.DeploymentListResponse, error) {
	path := "/v1/workflows/deployments"
	if params != nil {
		q := url.Values{}
		if params.ActiveOnly != nil {
			q.Set("active_only", strconv.FormatBool(*params.ActiveOnly))
		}
		if params.WorkflowName != nil {
			q.Set("workflow_name", *params.WorkflowName)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.DeploymentListResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowDeployment retrieves a workflow deployment by ID.
func (c *Client) GetWorkflowDeployment(ctx context.Context, deploymentID string) (*workflow.Deployment, error) {
	var resp workflow.Deployment
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/deployments/%s", deploymentID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
