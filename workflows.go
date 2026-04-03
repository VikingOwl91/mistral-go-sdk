package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// ListWorkflows lists workflows.
func (c *Client) ListWorkflows(ctx context.Context, params *workflow.WorkflowListParams) (*workflow.WorkflowListResponse, error) {
	path := "/v1/workflows"
	if params != nil {
		q := url.Values{}
		if params.ActiveOnly != nil {
			q.Set("active_only", strconv.FormatBool(*params.ActiveOnly))
		}
		if params.IncludeShared != nil {
			q.Set("include_shared", strconv.FormatBool(*params.IncludeShared))
		}
		if params.AvailableInChatAssistant != nil {
			q.Set("available_in_chat_assistant", strconv.FormatBool(*params.AvailableInChatAssistant))
		}
		if params.Archived != nil {
			q.Set("archived", strconv.FormatBool(*params.Archived))
		}
		if params.Cursor != nil {
			q.Set("cursor", *params.Cursor)
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.WorkflowListResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflow retrieves a workflow by identifier.
func (c *Client) GetWorkflow(ctx context.Context, workflowIdentifier string) (*workflow.Workflow, error) {
	var resp workflow.Workflow
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/%s", workflowIdentifier), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateWorkflow updates a workflow.
func (c *Client) UpdateWorkflow(ctx context.Context, workflowIdentifier string, req *workflow.WorkflowUpdateRequest) (*workflow.Workflow, error) {
	var resp workflow.Workflow
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/workflows/%s", workflowIdentifier), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ArchiveWorkflow archives a workflow.
func (c *Client) ArchiveWorkflow(ctx context.Context, workflowIdentifier string) (*workflow.WorkflowArchiveResponse, error) {
	var resp workflow.WorkflowArchiveResponse
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/workflows/%s/archive", workflowIdentifier), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnarchiveWorkflow unarchives a workflow.
func (c *Client) UnarchiveWorkflow(ctx context.Context, workflowIdentifier string) (*workflow.WorkflowArchiveResponse, error) {
	var resp workflow.WorkflowArchiveResponse
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/workflows/%s/unarchive", workflowIdentifier), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExecuteWorkflow executes a workflow.
func (c *Client) ExecuteWorkflow(ctx context.Context, workflowIdentifier string, req *workflow.ExecutionRequest) (*workflow.ExecutionResponse, error) {
	var resp workflow.ExecutionResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/%s/execute", workflowIdentifier), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListWorkflowRegistrations lists workflow registrations.
func (c *Client) ListWorkflowRegistrations(ctx context.Context, params *workflow.RegistrationListParams) (*workflow.RegistrationListResponse, error) {
	path := "/v1/workflows/registrations"
	if params != nil {
		q := url.Values{}
		if params.WorkflowID != nil {
			q.Set("workflow_id", *params.WorkflowID)
		}
		if params.TaskQueue != nil {
			q.Set("task_queue", *params.TaskQueue)
		}
		if params.ActiveOnly != nil {
			q.Set("active_only", strconv.FormatBool(*params.ActiveOnly))
		}
		if params.IncludeShared != nil {
			q.Set("include_shared", strconv.FormatBool(*params.IncludeShared))
		}
		if params.WorkflowSearch != nil {
			q.Set("workflow_search", *params.WorkflowSearch)
		}
		if params.Archived != nil {
			q.Set("archived", strconv.FormatBool(*params.Archived))
		}
		if params.WithWorkflow != nil {
			q.Set("with_workflow", strconv.FormatBool(*params.WithWorkflow))
		}
		if params.AvailableInChatAssistant != nil {
			q.Set("available_in_chat_assistant", strconv.FormatBool(*params.AvailableInChatAssistant))
		}
		if params.Limit != nil {
			q.Set("limit", strconv.Itoa(*params.Limit))
		}
		if params.Cursor != nil {
			q.Set("cursor", *params.Cursor)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.RegistrationListResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowRegistration retrieves a workflow registration by ID.
func (c *Client) GetWorkflowRegistration(ctx context.Context, registrationID string, params *workflow.RegistrationGetParams) (*workflow.Registration, error) {
	path := fmt.Sprintf("/v1/workflows/registrations/%s", registrationID)
	if params != nil {
		q := url.Values{}
		if params.WithWorkflow != nil {
			q.Set("with_workflow", strconv.FormatBool(*params.WithWorkflow))
		}
		if params.IncludeShared != nil {
			q.Set("include_shared", strconv.FormatBool(*params.IncludeShared))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.Registration
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ExecuteWorkflowRegistration executes a workflow via its registration.
//
// Deprecated: Use ExecuteWorkflow instead. This method will be removed in a future release.
func (c *Client) ExecuteWorkflowRegistration(ctx context.Context, registrationID string, req *workflow.ExecutionRequest) (*workflow.ExecutionResponse, error) {
	var resp workflow.ExecutionResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/registrations/%s/execute", registrationID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
