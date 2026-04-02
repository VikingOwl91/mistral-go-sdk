package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/workflow"
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
