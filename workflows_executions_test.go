package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
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
	var lastPayload *workflow.StreamPayload
	for stream.Next() {
		events = append(events, stream.Current())
		lastPayload = stream.CurrentPayload()
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
	if lastPayload.WorkflowContext.WorkflowName != "my-flow" {
		t.Errorf("got workflow context name %q", lastPayload.WorkflowContext.WorkflowName)
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
