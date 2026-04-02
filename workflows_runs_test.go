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
