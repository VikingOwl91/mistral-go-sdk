# Workflows API Integration — Design Spec

**Date:** 2026-04-02
**Upstream:** Mistral Python SDK v2.2.0 (released 2026-03-31)
**SDK version:** v1.2.0
**Scope:** Full parity with Python SDK v2.2.0 changes since v2.1.3

## Summary

Add the Workflows API (37 new methods) and `DeleteBatchJob` (1 method) to the Go SDK.
This is purely additive — no breaking changes to existing API surface.

## New Package: `workflow/`

Types-only package following the two-layer architecture. 8 type files + `doc.go`.

### `workflow/doc.go`

Package documentation.

### `workflow/workflow.go` — Core CRUD types

```go
type Workflow struct {
    ID                       string  `json:"id"`
    Name                     string  `json:"name"`
    DisplayName              *string `json:"display_name,omitempty"`
    Description              *string `json:"description,omitempty"`
    OwnerID                  string  `json:"owner_id"`
    WorkspaceID              string  `json:"workspace_id"`
    AvailableInChatAssistant bool    `json:"available_in_chat_assistant"`
    Archived                 bool    `json:"archived"`
    CreatedAt                string  `json:"created_at"`
    UpdatedAt                string  `json:"updated_at"`
}

type WorkflowUpdateRequest struct {
    DisplayName              *string `json:"display_name,omitempty"`
    Description              *string `json:"description,omitempty"`
    AvailableInChatAssistant *bool   `json:"available_in_chat_assistant,omitempty"`
}

type WorkflowListResponse struct {
    Workflows  []Workflow `json:"workflows"`
    NextCursor *string    `json:"next_cursor,omitempty"`
}

type WorkflowListParams struct {
    ActiveOnly               *bool
    IncludeShared            *bool
    AvailableInChatAssistant *bool
    Archived                 *bool
    Cursor                   *string
    Limit                    *int
}

type WorkflowArchiveResponse struct {
    ID       string `json:"id"`
    Archived bool   `json:"archived"`
}
```

### `workflow/execution.go` — Execution types

