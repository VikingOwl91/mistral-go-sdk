package workflow

import (
	"encoding/json"
	"fmt"
)

// EventType identifies the kind of workflow event.
type EventType string

const (
	EventWorkflowStarted         EventType = "WORKFLOW_EXECUTION_STARTED"
	EventWorkflowCompleted       EventType = "WORKFLOW_EXECUTION_COMPLETED"
	EventWorkflowFailed          EventType = "WORKFLOW_EXECUTION_FAILED"
	EventWorkflowCanceled        EventType = "WORKFLOW_EXECUTION_CANCELED"
	EventWorkflowContinuedAsNew  EventType = "WORKFLOW_EXECUTION_CONTINUED_AS_NEW"
	EventWorkflowTaskTimedOut    EventType = "WORKFLOW_TASK_TIMED_OUT"
	EventWorkflowTaskFailed      EventType = "WORKFLOW_TASK_FAILED"
	EventCustomTaskStarted       EventType = "CUSTOM_TASK_STARTED"
	EventCustomTaskInProgress    EventType = "CUSTOM_TASK_IN_PROGRESS"
	EventCustomTaskCompleted     EventType = "CUSTOM_TASK_COMPLETED"
	EventCustomTaskFailed        EventType = "CUSTOM_TASK_FAILED"
	EventCustomTaskTimedOut      EventType = "CUSTOM_TASK_TIMED_OUT"
	EventCustomTaskCanceled      EventType = "CUSTOM_TASK_CANCELED"
	EventActivityTaskStarted     EventType = "ACTIVITY_TASK_STARTED"
	EventActivityTaskCompleted   EventType = "ACTIVITY_TASK_COMPLETED"
	EventActivityTaskRetrying    EventType = "ACTIVITY_TASK_RETRYING"
	EventActivityTaskFailed      EventType = "ACTIVITY_TASK_FAILED"
)

// EventSource identifies where an event originated.
type EventSource string

const (
	EventSourceDatabase EventSource = "DATABASE"
	EventSourceLive     EventSource = "LIVE"
)

// Scope identifies the scope of an event subscription.
type Scope string

const (
	ScopeActivity Scope = "activity"
	ScopeWorkflow Scope = "workflow"
	ScopeAll      Scope = "*"
)

// Event is a sealed interface for workflow execution events.
type Event interface {
	workflowEvent()
	EventType() EventType
}

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

// WorkflowStartedAttributes holds typed attributes for workflow started events.
type WorkflowStartedAttributes struct {
	TaskID string `json:"task_id"`
}

// JSONPayload holds a typed JSON value.
type JSONPayload struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}

// WorkflowCompletedAttributes holds typed attributes for workflow completed events.
type WorkflowCompletedAttributes struct {
	TaskID string      `json:"task_id"`
	Result JSONPayload `json:"result"`
}

// WorkflowFailedAttributes holds typed attributes for workflow failed events.
type WorkflowFailedAttributes struct {
	TaskID  string `json:"task_id"`
	Failure any    `json:"failure"`
}

// WorkflowExecutionStartedEvent signals that a workflow execution has started.
type WorkflowExecutionStartedEvent struct {
	eventBase
	Attributes WorkflowStartedAttributes `json:"attributes"`
}

func (*WorkflowExecutionStartedEvent) workflowEvent()      {}
func (*WorkflowExecutionStartedEvent) EventType() EventType { return EventWorkflowStarted }

// WorkflowExecutionCompletedEvent signals that a workflow execution has completed.
type WorkflowExecutionCompletedEvent struct {
	eventBase
	Attributes WorkflowCompletedAttributes `json:"attributes"`
}

func (*WorkflowExecutionCompletedEvent) workflowEvent()      {}
func (*WorkflowExecutionCompletedEvent) EventType() EventType { return EventWorkflowCompleted }

// WorkflowExecutionFailedEvent signals that a workflow execution has failed.
type WorkflowExecutionFailedEvent struct {
	eventBase
	Attributes WorkflowFailedAttributes `json:"attributes"`
}

func (*WorkflowExecutionFailedEvent) workflowEvent()      {}
func (*WorkflowExecutionFailedEvent) EventType() EventType { return EventWorkflowFailed }

// WorkflowExecutionCanceledEvent signals that a workflow execution was canceled.
type WorkflowExecutionCanceledEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowExecutionCanceledEvent) workflowEvent()      {}
func (*WorkflowExecutionCanceledEvent) EventType() EventType { return EventWorkflowCanceled }

// WorkflowExecutionContinuedAsNewEvent signals that a workflow continued as a new execution.
type WorkflowExecutionContinuedAsNewEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowExecutionContinuedAsNewEvent) workflowEvent()      {}
func (*WorkflowExecutionContinuedAsNewEvent) EventType() EventType { return EventWorkflowContinuedAsNew }

// WorkflowTaskTimedOutEvent signals that a workflow task timed out.
type WorkflowTaskTimedOutEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowTaskTimedOutEvent) workflowEvent()      {}
func (*WorkflowTaskTimedOutEvent) EventType() EventType { return EventWorkflowTaskTimedOut }

// WorkflowTaskFailedEvent signals that a workflow task failed.
type WorkflowTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*WorkflowTaskFailedEvent) workflowEvent()      {}
func (*WorkflowTaskFailedEvent) EventType() EventType { return EventWorkflowTaskFailed }

// CustomTaskStartedEvent signals that a custom task has started.
type CustomTaskStartedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskStartedEvent) workflowEvent()      {}
func (*CustomTaskStartedEvent) EventType() EventType { return EventCustomTaskStarted }

