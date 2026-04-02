package mistral

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
)

// StreamWorkflowEvents streams workflow events via SSE.
func (c *Client) StreamWorkflowEvents(ctx context.Context, params *workflow.EventStreamParams) (*WorkflowEventStream, error) {
	path := "/v1/workflows/events/stream"
	if params != nil {
		q := url.Values{}
		if params.Scope != nil {
			q.Set("scope", string(*params.Scope))
		}
		if params.ActivityName != nil {
			q.Set("activity_name", *params.ActivityName)
		}
		if params.ActivityID != nil {
			q.Set("activity_id", *params.ActivityID)
		}
		if params.WorkflowName != nil {
			q.Set("workflow_name", *params.WorkflowName)
		}
		if params.WorkflowExecID != nil {
			q.Set("workflow_exec_id", *params.WorkflowExecID)
		}
		if params.RootWorkflowExecID != nil {
			q.Set("root_workflow_exec_id", *params.RootWorkflowExecID)
		}
		if params.ParentWorkflowExecID != nil {
			q.Set("parent_workflow_exec_id", *params.ParentWorkflowExecID)
		}
		if params.Stream != nil {
			q.Set("stream", *params.Stream)
		}
		if params.StartSeq != nil {
			q.Set("start_seq", strconv.Itoa(*params.StartSeq))
		}
		if params.MetadataFilters != nil {
			data, _ := json.Marshal(params.MetadataFilters)
			q.Set("metadata_filters", string(data))
		}
		if len(params.WorkflowEventTypes) > 0 {
			types := make([]string, len(params.WorkflowEventTypes))
			for i, et := range params.WorkflowEventTypes {
				types[i] = string(et)
			}
			q.Set("workflow_event_types", strings.Join(types, ","))
		}
		if params.LastEventID != nil {
			q.Set("last_event_id", *params.LastEventID)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	resp, err := c.doStream(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	return newWorkflowEventStream(resp.Body), nil
}

// ListWorkflowEvents lists workflow events.
func (c *Client) ListWorkflowEvents(ctx context.Context, params *workflow.EventListParams) (*workflow.EventListResponse, error) {
	path := "/v1/workflows/events/list"
	if params != nil {
		q := url.Values{}
		if params.RootWorkflowExecID != nil {
			q.Set("root_workflow_exec_id", *params.RootWorkflowExecID)
		}
		if params.WorkflowExecID != nil {
			q.Set("workflow_exec_id", *params.WorkflowExecID)
		}
		if params.WorkflowRunID != nil {
			q.Set("workflow_run_id", *params.WorkflowRunID)
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
	var resp workflow.EventListResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