```go
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

type ExecutionRequest struct {
    ExecutionID              *string           `json:"execution_id,omitempty"`
    Input                    map[string]any    `json:"input,omitempty"`
    EncodedInput             *NetworkEncodedInput `json:"encoded_input,omitempty"`
    WaitForResult            bool              `json:"wait_for_result,omitempty"`
    TimeoutSeconds           *float64          `json:"timeout_seconds,omitempty"`
    CustomTracingAttributes  map[string]string `json:"custom_tracing_attributes,omitempty"`
    DeploymentName           *string           `json:"deployment_name,omitempty"`
}

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

type NetworkEncodedInput struct {
    B64Payload      string   `json:"b64payload"`
    EncodingOptions []string `json:"encoding_options,omitempty"`
    Empty           bool     `json:"empty,omitempty"`
}

type SignalInvocationBody struct {
    Name  string `json:"name"`
    Input any    `json:"input"`
}

type SignalResponse struct {
    Message string `json:"message"` // default: "Signal accepted"
}

type QueryInvocationBody struct {
    Name  string `json:"name"`
    Input any    `json:"input,omitempty"`
}

type QueryResponse struct {
    QueryName string `json:"query_name"`
    Result    any    `json:"result"`
}

type UpdateInvocationBody struct {
    Name  string `json:"name"`
    Input any    `json:"input,omitempty"`
}

type UpdateResponse struct {
    UpdateName string `json:"update_name"`
    Result     any    `json:"result"`
}

// Trace response types

type TraceOTelResponse struct {
    WorkflowName      string          `json:"workflow_name"`
    ExecutionID       string          `json:"execution_id"`
    RootExecutionID   string          `json:"root_execution_id"`
    Status            *ExecutionStatus `json:"status"`
    StartTime         string          `json:"start_time"`
    EndTime           *string         `json:"end_time,omitempty"`
    Result            any             `json:"result"`
    DataSource        string          `json:"data_source"`
    ParentExecutionID *string         `json:"parent_execution_id,omitempty"`
    TotalDurationMs   *int            `json:"total_duration_ms,omitempty"`
    OTelTraceID       *string         `json:"otel_trace_id,omitempty"`
    OTelTraceData     any             `json:"otel_trace_data,omitempty"`
}

type TraceSummaryResponse struct {
    WorkflowName      string          `json:"workflow_name"`
    ExecutionID       string          `json:"execution_id"`
    RootExecutionID   string          `json:"root_execution_id"`
    Status            *ExecutionStatus `json:"status"`
    StartTime         string          `json:"start_time"`
    EndTime           *string         `json:"end_time,omitempty"`
    Result            any             `json:"result"`
    ParentExecutionID *string         `json:"parent_execution_id,omitempty"`
    TotalDurationMs   *int            `json:"total_duration_ms,omitempty"`
    SpanTree          any             `json:"span_tree,omitempty"`
}

type TraceEventsResponse struct {
    WorkflowName      string          `json:"workflow_name"`
    ExecutionID       string          `json:"execution_id"`
    RootExecutionID   string          `json:"root_execution_id"`
    Status            *ExecutionStatus `json:"status"`
    StartTime         string          `json:"start_time"`
    EndTime           *string         `json:"end_time,omitempty"`
    Result            any             `json:"result"`
    ParentExecutionID *string         `json:"parent_execution_id,omitempty"`
    TotalDurationMs   *int            `json:"total_duration_ms,omitempty"`
    Events            []json.RawMessage `json:"events,omitempty"`
}

type TraceEventsParams struct {
    MergeSameIDEvents     *bool
    IncludeInternalEvents *bool
}

type ResetInvocationBody struct {
    EventID        int     `json:"event_id"`
    Reason         *string `json:"reason,omitempty"`
    ExcludeSignals bool    `json:"exclude_signals,omitempty"`
    ExcludeUpdates bool    `json:"exclude_updates,omitempty"`
}

type BatchExecutionBody struct {
    ExecutionIDs []string `json:"execution_ids"`
}

type BatchExecutionResponse struct {
    Results map[string]BatchExecutionResult `json:"results,omitempty"`
}

type BatchExecutionResult struct {
    Status string  `json:"status"`
    Error  *string `json:"error,omitempty"`
}

type StreamParams struct {
    EventSource *EventSource
    LastEventID *string
}
```

### `workflow/event.go` — Sealed event interface + 17 variants

Discriminator field: `event_type`

```go
type Event interface {
    workflowEvent()
    EventType() EventType
}

type EventType string

const (
    EventWorkflowStarted       EventType = "WORKFLOW_EXECUTION_STARTED"
    EventWorkflowCompleted     EventType = "WORKFLOW_EXECUTION_COMPLETED"
    EventWorkflowFailed        EventType = "WORKFLOW_EXECUTION_FAILED"
    EventWorkflowCanceled      EventType = "WORKFLOW_EXECUTION_CANCELED"
    EventWorkflowContinuedAsNew EventType = "WORKFLOW_EXECUTION_CONTINUED_AS_NEW"
    EventWorkflowTaskTimedOut  EventType = "WORKFLOW_TASK_TIMED_OUT"
    EventWorkflowTaskFailed    EventType = "WORKFLOW_TASK_FAILED"
    EventCustomTaskStarted     EventType = "CUSTOM_TASK_STARTED"
    EventCustomTaskInProgress  EventType = "CUSTOM_TASK_IN_PROGRESS"
    EventCustomTaskCompleted   EventType = "CUSTOM_TASK_COMPLETED"
    EventCustomTaskFailed      EventType = "CUSTOM_TASK_FAILED"
    EventCustomTaskTimedOut    EventType = "CUSTOM_TASK_TIMED_OUT"
    EventCustomTaskCanceled    EventType = "CUSTOM_TASK_CANCELED"
    EventActivityTaskStarted   EventType = "ACTIVITY_TASK_STARTED"
    EventActivityTaskCompleted EventType = "ACTIVITY_TASK_COMPLETED"
    EventActivityTaskRetrying  EventType = "ACTIVITY_TASK_RETRYING"
    EventActivityTaskFailed    EventType = "ACTIVITY_TASK_FAILED"
)

type EventSource string

const (
    EventSourceDatabase EventSource = "DATABASE"
    EventSourceLive     EventSource = "LIVE"
)

type Scope string

const (
    ScopeActivity Scope = "activity"
    ScopeWorkflow Scope = "workflow"
    ScopeAll      Scope = "*"
)
```

