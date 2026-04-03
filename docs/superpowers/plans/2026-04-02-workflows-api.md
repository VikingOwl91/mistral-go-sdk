# Workflows API Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Workflows API (37 methods) and DeleteBatchJob (1 method) to the Go SDK, reaching full parity with Python SDK v2.2.0.

**Architecture:** Types-only `workflow/` sub-package with 9 files (doc.go + 8 type files). 8 service files on `*Client` in the root package, one per sub-resource. `WorkflowEventStream` wraps `Stream[json.RawMessage]` with `StreamPayload` envelope unwrapping. Sealed `Event` interface with 17 concrete types + `UnknownEvent` fallback.

**Tech Stack:** Go 1.26+, stdlib only (net/http, encoding/json, net/http/httptest for tests).

**Spec:** `docs/superpowers/specs/2026-04-02-workflows-api-design.md`

---

## File Map

### New files — `workflow/` package (types only)

| File | Responsibility |
|------|---------------|
| `workflow/doc.go` | Package documentation |
| `workflow/workflow.go` | Core CRUD types: Workflow, WorkflowUpdateRequest, WorkflowListResponse, WorkflowListParams, WorkflowArchiveResponse |
| `workflow/execution.go` | Execution types: ExecutionStatus enum, ExecutionRequest/Response, Signal/Query/Update bodies+responses, Reset, Batch, NetworkEncodedInput, StreamParams, Trace responses |
| `workflow/event.go` | Sealed Event interface, 17 concrete event types, UnknownEvent, UnmarshalEvent dispatcher, EventType/EventSource/Scope enums, StreamPayload, StreamWorkflowContext |
| `workflow/event_test.go` | UnmarshalEvent tests for all 17 types + unknown fallback |
| `workflow/deployment.go` | Deployment, DeploymentListResponse, DeploymentListParams |
| `workflow/metrics.go` | Metrics, ScalarMetric, TimeSeriesMetric, MetricsParams |
| `workflow/run.go` | Run, ListRunsResponse, RunListParams |
| `workflow/schedule.go` | ScheduleRequest, ScheduleDefinition, ScheduleCalendar, ScheduleRange, ScheduleInterval, SchedulePolicy, ScheduleResponse, ScheduleListResponse, Schedule |
| `workflow/registration.go` | Registration, RegistrationListResponse, RegistrationListParams, RegistrationGetParams, WorkerInfo |

### New files — root package (service methods + tests)

| File | Responsibility |
|------|---------------|
| `workflows.go` | 10 methods: List/Get/Update/Archive/Unarchive Workflow, ExecuteWorkflow, ExecuteWorkflowAndWait, List/Get WorkflowRegistration, ExecuteWorkflowRegistration |
| `workflows_test.go` | Tests for all 10 methods |
| `workflows_executions.go` | 14 methods: Get/History/Stream/Signal/Query/Update/Terminate/Cancel/Reset execution, BatchCancel/BatchTerminate, 3 trace endpoints. Plus WorkflowEventStream type. |
| `workflows_executions_test.go` | Tests for all 14 methods |
| `workflows_events.go` | 2 methods: StreamWorkflowEvents, ListWorkflowEvents |
| `workflows_events_test.go` | Tests for both methods |
| `workflows_deployments.go` | 2 methods: ListWorkflowDeployments, GetWorkflowDeployment |
| `workflows_deployments_test.go` | Tests for both methods |
| `workflows_metrics.go` | 1 method: GetWorkflowMetrics |
| `workflows_metrics_test.go` | Test for GetWorkflowMetrics |
| `workflows_runs.go` | 3 methods: ListWorkflowRuns, GetWorkflowRun, GetWorkflowRunHistory |
| `workflows_runs_test.go` | Tests for all 3 methods |
| `workflows_schedules.go` | 3 methods: ListWorkflowSchedules, ScheduleWorkflow, UnscheduleWorkflow |
| `workflows_schedules_test.go` | Tests for all 3 methods |
| `workflows_workers.go` | 1 method: GetWorkflowWorkerInfo |
| `workflows_workers_test.go` | Test for GetWorkflowWorkerInfo |

### Modified files

| File | Change |
|------|--------|
| `batch_api.go` | Add DeleteBatchJob method |
| `batch/batch.go` (or similar) | Add DeleteResponse type |
| `batch_test.go` | Add TestDeleteBatchJob_Success |
| `mistral.go` | Bump Version to "1.2.0" |
| `CHANGELOG.md` | Add v1.2.0 entry |
| `CLAUDE.md` | Add `workflow/` to sub-packages list |

---

## Task 1: Package scaffold and core types

**Files:**
- Create: `workflow/doc.go`
- Create: `workflow/workflow.go`

- [ ] **Step 1: Create `workflow/doc.go`**

```go
// Package workflow provides types for the Mistral workflows API.
//
// Workflows support orchestrating multi-step processes with execution
// management, scheduling, event streaming, and observability.
package workflow
```

- [ ] **Step 2: Create `workflow/workflow.go` with core CRUD types**

```go
package workflow

// Workflow represents a workflow definition.
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

// WorkflowUpdateRequest is the request body for updating a workflow.
type WorkflowUpdateRequest struct {
	DisplayName              *string `json:"display_name,omitempty"`
	Description              *string `json:"description,omitempty"`
	AvailableInChatAssistant *bool   `json:"available_in_chat_assistant,omitempty"`
}

// WorkflowListResponse is the response from listing workflows.
type WorkflowListResponse struct {
	Workflows  []Workflow `json:"workflows"`
	NextCursor *string    `json:"next_cursor,omitempty"`
}

// WorkflowListParams holds query parameters for listing workflows.
type WorkflowListParams struct {
	ActiveOnly               *bool
	IncludeShared            *bool
	AvailableInChatAssistant *bool
	Archived                 *bool
	Cursor                   *string
	Limit                    *int
}

// WorkflowArchiveResponse is the response from archiving/unarchiving a workflow.
type WorkflowArchiveResponse struct {
	ID       string `json:"id"`
	Archived bool   `json:"archived"`
}
```

- [ ] **Step 3: Verify it compiles**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go build ./workflow/`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add workflow/doc.go workflow/workflow.go
git commit -m "feat(workflow): add package scaffold and core CRUD types"
```

---

## Task 2: Execution types

**Files:**
- Create: `workflow/execution.go`

- [ ] **Step 1: Create `workflow/execution.go`**