// CustomTaskInProgressEvent signals that a custom task is in progress.
type CustomTaskInProgressEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskInProgressEvent) workflowEvent()      {}
func (*CustomTaskInProgressEvent) EventType() EventType { return EventCustomTaskInProgress }

// CustomTaskCompletedEvent signals that a custom task has completed.
type CustomTaskCompletedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskCompletedEvent) workflowEvent()      {}
func (*CustomTaskCompletedEvent) EventType() EventType { return EventCustomTaskCompleted }

// CustomTaskFailedEvent signals that a custom task has failed.
type CustomTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskFailedEvent) workflowEvent()      {}
func (*CustomTaskFailedEvent) EventType() EventType { return EventCustomTaskFailed }

// CustomTaskTimedOutEvent signals that a custom task timed out.
type CustomTaskTimedOutEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskTimedOutEvent) workflowEvent()      {}
func (*CustomTaskTimedOutEvent) EventType() EventType { return EventCustomTaskTimedOut }

// CustomTaskCanceledEvent signals that a custom task was canceled.
type CustomTaskCanceledEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*CustomTaskCanceledEvent) workflowEvent()      {}
func (*CustomTaskCanceledEvent) EventType() EventType { return EventCustomTaskCanceled }

// ActivityTaskStartedEvent signals that an activity task has started.
type ActivityTaskStartedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskStartedEvent) workflowEvent()      {}
func (*ActivityTaskStartedEvent) EventType() EventType { return EventActivityTaskStarted }

// ActivityTaskCompletedEvent signals that an activity task has completed.
type ActivityTaskCompletedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskCompletedEvent) workflowEvent()      {}
func (*ActivityTaskCompletedEvent) EventType() EventType { return EventActivityTaskCompleted }

// ActivityTaskRetryingEvent signals that an activity task is being retried.
type ActivityTaskRetryingEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskRetryingEvent) workflowEvent()      {}
func (*ActivityTaskRetryingEvent) EventType() EventType { return EventActivityTaskRetrying }

// ActivityTaskFailedEvent signals that an activity task has failed.
type ActivityTaskFailedEvent struct {
	eventBase
	Attributes json.RawMessage `json:"attributes"`
}

func (*ActivityTaskFailedEvent) workflowEvent()      {}
func (*ActivityTaskFailedEvent) EventType() EventType { return EventActivityTaskFailed }

// UnknownEvent holds an event with an unrecognized event_type.
// This prevents the SDK from breaking when new event types are added.
type UnknownEvent struct {
	eventBase
	RawType string
	Raw     json.RawMessage
}

func (*UnknownEvent) workflowEvent()       {}
func (e *UnknownEvent) EventType() EventType { return EventType(e.RawType) }

// UnmarshalEvent dispatches JSON to the concrete Event type
// based on the "event_type" discriminator field.
func UnmarshalEvent(data []byte) (Event, error) {
	var probe struct {
		Type string `json:"event_type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("mistral: unmarshal workflow event: %w", err)
	}
	switch probe.Type {
	case string(EventWorkflowStarted):
		var e WorkflowExecutionStartedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowCompleted):
		var e WorkflowExecutionCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowFailed):
		var e WorkflowExecutionFailedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowCanceled):
		var e WorkflowExecutionCanceledEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowContinuedAsNew):
		var e WorkflowExecutionContinuedAsNewEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowTaskTimedOut):
		var e WorkflowTaskTimedOutEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventWorkflowTaskFailed):
		var e WorkflowTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskStarted):
		var e CustomTaskStartedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskInProgress):
		var e CustomTaskInProgressEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskCompleted):
		var e CustomTaskCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskFailed):
		var e CustomTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskTimedOut):
		var e CustomTaskTimedOutEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventCustomTaskCanceled):
		var e CustomTaskCanceledEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventActivityTaskStarted):
		var e ActivityTaskStartedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventActivityTaskCompleted):
		var e ActivityTaskCompletedEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventActivityTaskRetrying):
		var e ActivityTaskRetryingEvent
		return &e, json.Unmarshal(data, &e)
	case string(EventActivityTaskFailed):
		var e ActivityTaskFailedEvent
		return &e, json.Unmarshal(data, &e)
	default:
		var base eventBase
		if err := json.Unmarshal(data, &base); err != nil {
			return nil, fmt.Errorf("mistral: unmarshal workflow event base: %w", err)
		}
		return &UnknownEvent{
			eventBase: base,
			RawType:   probe.Type,
			Raw:       json.RawMessage(data),
		}, nil
	}
}

// StreamPayload is a single SSE payload from the workflow event stream.
type StreamPayload struct {
	Data json.RawMessage `json:"data"`
}

// StreamWorkflowContext holds context for a workflow event stream.
type StreamWorkflowContext struct {
	WorkflowName string `json:"workflow_name"`
	ExecutionID  string `json:"execution_id"`
}

// EventStreamParams holds query parameters for streaming workflow events.
type EventStreamParams struct {
	Source      *EventSource
	LastEventID *string
	Scope       *Scope
}

// EventListParams holds query parameters for listing workflow events.
type EventListParams struct {
	Source *EventSource
	Scope  *Scope
	Cursor *string
	Limit  *int
}

// EventListResponse is the response from listing workflow events.
type EventListResponse struct {
	Events     []json.RawMessage `json:"events"`
	NextCursor *string           `json:"next_cursor,omitempty"`
}