Each concrete event type has common fields + type-specific attributes:

```go
// Common fields embedded in all event types
type eventBase struct {
    ID                    string  `json:"event_id"`
    Timestamp             int64   `json:"event_timestamp"`
    RootWorkflowExecID    string  `json:"root_workflow_exec_id"`
    ParentWorkflowExecID  *string `json:"parent_workflow_exec_id"`
    WorkflowExecID        string  `json:"workflow_exec_id"`
    WorkflowRunID         string  `json:"workflow_run_id"`
    WorkflowName          string  `json:"workflow_name"`
}

// Example concrete types:

type WorkflowExecutionStartedEvent struct {
    eventBase
    Attributes WorkflowStartedAttributes `json:"attributes"`
}
func (WorkflowExecutionStartedEvent) workflowEvent()          {}
func (WorkflowExecutionStartedEvent) EventType() EventType    { return EventWorkflowStarted }

type WorkflowExecutionCompletedEvent struct {
    eventBase
    Attributes WorkflowCompletedAttributes `json:"attributes"`
}
// ... pattern repeats for all 17 types

type UnknownEvent struct {
    eventBase
    RawType string
    Raw     json.RawMessage
}
```

SSE envelope types:

```go
type StreamPayload struct {
    Stream          string               `json:"stream"`
    Data            json.RawMessage      `json:"data"`
    WorkflowContext StreamWorkflowContext `json:"workflow_context"`
    BrokerSequence  int                  `json:"broker_sequence"`
    Timestamp       *string              `json:"timestamp,omitempty"`
    Metadata        map[string]any       `json:"metadata,omitempty"`
}

type StreamWorkflowContext struct {
    Namespace             string  `json:"namespace"`
    WorkflowName          string  `json:"workflow_name"`
    WorkflowExecID        string  `json:"workflow_exec_id"`
    ParentWorkflowExecID  *string `json:"parent_workflow_exec_id,omitempty"`
    RootWorkflowExecID    *string `json:"root_workflow_exec_id,omitempty"`
}

func UnmarshalEvent(data json.RawMessage) (Event, error)
// Probes event_type discriminator, dispatches to concrete type.
// Unknown event_type returns UnknownEvent.
```

### `workflow/deployment.go`

```go
type Deployment struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    IsActive  bool   `json:"is_active"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type DeploymentListResponse struct {
    Deployments []Deployment `json:"deployments"`
}

type DeploymentListParams struct {
    ActiveOnly   *bool
    WorkflowName *string
}
```

### `workflow/metrics.go`

```go
type Metrics struct {
    ExecutionCount   ScalarMetric     `json:"execution_count"`
    SuccessCount     ScalarMetric     `json:"success_count"`
    ErrorCount       ScalarMetric     `json:"error_count"`
    AverageLatencyMs ScalarMetric     `json:"average_latency_ms"`
    LatencyOverTime  TimeSeriesMetric `json:"latency_over_time"`
    RetryRate        ScalarMetric     `json:"retry_rate"`
}

type ScalarMetric struct {
    Value float64 `json:"value"`
}

type TimeSeriesMetric struct {
    Value [][]float64 `json:"value"`
}