```go
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
	ExecutionID             *string           `json:"execution_id,omitempty"`
	Input                   map[string]any    `json:"input,omitempty"`
	EncodedInput            *NetworkEncodedInput `json:"encoded_input,omitempty"`
	WaitForResult           bool              `json:"wait_for_result,omitempty"`
	TimeoutSeconds          *float64          `json:"timeout_seconds,omitempty"`
	CustomTracingAttributes map[string]string `json:"custom_tracing_attributes,omitempty"`
	DeploymentName          *string           `json:"deployment_name,omitempty"`
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

// NetworkEncodedInput holds a base64-encoded payload for workflow input.
type NetworkEncodedInput struct {
	B64Payload      string   `json:"b64payload"`
	EncodingOptions []string `json:"encoding_options,omitempty"`
	Empty           bool     `json:"empty,omitempty"`
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
	WorkflowName      string           `json:"workflow_name"`
	ExecutionID       string           `json:"execution_id"`
	RootExecutionID   string           `json:"root_execution_id"`
	Status            *ExecutionStatus `json:"status"`
	StartTime         string           `json:"start_time"`
	EndTime           *string          `json:"end_time,omitempty"`
	Result            any              `json:"result"`
	ParentExecutionID *string          `json:"parent_execution_id,omitempty"`
	TotalDurationMs   *int             `json:"total_duration_ms,omitempty"`
	Events            []json.RawMessage `json:"events,omitempty"`
}

// TraceEventsParams holds query parameters for the trace events endpoint.
type TraceEventsParams struct {
	MergeSameIDEvents     *bool
	IncludeInternalEvents *bool
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go build ./workflow/`
Expected: no errors

- [ ] **Step 3: Commit**

```bash
git add workflow/execution.go
git commit -m "feat(workflow): add execution types"
```

---

## Task 3: Event sealed interface and UnmarshalEvent

**Files:**
- Create: `workflow/event.go`
- Create: `workflow/event_test.go`

- [ ] **Step 1: Write the failing test in `workflow/event_test.go`**

```go
package workflow

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEvent_WorkflowExecutionCompleted(t *testing.T) {
	data := []byte(`{
		"event_id": "evt-1",
		"event_timestamp": 1711929600000000000,
		"root_workflow_exec_id": "root-1",
		"parent_workflow_exec_id": null,
		"workflow_exec_id": "exec-1",
		"workflow_run_id": "run-1",
		"workflow_name": "my-workflow",
		"event_type": "WORKFLOW_EXECUTION_COMPLETED",
		"attributes": {"task_id": "t1", "result": {"value": {"answer": 42}, "type": "json"}}
	}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*WorkflowExecutionCompletedEvent)
	if !ok {
		t.Fatalf("expected *WorkflowExecutionCompletedEvent, got %T", event)
	}
	if e.ID != "evt-1" {
		t.Errorf("got ID %q", e.ID)
	}
	if e.WorkflowName != "my-workflow" {
		t.Errorf("got WorkflowName %q", e.WorkflowName)
	}
	if e.EventType() != EventWorkflowCompleted {
		t.Errorf("got EventType %q", e.EventType())
	}
}

