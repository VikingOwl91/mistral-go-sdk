package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
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
