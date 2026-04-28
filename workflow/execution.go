package workflow

import "encoding/json"

// ExecutionStatus is the status of a workflow execution.
type ExecutionStatus string

const (
	ExecutionRunning          ExecutionStatus = "RUNNING"
	ExecutionCompleted        ExecutionStatus = "COMPLETED"
	ExecutionFailed           ExecutionStatus = "FAILED"
	ExecutionCanceled         ExecutionStatus = "CANCELED"
	ExecutionTerminated       ExecutionStatus = "TERMINATED"
	ExecutionContinuedAsNew   ExecutionStatus = "CONTINUED_AS_NEW"
	ExecutionTimedOut         ExecutionStatus = "TIMED_OUT"
	ExecutionRetryingAfterErr ExecutionStatus = "RETRYING_AFTER_ERROR"
)

// ExecutionRequest is the request body for executing a workflow.
type ExecutionRequest struct {
	ExecutionID             *string              `json:"execution_id,omitempty"`
	Input                   map[string]any       `json:"input,omitempty"`
	EncodedInput            *NetworkEncodedInput `json:"encoded_input,omitempty"`
	WaitForResult           bool                 `json:"wait_for_result,omitempty"`
	TimeoutSeconds          *float64             `json:"timeout_seconds,omitempty"`
	CustomTracingAttributes map[string]string    `json:"custom_tracing_attributes,omitempty"`
	DeploymentName          *string              `json:"deployment_name,omitempty"`
	// Extensions carries plugin-specific data such as connector bindings.
	// Use BuildConnectorExtensions to construct the standard connector shape.
	Extensions map[string]any `json:"extensions,omitempty"`
}

// ExecutionResponse is the response from a workflow execution.
type ExecutionResponse struct {
	WorkflowName      string          `json:"workflow_name"`
	ExecutionID       string          `json:"execution_id"`
	RootExecutionID   string          `json:"root_execution_id"`
	Status            ExecutionStatus `json:"status"`
	StartTime         string          `json:"start_time"`
	EndTime           *string         `json:"end_time,omitempty"`
	Result            any             `json:"result,omitempty"`
	ParentExecutionID *string         `json:"parent_execution_id,omitempty"`
	TotalDurationMs   *int            `json:"total_duration_ms,omitempty"`
}

// EncodedPayloadOption identifies how a workflow payload was encoded.
type EncodedPayloadOption string

const (
	EncodedPayloadOffloaded        EncodedPayloadOption = "offloaded"
	EncodedPayloadEncrypted        EncodedPayloadOption = "encrypted"
	EncodedPayloadEncryptedPartial EncodedPayloadOption = "encrypted-partial"
)

// NetworkEncodedInput holds a base64-encoded payload for workflow input.
type NetworkEncodedInput struct {
	B64Payload      string                 `json:"b64payload"`
	EncodingOptions []EncodedPayloadOption `json:"encoding_options,omitempty"`
	Empty           bool                   `json:"empty,omitempty"`
}

// SignalInvocationBody is the request body for signaling a workflow execution.
type SignalInvocationBody struct {
	Name  string `json:"name"`
	Input any    `json:"input"`
}

// SignalResponse is the response from signaling a workflow execution.
type SignalResponse struct {
	Message string `json:"message"`
}

// QueryInvocationBody is the request body for querying a workflow execution.
type QueryInvocationBody struct {
	Name  string `json:"name"`
	Input any    `json:"input,omitempty"`
}

// QueryResponse is the response from querying a workflow execution.
type QueryResponse struct {
	QueryName string `json:"query_name"`
	Result    any    `json:"result"`
}

// UpdateInvocationBody is the request body for updating a workflow execution.
type UpdateInvocationBody struct {
	Name  string `json:"name"`
	Input any    `json:"input,omitempty"`
}

// UpdateResponse is the response from updating a workflow execution.
type UpdateResponse struct {
	UpdateName string `json:"update_name"`
	Result     any    `json:"result"`
}

// ResetInvocationBody is the request body for resetting a workflow execution.
type ResetInvocationBody struct {
	EventID        int     `json:"event_id"`
	Reason         *string `json:"reason,omitempty"`
	ExcludeSignals bool    `json:"exclude_signals,omitempty"`
	ExcludeUpdates bool    `json:"exclude_updates,omitempty"`
}

// BatchExecutionBody is the request body for batch execution operations.
type BatchExecutionBody struct {
	ExecutionIDs []string `json:"execution_ids"`
}

// BatchExecutionResponse is the response from batch execution operations.
type BatchExecutionResponse struct {
	Results map[string]BatchExecutionResult `json:"results,omitempty"`
}

// BatchExecutionResult is the result of a single batch operation.
type BatchExecutionResult struct {
	Status string  `json:"status"`
	Error  *string `json:"error,omitempty"`
}

// StreamParams holds query parameters for streaming workflow executions.
type StreamParams struct {
	EventSource *EventSource
	LastEventID *string
}

// TraceOTelResponse is the response from the OTel trace endpoint.
type TraceOTelResponse struct {
	WorkflowName      string           `json:"workflow_name"`
	ExecutionID       string           `json:"execution_id"`
	RootExecutionID   string           `json:"root_execution_id"`
	Status            *ExecutionStatus `json:"status"`
	StartTime         string           `json:"start_time"`
	EndTime           *string          `json:"end_time,omitempty"`
	Result            any              `json:"result"`
	DataSource        string           `json:"data_source"`
	ParentExecutionID *string          `json:"parent_execution_id,omitempty"`
	TotalDurationMs   *int             `json:"total_duration_ms,omitempty"`
	OTelTraceID       *string          `json:"otel_trace_id,omitempty"`
	OTelTraceData     any              `json:"otel_trace_data,omitempty"`
}

// TraceSummaryResponse is the response from the trace summary endpoint.
type TraceSummaryResponse struct {
	WorkflowName      string           `json:"workflow_name"`
	ExecutionID       string           `json:"execution_id"`
	RootExecutionID   string           `json:"root_execution_id"`
	Status            *ExecutionStatus `json:"status"`
	StartTime         string           `json:"start_time"`
	EndTime           *string          `json:"end_time,omitempty"`
	Result            any              `json:"result"`
	ParentExecutionID *string          `json:"parent_execution_id,omitempty"`
	TotalDurationMs   *int             `json:"total_duration_ms,omitempty"`
	SpanTree          any              `json:"span_tree,omitempty"`
}

// TraceEventsResponse is the response from the trace events endpoint.
type TraceEventsResponse struct {
	WorkflowName      string            `json:"workflow_name"`
	ExecutionID       string            `json:"execution_id"`
	RootExecutionID   string            `json:"root_execution_id"`
	Status            *ExecutionStatus  `json:"status"`
	StartTime         string            `json:"start_time"`
	EndTime           *string           `json:"end_time,omitempty"`
	Result            any               `json:"result"`
	ParentExecutionID *string           `json:"parent_execution_id,omitempty"`
	TotalDurationMs   *int              `json:"total_duration_ms,omitempty"`
	Events            []json.RawMessage `json:"events,omitempty"`
}

// TraceEventsParams holds query parameters for the trace events endpoint.
type TraceEventsParams struct {
	MergeSameIDEvents     *bool
	IncludeInternalEvents *bool
}