func TestUnmarshalEvent_CustomTaskStarted(t *testing.T) {
	data := []byte(`{
		"event_id": "evt-2",
		"event_timestamp": 1711929600000000000,
		"root_workflow_exec_id": "root-1",
		"parent_workflow_exec_id": "parent-1",
		"workflow_exec_id": "exec-1",
		"workflow_run_id": "run-1",
		"workflow_name": "my-workflow",
		"event_type": "CUSTOM_TASK_STARTED",
		"attributes": {}
	}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*CustomTaskStartedEvent)
	if !ok {
		t.Fatalf("expected *CustomTaskStartedEvent, got %T", event)
	}
	parent := "parent-1"
	if e.ParentWorkflowExecID == nil || *e.ParentWorkflowExecID != parent {
		t.Errorf("expected parent %q, got %v", parent, e.ParentWorkflowExecID)
	}
}

func TestUnmarshalEvent_ActivityTaskRetrying(t *testing.T) {
	data := []byte(`{
		"event_id": "evt-3",
		"event_timestamp": 1711929600000000000,
		"root_workflow_exec_id": "root-1",
		"parent_workflow_exec_id": null,
		"workflow_exec_id": "exec-1",
		"workflow_run_id": "run-1",
		"workflow_name": "my-workflow",
		"event_type": "ACTIVITY_TASK_RETRYING",
		"attributes": {}
	}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := event.(*ActivityTaskRetryingEvent); !ok {
		t.Fatalf("expected *ActivityTaskRetryingEvent, got %T", event)
	}
}

func TestUnmarshalEvent_UnknownType(t *testing.T) {
	data := []byte(`{
		"event_id": "evt-4",
		"event_timestamp": 1711929600000000000,
		"root_workflow_exec_id": "root-1",
		"parent_workflow_exec_id": null,
		"workflow_exec_id": "exec-1",
		"workflow_run_id": "run-1",
		"workflow_name": "my-workflow",
		"event_type": "FUTURE_EVENT_TYPE",
		"attributes": {}
	}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	unk, ok := event.(*UnknownEvent)
	if !ok {
		t.Fatalf("expected *UnknownEvent, got %T", event)
	}
	if unk.RawType != "FUTURE_EVENT_TYPE" {
		t.Errorf("got RawType %q", unk.RawType)
	}
}

func TestUnmarshalEvent_AllTypes(t *testing.T) {
	types := []struct {
		eventType string
		wantType  string
	}{
		{"WORKFLOW_EXECUTION_STARTED", "*workflow.WorkflowExecutionStartedEvent"},
		{"WORKFLOW_EXECUTION_COMPLETED", "*workflow.WorkflowExecutionCompletedEvent"},
		{"WORKFLOW_EXECUTION_FAILED", "*workflow.WorkflowExecutionFailedEvent"},
		{"WORKFLOW_EXECUTION_CANCELED", "*workflow.WorkflowExecutionCanceledEvent"},
		{"WORKFLOW_EXECUTION_CONTINUED_AS_NEW", "*workflow.WorkflowExecutionContinuedAsNewEvent"},
		{"WORKFLOW_TASK_TIMED_OUT", "*workflow.WorkflowTaskTimedOutEvent"},
		{"WORKFLOW_TASK_FAILED", "*workflow.WorkflowTaskFailedEvent"},
		{"CUSTOM_TASK_STARTED", "*workflow.CustomTaskStartedEvent"},
		{"CUSTOM_TASK_IN_PROGRESS", "*workflow.CustomTaskInProgressEvent"},
		{"CUSTOM_TASK_COMPLETED", "*workflow.CustomTaskCompletedEvent"},
		{"CUSTOM_TASK_FAILED", "*workflow.CustomTaskFailedEvent"},
		{"CUSTOM_TASK_TIMED_OUT", "*workflow.CustomTaskTimedOutEvent"},
		{"CUSTOM_TASK_CANCELED", "*workflow.CustomTaskCanceledEvent"},
		{"ACTIVITY_TASK_STARTED", "*workflow.ActivityTaskStartedEvent"},
		{"ACTIVITY_TASK_COMPLETED", "*workflow.ActivityTaskCompletedEvent"},
		{"ACTIVITY_TASK_RETRYING", "*workflow.ActivityTaskRetryingEvent"},
		{"ACTIVITY_TASK_FAILED", "*workflow.ActivityTaskFailedEvent"},
	}
	for _, tt := range types {
		t.Run(tt.eventType, func(t *testing.T) {
			data, _ := json.Marshal(map[string]any{
				"event_id":                 "evt",
				"event_timestamp":          1711929600000000000,
				"root_workflow_exec_id":    "root",
				"parent_workflow_exec_id":  nil,
				"workflow_exec_id":         "exec",
				"workflow_run_id":          "run",
				"workflow_name":            "wf",
				"event_type":              tt.eventType,
				"attributes":              map[string]any{},
			})
			event, err := UnmarshalEvent(data)
			if err != nil {
				t.Fatal(err)
			}
			got := fmt.Sprintf("%T", event)
			if got != tt.wantType {
				t.Errorf("event_type %q: got %s, want %s", tt.eventType, got, tt.wantType)
			}
		})
	}
}
```

Add `"fmt"` to imports.

- [ ] **Step 2: Run the test to verify it fails**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test ./workflow/ -v -run TestUnmarshalEvent`
Expected: FAIL — `UnmarshalEvent` not defined

- [ ] **Step 3: Implement `workflow/event.go`**

```go
package workflow

import (
	"encoding/json"
	"fmt"
)

// Event is a sealed interface for workflow streaming events.
type Event interface {
	workflowEvent()
	EventType() EventType
}

// EventType identifies the type of workflow event.
type EventType string

const (
	EventWorkflowStarted        EventType = "WORKFLOW_EXECUTION_STARTED"
	EventWorkflowCompleted      EventType = "WORKFLOW_EXECUTION_COMPLETED"
	EventWorkflowFailed         EventType = "WORKFLOW_EXECUTION_FAILED"
	EventWorkflowCanceled       EventType = "WORKFLOW_EXECUTION_CANCELED"
	EventWorkflowContinuedAsNew EventType = "WORKFLOW_EXECUTION_CONTINUED_AS_NEW"
	EventWorkflowTaskTimedOut   EventType = "WORKFLOW_TASK_TIMED_OUT"
	EventWorkflowTaskFailed     EventType = "WORKFLOW_TASK_FAILED"
	EventCustomTaskStarted      EventType = "CUSTOM_TASK_STARTED"
	EventCustomTaskInProgress   EventType = "CUSTOM_TASK_IN_PROGRESS"
	EventCustomTaskCompleted    EventType = "CUSTOM_TASK_COMPLETED"
	EventCustomTaskFailed       EventType = "CUSTOM_TASK_FAILED"
	EventCustomTaskTimedOut     EventType = "CUSTOM_TASK_TIMED_OUT"
	EventCustomTaskCanceled     EventType = "CUSTOM_TASK_CANCELED"
	EventActivityTaskStarted    EventType = "ACTIVITY_TASK_STARTED"
	EventActivityTaskCompleted  EventType = "ACTIVITY_TASK_COMPLETED"
	EventActivityTaskRetrying   EventType = "ACTIVITY_TASK_RETRYING"
	EventActivityTaskFailed     EventType = "ACTIVITY_TASK_FAILED"
)

// EventSource identifies the source of events for streaming.
type EventSource string

const (
	EventSourceDatabase EventSource = "DATABASE"
	EventSourceLive     EventSource = "LIVE"
)

// Scope identifies the scope filter for event streaming.
type Scope string

const (
	ScopeActivity Scope = "activity"
	ScopeWorkflow Scope = "workflow"
	ScopeAll      Scope = "*"
)

// eventBase holds fields common to all workflow events.
type eventBase struct {
	ID                   string  `json:"event_id"`
	Timestamp            int64   `json:"event_timestamp"`
	RootWorkflowExecID   string  `json:"root_workflow_exec_id"`
	ParentWorkflowExecID *string `json:"parent_workflow_exec_id"`
	WorkflowExecID       string  `json:"workflow_exec_id"`
	WorkflowRunID        string  `json:"workflow_run_id"`
	WorkflowName         string  `json:"workflow_name"`
}

// JSONPayload wraps a JSON value with a type discriminator.
type JSONPayload struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

// --- Workflow execution events ---

type WorkflowStartedAttributes struct {
	TaskID string `json:"task_id,omitempty"`
}

type WorkflowExecutionStartedEvent struct {
	eventBase
	Attributes WorkflowStartedAttributes `json:"attributes"`
}

func (*WorkflowExecutionStartedEvent) workflowEvent()       {}
func (*WorkflowExecutionStartedEvent) EventType() EventType { return EventWorkflowStarted }

type WorkflowCompletedAttributes struct {
	TaskID string      `json:"task_id,omitempty"`
	Result JSONPayload `json:"result,omitempty"`
}

type WorkflowExecutionCompletedEvent struct {
	eventBase
	Attributes WorkflowCompletedAttributes `json:"attributes"`
}

func (*WorkflowExecutionCompletedEvent) workflowEvent()       {}
func (*WorkflowExecutionCompletedEvent) EventType() EventType { return EventWorkflowCompleted }

type WorkflowFailedAttributes struct {
	TaskID  string `json:"task_id,omitempty"`
	Failure any    `json:"failure,omitempty"`
}

type WorkflowExecutionFailedEvent struct {
	eventBase
	Attributes WorkflowFailedAttributes `json:"attributes"`
}

func (*WorkflowExecutionFailedEvent) workflowEvent()       {}
func (*WorkflowExecutionFailedEvent) EventType() EventType { return EventWorkflowFailed }

type WorkflowExecutionCanceledEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowExecutionCanceledEvent) workflowEvent()       {}
func (*WorkflowExecutionCanceledEvent) EventType() EventType { return EventWorkflowCanceled }

type WorkflowExecutionContinuedAsNewEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowExecutionContinuedAsNewEvent) workflowEvent()       {}
func (*WorkflowExecutionContinuedAsNewEvent) EventType() EventType { return EventWorkflowContinuedAsNew }

// --- Workflow task events ---

type WorkflowTaskTimedOutEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowTaskTimedOutEvent) workflowEvent()       {}
func (*WorkflowTaskTimedOutEvent) EventType() EventType { return EventWorkflowTaskTimedOut }

type WorkflowTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowTaskFailedEvent) workflowEvent()       {}
func (*WorkflowTaskFailedEvent) EventType() EventType { return EventWorkflowTaskFailed }

// --- Custom task events ---

type CustomTaskStartedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskStartedEvent) workflowEvent()       {}
func (*CustomTaskStartedEvent) EventType() EventType { return EventCustomTaskStarted }

type CustomTaskInProgressEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskInProgressEvent) workflowEvent()       {}
func (*CustomTaskInProgressEvent) EventType() EventType { return EventCustomTaskInProgress }

type CustomTaskCompletedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskCompletedEvent) workflowEvent()       {}
func (*CustomTaskCompletedEvent) EventType() EventType { return EventCustomTaskCompleted }

type CustomTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskFailedEvent) workflowEvent()       {}
func (*CustomTaskFailedEvent) EventType() EventType { return EventCustomTaskFailed }

type CustomTaskTimedOutEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskTimedOutEvent) workflowEvent()       {}
func (*CustomTaskTimedOutEvent) EventType() EventType { return EventCustomTaskTimedOut }

type CustomTaskCanceledEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskCanceledEvent) workflowEvent()       {}
func (*CustomTaskCanceledEvent) EventType() EventType { return EventCustomTaskCanceled }

// --- Activity task events ---

type ActivityTaskStartedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskStartedEvent) workflowEvent()       {}
func (*ActivityTaskStartedEvent) EventType() EventType { return EventActivityTaskStarted }

type ActivityTaskCompletedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskCompletedEvent) workflowEvent()       {}
func (*ActivityTaskCompletedEvent) EventType() EventType { return EventActivityTaskCompleted }

type ActivityTaskRetryingEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskRetryingEvent) workflowEvent()       {}
func (*ActivityTaskRetryingEvent) EventType() EventType { return EventActivityTaskRetrying }

type ActivityTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskFailedEvent) workflowEvent()       {}
func (*ActivityTaskFailedEvent) EventType() EventType { return EventActivityTaskFailed }

// --- Unknown event ---

// UnknownEvent holds a workflow event with an unrecognized event_type.
// This prevents the SDK from breaking when new event types are added.
type UnknownEvent struct {
	eventBase
	RawType string
	Raw     json.RawMessage
}

func (*UnknownEvent) workflowEvent()       {}
func (*UnknownEvent) EventType() EventType { return "UNKNOWN" }

// --- SSE envelope types ---

// StreamPayload is the envelope for workflow SSE events.
type StreamPayload struct {
	Stream          string               `json:"stream"`
	Data            json.RawMessage      `json:"data"`
	WorkflowContext StreamWorkflowContext `json:"workflow_context"`
	BrokerSequence  int                  `json:"broker_sequence"`
	Timestamp       *string              `json:"timestamp,omitempty"`
	Metadata        map[string]any       `json:"metadata,omitempty"`
}

// StreamWorkflowContext provides context about the workflow for a streamed event.
type StreamWorkflowContext struct {
	Namespace            string  `json:"namespace"`
	WorkflowName         string  `json:"workflow_name"`
	WorkflowExecID       string  `json:"workflow_exec_id"`
	ParentWorkflowExecID *string `json:"parent_workflow_exec_id,omitempty"`
	RootWorkflowExecID   *string `json:"root_workflow_exec_id,omitempty"`
}

// EventStreamParams holds query parameters for streaming workflow events.
type EventStreamParams struct {
	Scope                  *Scope
	ActivityName           *string
	ActivityID             *string
	WorkflowName           *string
	WorkflowExecID         *string
	RootWorkflowExecID     *string
	ParentWorkflowExecID   *string
	Stream                 *string
	StartSeq               *int
	MetadataFilters        map[string]any
	WorkflowEventTypes     []EventType
	LastEventID            *string
}

// EventListParams holds query parameters for listing workflow events.
type EventListParams struct {
	RootWorkflowExecID *string
	WorkflowExecID     *string
	WorkflowRunID      *string
	Limit              *int
	Cursor             *string
}

// EventListResponse is the response from listing workflow events.
type EventListResponse struct {
	Events     []json.RawMessage `json:"events"`
	NextCursor *string           `json:"next_cursor,omitempty"`
}