type MetricsParams struct {
    StartTime *string
    EndTime   *string
}
```

### `workflow/run.go`

```go
type Run struct {
    ID                string          `json:"id"`
    WorkflowName      string          `json:"workflow_name"`
    ExecutionID       string          `json:"execution_id"`
    Status            ExecutionStatus `json:"status"`
    StartTime         string          `json:"start_time"`
    EndTime           *string         `json:"end_time,omitempty"`
}

type ListRunsResponse struct {
    Runs          []Run   `json:"runs"`
    NextPageToken *string `json:"next_page_token,omitempty"`
}

type RunListParams struct {
    WorkflowIdentifier *string
    Search             *string
    Status             *string
    PageSize           *int
    NextPageToken      *string
}
```

### `workflow/schedule.go`

```go
type ScheduleRequest struct {
    Schedule               ScheduleDefinition `json:"schedule"`
    WorkflowRegistrationID *string            `json:"workflow_registration_id,omitempty"`
    WorkflowIdentifier     *string            `json:"workflow_identifier,omitempty"`
    ScheduleID             *string            `json:"schedule_id,omitempty"`
    DeploymentName         *string            `json:"deployment_name,omitempty"`
}

type ScheduleDefinition struct {
    Input           any                `json:"input"`
    Calendars       []ScheduleCalendar `json:"calendars,omitempty"`
    Intervals       []ScheduleInterval `json:"intervals,omitempty"`
    CronExpressions []string           `json:"cron_expressions,omitempty"`
    Skip            []ScheduleCalendar `json:"skip,omitempty"`
    StartAt         *string            `json:"start_at,omitempty"`
    EndAt           *string            `json:"end_at,omitempty"`
    Jitter          *string            `json:"jitter,omitempty"`
    TimeZoneName    *string            `json:"time_zone_name,omitempty"`
    Policy          *SchedulePolicy    `json:"policy,omitempty"`
}

type ScheduleCalendar struct {
    Second     []ScheduleRange `json:"second,omitempty"`
    Minute     []ScheduleRange `json:"minute,omitempty"`
    Hour       []ScheduleRange `json:"hour,omitempty"`
    DayOfMonth []ScheduleRange `json:"day_of_month,omitempty"`
    Month      []ScheduleRange `json:"month,omitempty"`
    Year       []ScheduleRange `json:"year,omitempty"`
    DayOfWeek  []ScheduleRange `json:"day_of_week,omitempty"`
    Comment    *string         `json:"comment,omitempty"`
}

type ScheduleRange struct {
    Start int `json:"start"`
    End   int `json:"end,omitempty"`
    Step  int `json:"step,omitempty"`
}

type ScheduleInterval struct {
    Every  string  `json:"every"`
    Offset *string `json:"offset,omitempty"`
}

type SchedulePolicy struct {
    CatchupWindowSeconds int  `json:"catchup_window_seconds,omitempty"`
    Overlap              *int `json:"overlap,omitempty"`
    PauseOnFailure       bool `json:"pause_on_failure,omitempty"`
}

type ScheduleResponse struct {
    ScheduleID string `json:"schedule_id"`
}

type ScheduleListResponse struct {
    Schedules []Schedule `json:"schedules"`
}

type Schedule struct {
    ScheduleID   string             `json:"schedule_id"`
    Definition   ScheduleDefinition `json:"definition"`
    WorkflowName string             `json:"workflow_name"`
    CreatedAt    string             `json:"created_at"`
    UpdatedAt    string             `json:"updated_at"`
}
```

### `workflow/registration.go`

```go
type Registration struct {
    ID          string    `json:"id"`
    WorkflowID  string    `json:"workflow_id"`
    TaskQueue   string    `json:"task_queue"`
    Workflow    *Workflow  `json:"workflow,omitempty"`
    CreatedAt   string    `json:"created_at"`
    UpdatedAt   string    `json:"updated_at"`
}

