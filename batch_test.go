package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/batch"
)

func TestCreateBatchJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/batch/jobs" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "batch-123", "object": "batch",
			"input_files": []string{"file-1"}, "endpoint": "/v1/chat/completions",
			"model": "mistral-small-latest", "status": "QUEUED",
			"errors": []any{}, "created_at": 1700000000,
			"total_requests": 10, "completed_requests": 0,
			"succeeded_requests": 0, "failed_requests": 0,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	model := "mistral-small-latest"
	job, err := client.CreateBatchJob(context.Background(), &batch.JobIn{
		Endpoint:   "/v1/chat/completions",
		InputFiles: []string{"file-1"},
		Model:      &model,
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "batch-123" {
		t.Errorf("got id %q", job.ID)
	}
	if job.Status != "QUEUED" {
		t.Errorf("got status %q", job.Status)
	}
	if job.TotalRequests != 10 {
		t.Errorf("got total_requests %d", job.TotalRequests)
	}
}

func TestListBatchJobs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id": "batch-1", "input_files": []string{}, "endpoint": "/v1/embeddings",
				"status": "SUCCESS", "errors": []any{}, "created_at": 1700000000,
				"total_requests": 5, "completed_requests": 5,
				"succeeded_requests": 5, "failed_requests": 0,
			}},
			"object": "list", "total": 1,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	jobs, err := client.ListBatchJobs(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if jobs.Total != 1 {
		t.Errorf("got total %d", jobs.Total)
	}
}

func TestGetBatchJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/batch/jobs/batch-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "batch-123", "input_files": []string{}, "endpoint": "/v1/chat/completions",
			"status": "RUNNING", "errors": []any{}, "created_at": 1700000000,
			"total_requests": 10, "completed_requests": 3,
			"succeeded_requests": 3, "failed_requests": 0,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	job, err := client.GetBatchJob(context.Background(), "batch-123")
	if err != nil {
		t.Fatal(err)
	}
	if job.CompletedRequests != 3 {
		t.Errorf("got completed_requests %d", job.CompletedRequests)
	}
}

func TestCancelBatchJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "batch-123", "input_files": []string{}, "endpoint": "/v1/chat/completions",
			"status": "CANCELLATION_REQUESTED", "errors": []any{}, "created_at": 1700000000,
			"total_requests": 10, "completed_requests": 3,
			"succeeded_requests": 3, "failed_requests": 0,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	job, err := client.CancelBatchJob(context.Background(), "batch-123")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "CANCELLATION_REQUESTED" {
		t.Errorf("got status %q", job.Status)
	}
}

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