// UnmarshalEvent dispatches JSON to the concrete Event type
// based on the "event_type" discriminator field.
func UnmarshalEvent(data []byte) (Event, error) {
	var probe struct {
		Type string `json:"event_type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal workflow event: %w", err)
	}
	switch EventType(probe.Type) {
	case EventWorkflowStarted:
		var e WorkflowExecutionStartedEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowCompleted:
		var e WorkflowExecutionCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowFailed:
		var e WorkflowExecutionFailedEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowCanceled:
		var e WorkflowExecutionCanceledEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowContinuedAsNew:
		var e WorkflowExecutionContinuedAsNewEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowTaskTimedOut:
		var e WorkflowTaskTimedOutEvent
		return &e, json.Unmarshal(data, &e)
	case EventWorkflowTaskFailed:
		var e WorkflowTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskStarted:
		var e CustomTaskStartedEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskInProgress:
		var e CustomTaskInProgressEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskCompleted:
		var e CustomTaskCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskFailed:
		var e CustomTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskTimedOut:
		var e CustomTaskTimedOutEvent
		return &e, json.Unmarshal(data, &e)
	case EventCustomTaskCanceled:
		var e CustomTaskCanceledEvent
		return &e, json.Unmarshal(data, &e)
	case EventActivityTaskStarted:
		var e ActivityTaskStartedEvent
		return &e, json.Unmarshal(data, &e)
	case EventActivityTaskCompleted:
		var e ActivityTaskCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case EventActivityTaskRetrying:
		var e ActivityTaskRetryingEvent
		return &e, json.Unmarshal(data, &e)
	case EventActivityTaskFailed:
		var e ActivityTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	default:
		return &UnknownEvent{
			RawType: probe.Type,
			Raw:     json.RawMessage(data),
		}, nil
	}
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test ./workflow/ -v -run TestUnmarshalEvent`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add workflow/event.go workflow/event_test.go
git commit -m "feat(workflow): add sealed Event interface with 17 types and UnmarshalEvent"
```

---

## Task 4: Remaining type files

**Files:**
- Create: `workflow/deployment.go`
- Create: `workflow/metrics.go`
- Create: `workflow/run.go`
- Create: `workflow/schedule.go`
- Create: `workflow/registration.go`

- [ ] **Step 1: Create `workflow/deployment.go`**

```go
package workflow

// Deployment represents a workflow deployment.
type Deployment struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// DeploymentListResponse is the response from listing deployments.
type DeploymentListResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// DeploymentListParams holds query parameters for listing deployments.
type DeploymentListParams struct {
	ActiveOnly   *bool
	WorkflowName *string
}
```

- [ ] **Step 2: Create `workflow/metrics.go`**

```go
package workflow

// Metrics holds workflow performance metrics.
type Metrics struct {
	ExecutionCount   ScalarMetric     `json:"execution_count"`
	SuccessCount     ScalarMetric     `json:"success_count"`
	ErrorCount       ScalarMetric     `json:"error_count"`
	AverageLatencyMs ScalarMetric     `json:"average_latency_ms"`
	LatencyOverTime  TimeSeriesMetric `json:"latency_over_time"`
	RetryRate        ScalarMetric     `json:"retry_rate"`
}

// ScalarMetric holds a single numeric metric value.
type ScalarMetric struct {
	Value float64 `json:"value"`
}

// TimeSeriesMetric holds a time series of [timestamp, value] pairs.
type TimeSeriesMetric struct {
	Value [][]float64 `json:"value"`
}

// MetricsParams holds query parameters for workflow metrics.
type MetricsParams struct {
	StartTime *string
	EndTime   *string
}
```

- [ ] **Step 3: Create `workflow/run.go`**

```go
package workflow

// Run represents a workflow run.
type Run struct {
	ID           string          `json:"id"`
	WorkflowName string          `json:"workflow_name"`
	ExecutionID  string          `json:"execution_id"`
	Status       ExecutionStatus `json:"status"`
	StartTime    string          `json:"start_time"`
	EndTime      *string         `json:"end_time,omitempty"`
}

// ListRunsResponse is the response from listing workflow runs.
type ListRunsResponse struct {
	Runs          []Run   `json:"runs"`
	NextPageToken *string `json:"next_page_token,omitempty"`
}

// RunListParams holds query parameters for listing workflow runs.
type RunListParams struct {
	WorkflowIdentifier *string
	Search             *string
	Status             *string
	PageSize           *int
	NextPageToken      *string
}
```

- [ ] **Step 4: Create `workflow/schedule.go`**

```go
package workflow

// ScheduleRequest is the request body for scheduling a workflow.
type ScheduleRequest struct {
	Schedule               ScheduleDefinition `json:"schedule"`
	WorkflowRegistrationID *string            `json:"workflow_registration_id,omitempty"`
	WorkflowIdentifier     *string            `json:"workflow_identifier,omitempty"`
	ScheduleID             *string            `json:"schedule_id,omitempty"`
	DeploymentName         *string            `json:"deployment_name,omitempty"`
}

// ScheduleDefinition describes when and how a workflow should be scheduled.
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

// ScheduleCalendar defines calendar-based schedule entries.
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

// ScheduleRange defines a numeric range for calendar schedules.
type ScheduleRange struct {
	Start int `json:"start"`
	End   int `json:"end,omitempty"`
	Step  int `json:"step,omitempty"`
}

// ScheduleInterval defines an interval-based schedule.
type ScheduleInterval struct {
	Every  string  `json:"every"`
	Offset *string `json:"offset,omitempty"`
}

// SchedulePolicy controls schedule overlap and failure behavior.
type SchedulePolicy struct {
	CatchupWindowSeconds int  `json:"catchup_window_seconds,omitempty"`
	Overlap              *int `json:"overlap,omitempty"`
	PauseOnFailure       bool `json:"pause_on_failure,omitempty"`
}

// ScheduleResponse is the response from creating a workflow schedule.
type ScheduleResponse struct {
	ScheduleID string `json:"schedule_id"`
}

// ScheduleListResponse is the response from listing workflow schedules.
type ScheduleListResponse struct {
	Schedules []Schedule `json:"schedules"`
}

// Schedule represents a workflow schedule.
type Schedule struct {
	ScheduleID   string             `json:"schedule_id"`
	Definition   ScheduleDefinition `json:"definition"`
	WorkflowName string             `json:"workflow_name"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}
```

- [ ] **Step 5: Create `workflow/registration.go`**

```go
package workflow

// Registration represents a workflow registration.
type Registration struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflow_id"`
	TaskQueue  string    `json:"task_queue"`
	Workflow   *Workflow `json:"workflow,omitempty"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

// RegistrationListResponse is the response from listing workflow registrations.
type RegistrationListResponse struct {
	Registrations []Registration `json:"registrations"`
	NextCursor    *string        `json:"next_cursor,omitempty"`
}

// RegistrationListParams holds query parameters for listing registrations.
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

// RegistrationGetParams holds query parameters for getting a registration.
type RegistrationGetParams struct {
	WithWorkflow  *bool
	IncludeShared *bool
}

// WorkerInfo holds information about the current worker.
type WorkerInfo struct {
	SchedulerURL string `json:"scheduler_url"`
	Namespace    string `json:"namespace"`
	TLS          bool   `json:"tls"`
}
```

- [ ] **Step 6: Verify everything compiles**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go build ./workflow/ && go test ./workflow/ -v`
Expected: build succeeds, all tests pass

- [ ] **Step 7: Commit**

```bash
git add workflow/deployment.go workflow/metrics.go workflow/run.go workflow/schedule.go workflow/registration.go
git commit -m "feat(workflow): add deployment, metrics, run, schedule, and registration types"
```

---

## Task 5: Workflows CRUD service methods + tests

**Files:**
- Create: `workflows.go`
- Create: `workflows_test.go`

- [ ] **Step 1: Write failing tests in `workflows_test.go`**

```go
package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestListWorkflows_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("got limit %q", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("active_only") != "true" {
			t.Errorf("got active_only %q", r.URL.Query().Get("active_only"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"workflows": []map[string]any{
				{"id": "wf-1", "name": "my-flow", "owner_id": "u1", "workspace_id": "ws1", "created_at": "2026-01-01", "updated_at": "2026-01-01"},
			},
			"next_cursor": "cur-abc",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	active := true
	limit := 10
	resp, err := client.ListWorkflows(context.Background(), &workflow.WorkflowListParams{
		ActiveOnly: &active,
		Limit:      &limit,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Workflows) != 1 {
		t.Fatalf("got %d workflows", len(resp.Workflows))
	}
	if resp.Workflows[0].ID != "wf-1" {
		t.Errorf("got id %q", resp.Workflows[0].ID)
	}
	if resp.NextCursor == nil || *resp.NextCursor != "cur-abc" {
		t.Errorf("got cursor %v", resp.NextCursor)
	}
}

func TestGetWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/wf-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "wf-1", "name": "my-flow", "owner_id": "u1", "workspace_id": "ws1",
			"created_at": "2026-01-01", "updated_at": "2026-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	wf, err := client.GetWorkflow(context.Background(), "wf-1")
	if err != nil {
		t.Fatal(err)
	}
	if wf.Name != "my-flow" {
		t.Errorf("got name %q", wf.Name)
	}
}

func TestUpdateWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("got method %s", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["display_name"] != "New Name" {
			t.Errorf("got display_name %v", body["display_name"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "wf-1", "name": "my-flow", "display_name": "New Name",
			"owner_id": "u1", "workspace_id": "ws1",
			"created_at": "2026-01-01", "updated_at": "2026-01-02",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	name := "New Name"
	wf, err := client.UpdateWorkflow(context.Background(), "wf-1", &workflow.WorkflowUpdateRequest{
		DisplayName: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if wf.DisplayName == nil || *wf.DisplayName != "New Name" {
		t.Errorf("got display_name %v", wf.DisplayName)
	}
}

func TestArchiveWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/wf-1/archive" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"id": "wf-1", "archived": true})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ArchiveWorkflow(context.Background(), "wf-1")
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Archived {
		t.Error("expected archived=true")
	}
}

func TestExecuteWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/wf-1/execute" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		input, _ := body["input"].(map[string]any)
		if input["prompt"] != "hello" {
			t.Errorf("got input %v", body["input"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"workflow_name": "my-flow", "execution_id": "exec-1",
			"root_execution_id": "exec-1", "status": "RUNNING",
			"start_time": "2026-01-01T00:00:00Z",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ExecuteWorkflow(context.Background(), "wf-1", &workflow.ExecutionRequest{
		Input: map[string]any{"prompt": "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExecutionID != "exec-1" {
		t.Errorf("got execution_id %q", resp.ExecutionID)
	}
	if resp.Status != workflow.ExecutionRunning {
		t.Errorf("got status %q", resp.Status)
	}
}

func TestExecuteWorkflowAndWait_Success(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && r.URL.Path == "/v1/workflows/wf-1/execute":
			json.NewEncoder(w).Encode(map[string]any{
				"workflow_name": "my-flow", "execution_id": "exec-1",
				"root_execution_id": "exec-1", "status": "RUNNING",
				"start_time": "2026-01-01T00:00:00Z",
			})
		case r.Method == "GET" && r.URL.Path == "/v1/workflows/executions/exec-1":
			calls++
			status := "RUNNING"
			if calls >= 2 {
				status = "COMPLETED"
			}
			resp := map[string]any{
				"workflow_name": "my-flow", "execution_id": "exec-1",
				"root_execution_id": "exec-1", "status": status,
				"start_time": "2026-01-01T00:00:00Z",
			}
			if status == "COMPLETED" {
				resp["result"] = map[string]any{"answer": 42}
			}
			json.NewEncoder(w).Encode(resp)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ExecuteWorkflowAndWait(context.Background(), "wf-1", &workflow.ExecutionRequest{
		Input: map[string]any{"prompt": "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != workflow.ExecutionCompleted {
		t.Errorf("got status %q", resp.Status)
	}
}

func TestListWorkflowRegistrations_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/registrations" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"registrations": []map[string]any{
				{"id": "reg-1", "workflow_id": "wf-1", "task_queue": "default", "created_at": "2026-01-01", "updated_at": "2026-01-01"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListWorkflowRegistrations(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Registrations) != 1 {
		t.Fatalf("got %d registrations", len(resp.Registrations))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflows|TestGetWorkflow|TestUpdateWorkflow|TestArchiveWorkflow|TestExecuteWorkflow|TestListWorkflowRegistrations" -v`
Expected: FAIL — methods not defined

- [ ] **Step 3: Implement `workflows.go`**

```go
package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

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
```

Note: `ExecuteWorkflowAndWait` calls `GetWorkflowExecution` which is defined in Task 6. The compiler won't complain during test runs because both files are in the same package, but **Task 6 must be implemented before `TestExecuteWorkflowAndWait_Success` can pass**. Run only the other tests first.

- [ ] **Step 4: Run tests (excluding ExecuteWorkflowAndWait)**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflows_Success|TestGetWorkflow_Success|TestUpdateWorkflow_Success|TestArchiveWorkflow_Success|TestExecuteWorkflow_Success|TestListWorkflowRegistrations_Success" -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add workflows.go workflows_test.go
git commit -m "feat: add workflows CRUD and registration service methods"
```

---

## Task 6: Execution service methods + WorkflowEventStream + tests

**Files:**
- Create: `workflows_executions.go`
- Create: `workflows_executions_test.go`

- [ ] **Step 1: Write failing tests in `workflows_executions_test.go`**

```go
package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestGetWorkflowExecution_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/executions/exec-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"workflow_name": "my-flow", "execution_id": "exec-1",
			"root_execution_id": "exec-1", "status": "COMPLETED",
			"start_time": "2026-01-01T00:00:00Z",
			"end_time":   "2026-01-01T00:01:00Z",
			"result":     map[string]any{"answer": 42},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetWorkflowExecution(context.Background(), "exec-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != workflow.ExecutionCompleted {
		t.Errorf("got status %q", resp.Status)
	}
}

func TestSignalWorkflowExecution_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/executions/exec-1/signals" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "approval" {
			t.Errorf("got name %v", body["name"])
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]any{"message": "Signal accepted"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.SignalWorkflowExecution(context.Background(), "exec-1", &workflow.SignalInvocationBody{
		Name:  "approval",
		Input: map[string]any{"approved": true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Signal accepted" {
		t.Errorf("got message %q", resp.Message)
	}
}

func TestTerminateWorkflowExecution_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/executions/exec-1/terminate" {
			t.Errorf("got path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.TerminateWorkflowExecution(context.Background(), "exec-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBatchCancelWorkflowExecutions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/executions/cancel" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"results": map[string]any{
				"exec-1": map[string]any{"status": "success"},
				"exec-2": map[string]any{"status": "failure", "error": "not found"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.BatchCancelWorkflowExecutions(context.Background(), &workflow.BatchExecutionBody{
		ExecutionIDs: []string{"exec-1", "exec-2"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Results["exec-1"].Status != "success" {
		t.Errorf("got exec-1 status %q", resp.Results["exec-1"].Status)
	}
	if resp.Results["exec-2"].Error == nil || *resp.Results["exec-2"].Error != "not found" {
		t.Errorf("expected exec-2 error")
	}
}

func TestStreamWorkflowExecution_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/executions/exec-1/stream" {
			t.Errorf("got path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		payloads := []map[string]any{
			{
				"stream": "events",
				"data": map[string]any{
					"event_id": "evt-1", "event_timestamp": 1711929600000000000,
					"root_workflow_exec_id": "exec-1", "parent_workflow_exec_id": nil,
					"workflow_exec_id": "exec-1", "workflow_run_id": "run-1",
					"workflow_name": "my-flow", "event_type": "WORKFLOW_EXECUTION_STARTED",
					"attributes": map[string]any{},
				},
				"workflow_context": map[string]any{
					"namespace": "default", "workflow_name": "my-flow", "workflow_exec_id": "exec-1",
				},
				"broker_sequence": 1,
			},
			{
				"stream": "events",
				"data": map[string]any{
					"event_id": "evt-2", "event_timestamp": 1711929601000000000,
					"root_workflow_exec_id": "exec-1", "parent_workflow_exec_id": nil,
					"workflow_exec_id": "exec-1", "workflow_run_id": "run-1",
					"workflow_name": "my-flow", "event_type": "WORKFLOW_EXECUTION_COMPLETED",
					"attributes": map[string]any{"result": map[string]any{"value": 42, "type": "json"}},
				},
				"workflow_context": map[string]any{
					"namespace": "default", "workflow_name": "my-flow", "workflow_exec_id": "exec-1",
				},
				"broker_sequence": 2,
			},
		}
		for _, p := range payloads {
			data, _ := json.Marshal(p)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.StreamWorkflowExecution(context.Background(), "exec-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var events []workflow.Event
	for stream.Next() {
		events = append(events, stream.Current())
	}
	if stream.Err() != nil {
		t.Fatal(stream.Err())
	}
	if len(events) != 2 {
		t.Fatalf("got %d events, want 2", len(events))
	}
	if _, ok := events[0].(*workflow.WorkflowExecutionStartedEvent); !ok {
		t.Errorf("expected *WorkflowExecutionStartedEvent, got %T", events[0])
	}
	if _, ok := events[1].(*workflow.WorkflowExecutionCompletedEvent); !ok {
		t.Errorf("expected *WorkflowExecutionCompletedEvent, got %T", events[1])
	}

	// Verify payload metadata is accessible
	payload := stream.CurrentPayload()
	if payload.WorkflowContext.WorkflowName != "my-flow" {
		t.Errorf("got workflow context name %q", payload.WorkflowContext.WorkflowName)
	}
}

func TestGetWorkflowExecutionTraceOTel_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/executions/exec-1/trace/otel" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"workflow_name": "my-flow", "execution_id": "exec-1",
			"root_execution_id": "exec-1", "status": "COMPLETED",
			"start_time": "2026-01-01T00:00:00Z", "data_source": "temporal",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetWorkflowExecutionTraceOTel(context.Background(), "exec-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.DataSource != "temporal" {
		t.Errorf("got data_source %q", resp.DataSource)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestGetWorkflowExecution_Success|TestSignalWorkflow|TestTerminateWorkflow|TestBatchCancel|TestStreamWorkflowExecution|TestGetWorkflowExecutionTrace" -v`
Expected: FAIL — methods not defined

- [ ] **Step 3: Implement `workflows_executions.go`**

```go
package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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
```

- [ ] **Step 4: Run all tests including ExecuteWorkflowAndWait**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestGetWorkflowExecution|TestSignalWorkflow|TestTerminateWorkflow|TestBatchCancel|TestStreamWorkflowExecution|TestGetWorkflowExecutionTrace|TestExecuteWorkflowAndWait" -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add workflows_executions.go workflows_executions_test.go
git commit -m "feat: add workflow execution service methods and WorkflowEventStream"
```

---

## Task 7: Events, Deployments, Metrics service methods + tests

**Files:**
- Create: `workflows_events.go`, `workflows_events_test.go`
- Create: `workflows_deployments.go`, `workflows_deployments_test.go`
- Create: `workflows_metrics.go`, `workflows_metrics_test.go`

- [ ] **Step 1: Write failing tests for all three files**

`workflows_events_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestListWorkflowEvents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/events/list" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("got limit %q", r.URL.Query().Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"events":      []map[string]any{{"event_type": "WORKFLOW_EXECUTION_STARTED"}},
			"next_cursor": "cur-1",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	limit := 50
	resp, err := client.ListWorkflowEvents(context.Background(), &workflow.EventListParams{Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("got %d events", len(resp.Events))
	}
	if resp.NextCursor == nil || *resp.NextCursor != "cur-1" {
		t.Errorf("got cursor %v", resp.NextCursor)
	}
}
```

`workflows_deployments_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkflowDeployments_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/deployments" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"deployments": []map[string]any{
				{"id": "dep-1", "name": "prod", "is_active": true, "created_at": "2026-01-01", "updated_at": "2026-01-01"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListWorkflowDeployments(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Deployments) != 1 {
		t.Fatalf("got %d deployments", len(resp.Deployments))
	}
	if resp.Deployments[0].Name != "prod" {
		t.Errorf("got name %q", resp.Deployments[0].Name)
	}
}

func TestGetWorkflowDeployment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/deployments/dep-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "dep-1", "name": "prod", "is_active": true,
			"created_at": "2026-01-01", "updated_at": "2026-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	dep, err := client.GetWorkflowDeployment(context.Background(), "dep-1")
	if err != nil {
		t.Fatal(err)
	}
	if dep.ID != "dep-1" {
		t.Errorf("got id %q", dep.ID)
	}
}
```

`workflows_metrics_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkflowMetrics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/my-flow/metrics" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_time") != "2026-01-01T00:00:00Z" {
			t.Errorf("got start_time %q", r.URL.Query().Get("start_time"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"execution_count":    map[string]any{"value": 100},
			"success_count":      map[string]any{"value": 95},
			"error_count":        map[string]any{"value": 5},
			"average_latency_ms": map[string]any{"value": 1234.5},
			"latency_over_time":  map[string]any{"value": [][]float64{{1711929600, 1200}, {1711929660, 1300}}},
			"retry_rate":         map[string]any{"value": 0.02},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	start := "2026-01-01T00:00:00Z"
	resp, err := client.GetWorkflowMetrics(context.Background(), "my-flow", &workflow.MetricsParams{StartTime: &start})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExecutionCount.Value != 100 {
		t.Errorf("got execution_count %v", resp.ExecutionCount.Value)
	}
	if resp.AverageLatencyMs.Value != 1234.5 {
		t.Errorf("got average_latency_ms %v", resp.AverageLatencyMs.Value)
	}
}
```

Add `"github.com/VikingOwl91/mistral-go-sdk/workflow"` import to `workflows_metrics_test.go`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflowEvents|TestListWorkflowDeployments|TestGetWorkflowDeployment|TestGetWorkflowMetrics" -v`
Expected: FAIL

- [ ] **Step 3: Implement all three service files**

`workflows_events.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
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
```

`workflows_deployments.go`:

```go
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
```

`workflows_metrics.go`:

```go
package mistral

import (
	"context"
	"fmt"
	"net/url"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
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
```

- [ ] **Step 4: Run all tests**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflowEvents|TestListWorkflowDeployments|TestGetWorkflowDeployment|TestGetWorkflowMetrics" -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add workflows_events.go workflows_events_test.go workflows_deployments.go workflows_deployments_test.go workflows_metrics.go workflows_metrics_test.go
git commit -m "feat: add workflow events, deployments, and metrics service methods"
```

---

## Task 8: Runs, Schedules, Workers service methods + tests

**Files:**
- Create: `workflows_runs.go`, `workflows_runs_test.go`
- Create: `workflows_schedules.go`, `workflows_schedules_test.go`
- Create: `workflows_workers.go`, `workflows_workers_test.go`

- [ ] **Step 1: Write failing tests**

`workflows_runs_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkflowRuns_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/runs" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"runs": []map[string]any{
				{"id": "run-1", "workflow_name": "my-flow", "execution_id": "exec-1", "status": "COMPLETED", "start_time": "2026-01-01"},
			},
			"next_page_token": "tok-1",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListWorkflowRuns(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Runs) != 1 {
		t.Fatalf("got %d runs", len(resp.Runs))
	}
	if resp.NextPageToken == nil || *resp.NextPageToken != "tok-1" {
		t.Errorf("got token %v", resp.NextPageToken)
	}
}

func TestGetWorkflowRun_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/runs/run-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "run-1", "workflow_name": "my-flow", "execution_id": "exec-1",
			"status": "COMPLETED", "start_time": "2026-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	run, err := client.GetWorkflowRun(context.Background(), "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if run.ID != "run-1" {
		t.Errorf("got id %q", run.ID)
	}
}
```

`workflows_schedules_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestScheduleWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/schedules" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		schedule, _ := body["schedule"].(map[string]any)
		cronExprs, _ := schedule["cron_expressions"].([]any)
		if len(cronExprs) != 1 || cronExprs[0] != "0 9 * * MON-FRI" {
			t.Errorf("got cron_expressions %v", cronExprs)
		}
		json.NewEncoder(w).Encode(map[string]any{"schedule_id": "sched-1"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	wfID := "wf-1"
	resp, err := client.ScheduleWorkflow(context.Background(), &workflow.ScheduleRequest{
		WorkflowIdentifier: &wfID,
		Schedule: workflow.ScheduleDefinition{
			Input:           map[string]any{"prompt": "daily report"},
			CronExpressions: []string{"0 9 * * MON-FRI"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ScheduleID != "sched-1" {
		t.Errorf("got schedule_id %q", resp.ScheduleID)
	}
}

func TestUnscheduleWorkflow_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/workflows/schedules/sched-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.UnscheduleWorkflow(context.Background(), "sched-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListWorkflowSchedules_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/schedules" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"schedules": []map[string]any{
				{"schedule_id": "sched-1", "workflow_name": "my-flow", "created_at": "2026-01-01", "updated_at": "2026-01-01"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListWorkflowSchedules(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Schedules) != 1 {
		t.Fatalf("got %d schedules", len(resp.Schedules))
	}
}
```

`workflows_workers_test.go`:

```go
package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkflowWorkerInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/workers/whoami" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"scheduler_url": "https://scheduler.mistral.ai",
			"namespace":     "default",
			"tls":           true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	info, err := client.GetWorkflowWorkerInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Namespace != "default" {
		t.Errorf("got namespace %q", info.Namespace)
	}
	if !info.TLS {
		t.Error("expected tls=true")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflowRuns|TestGetWorkflowRun|TestScheduleWorkflow|TestUnscheduleWorkflow|TestListWorkflowSchedules|TestGetWorkflowWorkerInfo" -v`
Expected: FAIL

- [ ] **Step 3: Implement all three service files**

`workflows_runs.go`:

```go
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
```

`workflows_schedules.go`:

```go
package mistral

import (
	"context"
	"fmt"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// ListWorkflowSchedules lists workflow schedules.
func (c *Client) ListWorkflowSchedules(ctx context.Context) (*workflow.ScheduleListResponse, error) {
	var resp workflow.ScheduleListResponse
	if err := c.doJSON(ctx, "GET", "/v1/workflows/schedules", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ScheduleWorkflow creates a workflow schedule.
func (c *Client) ScheduleWorkflow(ctx context.Context, req *workflow.ScheduleRequest) (*workflow.ScheduleResponse, error) {
	var resp workflow.ScheduleResponse
	if err := c.doJSON(ctx, "POST", "/v1/workflows/schedules", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnscheduleWorkflow removes a workflow schedule.
func (c *Client) UnscheduleWorkflow(ctx context.Context, scheduleID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/workflows/schedules/%s", scheduleID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}
```

`workflows_workers.go`:

```go
package mistral

import (
	"context"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// GetWorkflowWorkerInfo retrieves information about the current worker.
func (c *Client) GetWorkflowWorkerInfo(ctx context.Context) (*workflow.WorkerInfo, error) {
	var resp workflow.WorkerInfo
	if err := c.doJSON(ctx, "GET", "/v1/workflows/workers/whoami", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
```

- [ ] **Step 4: Run all tests**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run "TestListWorkflowRuns|TestGetWorkflowRun|TestScheduleWorkflow|TestUnscheduleWorkflow|TestListWorkflowSchedules|TestGetWorkflowWorkerInfo" -v`
Expected: all PASS

- [ ] **Step 5: Commit**

```bash
git add workflows_runs.go workflows_runs_test.go workflows_schedules.go workflows_schedules_test.go workflows_workers.go workflows_workers_test.go
git commit -m "feat: add workflow runs, schedules, and workers service methods"
```

---

## Task 9: DeleteBatchJob

**Files:**
- Modify: `batch_api.go`
- Modify or create type in `batch/` package
- Modify: `batch_test.go`

- [ ] **Step 1: Add `DeleteResponse` to batch package**

Find the batch types file (likely `batch/batch.go` or similar — check first) and add:

```go
// DeleteResponse is the response from deleting a batch job.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
```

- [ ] **Step 2: Write failing test in `batch_test.go`**

Append to existing `batch_test.go`:

```go
func TestDeleteBatchJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/batch/jobs/batch-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "batch-123", "object": "batch", "deleted": true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.DeleteBatchJob(context.Background(), "batch-123")
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "batch-123" {
		t.Errorf("got id %q", resp.ID)
	}
	if !resp.Deleted {
		t.Error("expected deleted=true")
	}
}
```

- [ ] **Step 3: Run test to verify it fails**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run TestDeleteBatchJob_Success -v`
Expected: FAIL

- [ ] **Step 4: Add `DeleteBatchJob` to `batch_api.go`**

Append to `batch_api.go`:

```go
// DeleteBatchJob deletes a batch job.
func (c *Client) DeleteBatchJob(ctx context.Context, jobID string) (*batch.DeleteResponse, error) {
	var resp batch.DeleteResponse
	if err := c.doJSON(ctx, "DELETE", fmt.Sprintf("/v1/batch/jobs/%s", jobID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go test -run TestDeleteBatchJob_Success -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add batch/ batch_api.go batch_test.go
git commit -m "feat(batch): add DeleteBatchJob method"
```

---

## Task 10: Version bump, changelog, docs, final verification

**Files:**
- Modify: `mistral.go`
- Modify: `CHANGELOG.md`
- Modify: `CLAUDE.md`

- [ ] **Step 1: Bump version in `mistral.go`**

Change `const Version = "1.1.0"` to `const Version = "1.2.0"`

- [ ] **Step 2: Update `CHANGELOG.md`**

Prepend to the file:

```markdown
## v1.2.0 — 2026-04-02

Upstream sync with Python SDK v2.2.0. Adds Workflows API and DeleteBatchJob.

### Added

- **Workflows API** (new `workflow/` package) — complete workflow orchestration
  support with 37 service methods across 8 sub-resources:
  - **Workflows CRUD** — `ListWorkflows`, `GetWorkflow`, `UpdateWorkflow`,
    `ArchiveWorkflow`, `UnarchiveWorkflow`, `ExecuteWorkflow`,
    `ExecuteWorkflowAndWait`.
  - **Registrations** — `ListWorkflowRegistrations`, `GetWorkflowRegistration`,
    `ExecuteWorkflowRegistration` (deprecated).
  - **Executions** — `GetWorkflowExecution`, `GetWorkflowExecutionHistory`,
    `StreamWorkflowExecution`, `SignalWorkflowExecution`,
    `QueryWorkflowExecution`, `UpdateWorkflowExecution`,
    `TerminateWorkflowExecution`, `CancelWorkflowExecution`,
    `ResetWorkflowExecution`, `BatchCancelWorkflowExecutions`,
    `BatchTerminateWorkflowExecutions`.
  - **Trace** — `GetWorkflowExecutionTraceOTel`,
    `GetWorkflowExecutionTraceSummary`, `GetWorkflowExecutionTraceEvents`.
  - **Events** — `StreamWorkflowEvents`, `ListWorkflowEvents`.
  - **Deployments** — `ListWorkflowDeployments`, `GetWorkflowDeployment`.
  - **Metrics** — `GetWorkflowMetrics`.
  - **Runs** — `ListWorkflowRuns`, `GetWorkflowRun`, `GetWorkflowRunHistory`.
  - **Schedules** — `ListWorkflowSchedules`, `ScheduleWorkflow`,
    `UnscheduleWorkflow`.
  - **Workers** — `GetWorkflowWorkerInfo`.
- **`WorkflowEventStream`** — typed SSE stream wrapper with `StreamPayload`
  envelope, sealed `Event` interface (17 concrete types + `UnknownEvent`).
- **`DeleteBatchJob`** — delete a batch job by ID.

```

- [ ] **Step 3: Update `CLAUDE.md` sub-packages list**

In the "Two-layer design" section, add `workflow/` to the list of sub-packages.

- [ ] **Step 4: Run full test suite**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && go vet ./... && go test ./... -v`
Expected: all pass, no vet issues

- [ ] **Step 5: Verify method count**

Run: `cd /home/cnachtigall/ssd/git/active/mistral-go-sdk/mistral-go-sdk && grep -c 'func (c \*Client)' workflows*.go batch_api.go`
Expected: total of 38 new methods (37 workflow + 1 batch)

- [ ] **Step 6: Commit**

```bash
git add mistral.go CHANGELOG.md CLAUDE.md
git commit -m "chore: bump version to v1.2.0, update changelog and docs"
```
