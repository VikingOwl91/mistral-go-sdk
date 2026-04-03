package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// GetWorkflowExecution retrieves a workflow execution by ID.
func (c *Client) GetWorkflowExecution(ctx context.Context, executionID string) (*workflow.ExecutionResponse, error) {
	var resp workflow.ExecutionResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/executions/%s", executionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowExecutionHistory retrieves the history of a workflow execution.
func (c *Client) GetWorkflowExecutionHistory(ctx context.Context, executionID string, decodePayloads *bool) (json.RawMessage, error) {
	path := fmt.Sprintf("/v1/workflows/executions/%s/history", executionID)
	if decodePayloads != nil {
		path += "?decode_payloads=" + strconv.FormatBool(*decodePayloads)
	}
	var resp json.RawMessage
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// StreamWorkflowExecution streams events for a workflow execution via SSE.
func (c *Client) StreamWorkflowExecution(ctx context.Context, executionID string, params *workflow.StreamParams) (*WorkflowEventStream, error) {
	path := fmt.Sprintf("/v1/workflows/executions/%s/stream", executionID)
	if params != nil {
		q := url.Values{}
		if params.EventSource != nil {
			q.Set("event_source", string(*params.EventSource))
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

// SignalWorkflowExecution sends a signal to a workflow execution.
func (c *Client) SignalWorkflowExecution(ctx context.Context, executionID string, req *workflow.SignalInvocationBody) (*workflow.SignalResponse, error) {
	var resp workflow.SignalResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/signals", executionID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// QueryWorkflowExecution queries a workflow execution.
func (c *Client) QueryWorkflowExecution(ctx context.Context, executionID string, req *workflow.QueryInvocationBody) (*workflow.QueryResponse, error) {
	var resp workflow.QueryResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/queries", executionID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateWorkflowExecution sends an update to a workflow execution.
func (c *Client) UpdateWorkflowExecution(ctx context.Context, executionID string, req *workflow.UpdateInvocationBody) (*workflow.UpdateResponse, error) {
	var resp workflow.UpdateResponse
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/updates", executionID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TerminateWorkflowExecution terminates a workflow execution.
func (c *Client) TerminateWorkflowExecution(ctx context.Context, executionID string) error {
	resp, err := c.do(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/terminate", executionID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// CancelWorkflowExecution cancels a workflow execution.
func (c *Client) CancelWorkflowExecution(ctx context.Context, executionID string) error {
	resp, err := c.do(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/cancel", executionID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// ResetWorkflowExecution resets a workflow execution to a specific event.
func (c *Client) ResetWorkflowExecution(ctx context.Context, executionID string, req *workflow.ResetInvocationBody) error {
	return c.doJSON(ctx, "POST", fmt.Sprintf("/v1/workflows/executions/%s/reset", executionID), req, nil)
}

// BatchCancelWorkflowExecutions cancels multiple workflow executions.
func (c *Client) BatchCancelWorkflowExecutions(ctx context.Context, req *workflow.BatchExecutionBody) (*workflow.BatchExecutionResponse, error) {
	var resp workflow.BatchExecutionResponse
	if err := c.doJSON(ctx, "POST", "/v1/workflows/executions/cancel", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// BatchTerminateWorkflowExecutions terminates multiple workflow executions.
func (c *Client) BatchTerminateWorkflowExecutions(ctx context.Context, req *workflow.BatchExecutionBody) (*workflow.BatchExecutionResponse, error) {
	var resp workflow.BatchExecutionResponse
	if err := c.doJSON(ctx, "POST", "/v1/workflows/executions/terminate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowExecutionTraceOTel retrieves the OpenTelemetry trace for a workflow execution.
func (c *Client) GetWorkflowExecutionTraceOTel(ctx context.Context, executionID string) (*workflow.TraceOTelResponse, error) {
	var resp workflow.TraceOTelResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/executions/%s/trace/otel", executionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowExecutionTraceSummary retrieves the trace summary for a workflow execution.
func (c *Client) GetWorkflowExecutionTraceSummary(ctx context.Context, executionID string) (*workflow.TraceSummaryResponse, error) {
	var resp workflow.TraceSummaryResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/workflows/executions/%s/trace/summary", executionID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetWorkflowExecutionTraceEvents retrieves the trace events for a workflow execution.
func (c *Client) GetWorkflowExecutionTraceEvents(ctx context.Context, executionID string, params *workflow.TraceEventsParams) (*workflow.TraceEventsResponse, error) {
	path := fmt.Sprintf("/v1/workflows/executions/%s/trace/events", executionID)
	if params != nil {
		q := url.Values{}
		if params.MergeSameIDEvents != nil {
			q.Set("merge_same_id_events", strconv.FormatBool(*params.MergeSameIDEvents))
		}
		if params.IncludeInternalEvents != nil {
			q.Set("include_internal_events", strconv.FormatBool(*params.IncludeInternalEvents))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp workflow.TraceEventsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// WorkflowEventStream wraps the generic Stream to provide typed workflow events
// with StreamPayload envelope metadata.
type WorkflowEventStream struct {
	stream  *Stream[json.RawMessage]
	event   workflow.Event
	payload *workflow.StreamPayload
	err     error
}

func newWorkflowEventStream(body readCloser) *WorkflowEventStream {
	return &WorkflowEventStream{
		stream: newStream[json.RawMessage](body),
	}
}

// Next advances to the next event. Returns false when done or on error.
func (s *WorkflowEventStream) Next() bool {
	if s.err != nil {
		return false
	}
	if !s.stream.Next() {
		s.err = s.stream.Err()
		return false
	}
	var payload workflow.StreamPayload
	if err := json.Unmarshal(s.stream.Current(), &payload); err != nil {
		s.err = fmt.Errorf("mistral: decode workflow stream payload: %w", err)
		return false
	}
	event, err := workflow.UnmarshalEvent(payload.Data)
	if err != nil {
		s.err = err
		return false
	}
	s.event = event
	s.payload = &payload
	return true
}

// Current returns the most recently read workflow event.
func (s *WorkflowEventStream) Current() workflow.Event { return s.event }

// CurrentPayload returns the full StreamPayload envelope of the current event.
func (s *WorkflowEventStream) CurrentPayload() *workflow.StreamPayload { return s.payload }

// Err returns any error encountered during streaming.
func (s *WorkflowEventStream) Err() error { return s.err }

// Close releases the underlying connection.
func (s *WorkflowEventStream) Close() error { return s.stream.Close() }

// ExecuteWorkflowAndWait executes a workflow and polls until completion.
func (c *Client) ExecuteWorkflowAndWait(ctx context.Context, workflowIdentifier string, req *workflow.ExecutionRequest) (*workflow.ExecutionResponse, error) {
	execResp, err := c.ExecuteWorkflow(ctx, workflowIdentifier, req)
	if err != nil {
		return nil, err
	}
	for {
		if isTerminal(execResp.Status) {
			return execResp, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(500 * time.Millisecond):
		}
		execResp, err = c.GetWorkflowExecution(ctx, execResp.ExecutionID)
		if err != nil {
			return nil, err
		}
	}
}

func isTerminal(s workflow.ExecutionStatus) bool {
	switch s {
	case workflow.ExecutionCompleted, workflow.ExecutionFailed,
		workflow.ExecutionCanceled, workflow.ExecutionTerminated,
		workflow.ExecutionTimedOut:
		return true
	}
	return false
}
