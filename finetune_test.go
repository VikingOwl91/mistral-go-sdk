package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/finetune"
)

func TestCreateFineTuningJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		if r.URL.Path != "/v1/fine_tuning/jobs" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft-job-123", "object": "job", "model": "mistral-small-latest",
			"status": "QUEUED", "job_type": "completion",
			"created_at": 1700000000, "modified_at": 1700000000,
			"training_files": []string{"file-1"}, "auto_start": true,
			"hyperparameters": map[string]any{"training_steps": 100},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	job, err := client.CreateFineTuningJob(context.Background(), &finetune.JobIn{
		Model:           "mistral-small-latest",
		AutoStart:       true,
		Hyperparameters: json.RawMessage(`{"training_steps":100}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if job.ID != "ft-job-123" {
		t.Errorf("got id %q", job.ID)
	}
	if job.Status != "QUEUED" {
		t.Errorf("got status %q", job.Status)
	}
}

func TestListFineTuningJobs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"id": "ft-job-1", "object": "job", "model": "m",
				"status": "SUCCESS", "job_type": "completion",
				"created_at": 1700000000, "modified_at": 1700000000,
				"training_files": []string{}, "auto_start": true,
				"hyperparameters": map[string]any{},
			}},
			"object": "list", "total": 1,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	jobs, err := client.ListFineTuningJobs(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if jobs.Total != 1 {
		t.Errorf("got total %d", jobs.Total)
	}
	if len(jobs.Data) != 1 {
		t.Fatalf("got %d jobs", len(jobs.Data))
	}
}

func TestGetFineTuningJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/fine_tuning/jobs/ft-job-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft-job-123", "model": "m", "status": "RUNNING",
			"job_type": "completion", "created_at": 1700000000, "modified_at": 1700000000,
			"training_files": []string{}, "auto_start": true,
			"hyperparameters": map[string]any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	job, err := client.GetFineTuningJob(context.Background(), "ft-job-123")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "RUNNING" {
		t.Errorf("got status %q", job.Status)
	}
}

func TestCancelFineTuningJob_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft-job-123", "model": "m", "status": "CANCELLATION_REQUESTED",
			"job_type": "completion", "created_at": 1700000000, "modified_at": 1700000000,
			"training_files": []string{}, "auto_start": true,
			"hyperparameters": map[string]any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	job, err := client.CancelFineTuningJob(context.Background(), "ft-job-123")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "CANCELLATION_REQUESTED" {
		t.Errorf("got status %q", job.Status)
	}
}

func TestArchiveFineTunedModel_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft:model:123", "object": "model", "archived": true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ArchiveFineTunedModel(context.Background(), "ft:model:123")
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Archived {
		t.Error("expected archived=true")
	}
}

func TestUnarchiveFineTunedModel_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft:model:123", "object": "model", "archived": false,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.UnarchiveFineTunedModel(context.Background(), "ft:model:123")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Archived {
		t.Error("expected archived=false")
	}
}
