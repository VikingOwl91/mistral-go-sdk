package workflow

import (
	"encoding/json"
	"fmt"
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
				"event_id":                "evt",
				"event_timestamp":         1711929600000000000,
				"root_workflow_exec_id":   "root",
				"parent_workflow_exec_id": nil,
				"workflow_exec_id":        "exec",
				"workflow_run_id":         "run",
				"workflow_name":           "wf",
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