type RegistrationListResponse struct {
    Registrations []Registration `json:"registrations"`
    NextCursor    *string        `json:"next_cursor,omitempty"`
}

type RegistrationListParams struct {
    WorkflowID               *string
    TaskQueue                *string
    ActiveOnly               *bool
    IncludeShared            *bool
    WorkflowSearch           *string
    Archived                 *bool
    WithWorkflow             *bool
    AvailableInChatAssistant *bool
    Limit                    *int
    Cursor                   *string
}

type RegistrationGetParams struct {
    WithWorkflow  *bool
    IncludeShared *bool
}

type WorkerInfo struct {
    SchedulerURL string `json:"scheduler_url"`
    Namespace    string `json:"namespace"`
    TLS          bool   `json:"tls"`
}
```

## Service Methods (root package)

### `workflows.go` — 10 methods

| Method | HTTP | Path |
|--------|------|------|
| `ListWorkflows` | GET | `/v1/workflows` |
| `GetWorkflow` | GET | `/v1/workflows/{id}` |
| `UpdateWorkflow` | PUT | `/v1/workflows/{id}` |
| `ArchiveWorkflow` | PUT | `/v1/workflows/{id}/archive` |
| `UnarchiveWorkflow` | PUT | `/v1/workflows/{id}/unarchive` |
| `ExecuteWorkflow` | POST | `/v1/workflows/{id}/execute` |
| `ListWorkflowRegistrations` | GET | `/v1/workflows/registrations` |
| `GetWorkflowRegistration` | GET | `/v1/workflows/registrations/{id}` |
| `ExecuteWorkflowRegistration` | POST | `/v1/workflows/registrations/{id}/execute` |
| `ExecuteWorkflowAndWait` | (composite) | execute + poll |

`ExecuteWorkflowRegistration` is deprecated (doc comment only).

`ExecuteWorkflowAndWait` calls `ExecuteWorkflow`, then polls `GetWorkflowExecution`
in a loop until status is terminal or context is canceled.

### `workflows_executions.go` — 14 methods

| Method | HTTP | Path |
|--------|------|------|
| `GetWorkflowExecution` | GET | `/v1/workflows/executions/{id}` |
| `GetWorkflowExecutionHistory` | GET | `/v1/workflows/executions/{id}/history` |
| `StreamWorkflowExecution` | GET (SSE) | `/v1/workflows/executions/{id}/stream` |
| `SignalWorkflowExecution` | POST | `/v1/workflows/executions/{id}/signals` |
| `QueryWorkflowExecution` | POST | `/v1/workflows/executions/{id}/queries` |
| `UpdateWorkflowExecution` | POST | `/v1/workflows/executions/{id}/updates` |
| `TerminateWorkflowExecution` | POST (204) | `/v1/workflows/executions/{id}/terminate` |
| `CancelWorkflowExecution` | POST (204) | `/v1/workflows/executions/{id}/cancel` |
| `ResetWorkflowExecution` | POST (204) | `/v1/workflows/executions/{id}/reset` |
| `BatchCancelWorkflowExecutions` | POST | `/v1/workflows/executions/cancel` |
| `BatchTerminateWorkflowExecutions` | POST | `/v1/workflows/executions/terminate` |
| `GetWorkflowExecutionTraceOTel` | GET | `/v1/workflows/executions/{id}/trace/otel` |
| `GetWorkflowExecutionTraceSummary` | GET | `/v1/workflows/executions/{id}/trace/summary` |
| `GetWorkflowExecutionTraceEvents` | GET | `/v1/workflows/executions/{id}/trace/events` |

Also contains `WorkflowEventStream` type (wraps `Stream[json.RawMessage]`,
dispatches via `workflow.UnmarshalEvent`).

### `workflows_events.go` — 2 methods

| Method | HTTP | Path |
|--------|------|------|
| `StreamWorkflowEvents` | GET (SSE) | `/v1/workflows/events/stream` |
| `ListWorkflowEvents` | GET | `/v1/workflows/events/list` |

### `workflows_deployments.go` — 2 methods

| Method | HTTP | Path |
|--------|------|------|
| `ListWorkflowDeployments` | GET | `/v1/workflows/deployments` |
| `GetWorkflowDeployment` | GET | `/v1/workflows/deployments/{id}` |

### `workflows_metrics.go` — 1 method

| Method | HTTP | Path |
|--------|------|------|
| `GetWorkflowMetrics` | GET | `/v1/workflows/{name}/metrics` |

### `workflows_runs.go` — 3 methods

| Method | HTTP | Path |
|--------|------|------|
| `ListWorkflowRuns` | GET | `/v1/workflows/runs` |
| `GetWorkflowRun` | GET | `/v1/workflows/runs/{id}` |
| `GetWorkflowRunHistory` | GET | `/v1/workflows/runs/{id}/history` |

### `workflows_schedules.go` — 3 methods

| Method | HTTP | Path |
|--------|------|------|
| `ListWorkflowSchedules` | GET | `/v1/workflows/schedules` |
| `ScheduleWorkflow` | POST | `/v1/workflows/schedules` |
| `UnscheduleWorkflow` | DELETE | `/v1/workflows/schedules/{id}` |

### `workflows_workers.go` — 1 method

| Method | HTTP | Path |
|--------|------|------|
| `GetWorkflowWorkerInfo` | GET | `/v1/workflows/workers/whoami` |

### `batch_api.go` — 1 new method

| Method | HTTP | Path |
|--------|------|------|
| `DeleteBatchJob` | DELETE | `/v1/batch/jobs/{id}` |

New type in `batch/`: `DeleteResponse { ID, Object, Deleted }`.

## Streaming Design

`WorkflowEventStream` wraps `Stream[json.RawMessage]` like `EventStream` does for conversations.

SSE data arrives as `StreamPayload` envelope:
```json
{
  "stream": "...",
  "data": { "event_type": "WORKFLOW_EXECUTION_COMPLETED", ... },
  "workflow_context": { ... },
  "broker_sequence": 42,
  "timestamp": "...",
  "metadata": {}
}
```

`WorkflowEventStream.Next()`:
1. Read next SSE `data:` line via inner `Stream[json.RawMessage]`
2. Unmarshal into `workflow.StreamPayload`
3. Dispatch `payload.Data` via `workflow.UnmarshalEvent` (probes `event_type`)
4. Expose both `Current() workflow.Event` and `CurrentPayload() *workflow.StreamPayload`

Both `StreamWorkflowExecution` and `StreamWorkflowEvents` use GET (not POST)
with SSE response. They use `doStream` without a request body — the stream method
needs to support GET + query params (verify `doStream` handles nil body for GET).

## Testing

One test file per service file. `httptest.NewServer` with inline handlers. Stdlib only.

Key scenarios:
- Query param encoding for list/filter endpoints
- PUT body marshaling for update/archive
- 204 no-body responses for terminate/cancel/reset
- 202 response for signal
- SSE streaming with StreamPayload envelope + event type dispatch
- UnknownEvent forward compatibility
- ExecuteWorkflowAndWait polling loop (mock multiple get-execution responses)
- Sealed interface UnmarshalEvent for all 17 event types + unknown
- Batch operations with map response

## Version & Docs

- Bump version constant in `mistral.go` to `1.2.0`
- Update `CLAUDE.md` sub-packages list to include `workflow/`
- Update `CHANGELOG.md` with v1.2.0 entry
- Upstream sync reference: Python SDK v2.2.0

## Non-Goals

- No pagination helpers (cursor chaining) — callers manage pagination manually, same as existing endpoints
- No traceparent injection hook — Go callers manage their own tracing headers
- No `execute_workflow_and_wait_async` — Go has context cancellation instead
- No WebSocket/realtime workflow support (not in Python SDK either)
