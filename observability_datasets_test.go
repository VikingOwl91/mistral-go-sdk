package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/observability"
)

func datasetJSON() map[string]any {
	return map[string]any{
		"id": "ds-1", "created_at": "t", "updated_at": "t",
		"name": "test-ds", "description": "d",
		"owner_id": "o", "workspace_id": "w",
	}
}

func TestCreateDataset_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/datasets" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(datasetJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateDataset(context.Background(), &observability.CreateDatasetRequest{
		Name:        "test-ds",
		Description: "d",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "ds-1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestListDatasets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"count":   1,
			"results": []any{datasetJSON()},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListDatasets(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 1 {
		t.Errorf("got count %d", resp.Count)
	}
}

func TestGetDataset_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/datasets/ds-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(datasetJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetDataset(context.Background(), "ds-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "test-ds" {
		t.Errorf("got name %q", resp.Name)
	}
}

func TestDeleteDataset_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	if err := client.DeleteDataset(context.Background(), "ds-1"); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDatasetRecord_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/datasets/ds-1/records" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": "rec-1", "created_at": "t", "updated_at": "t",
			"dataset_id": "ds-1", "source": "DIRECT_INPUT",
			"payload":    map[string]any{"messages": []any{}},
			"properties": map[string]any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateDatasetRecord(context.Background(), "ds-1", &observability.CreateRecordRequest{
		Payload:    observability.ConversationPayload{Messages: []map[string]any{}},
		Properties: map[string]any{},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "rec-1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestListDatasetRecords_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/datasets/ds-1/records" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"count": 0, "results": []any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListDatasetRecords(context.Background(), "ds-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Errorf("got count %d", resp.Count)
	}
}

func TestImportDatasetFromCampaign_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/datasets/ds-1/imports/from-campaign" {
			t.Errorf("got path %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": "task-1", "created_at": "t", "updated_at": "t",
			"creator_id": "u", "dataset_id": "ds-1", "workspace_id": "w",
			"status": "RUNNING",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ImportDatasetFromCampaign(context.Background(), "ds-1", &observability.ImportFromCampaignRequest{
		CampaignID: "camp-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != observability.TaskStatusRunning {
		t.Errorf("got status %q", resp.Status)
	}
}

func TestExportDatasetToJSONL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/datasets/ds-1/exports/to-jsonl" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"file_url": "https://storage.example.com/export.jsonl",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ExportDatasetToJSONL(context.Background(), "ds-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.FileURL != "https://storage.example.com/export.jsonl" {
		t.Errorf("got file_url %q", resp.FileURL)
	}
}

func TestGetDatasetTask_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/datasets/ds-1/tasks/task-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "task-1", "created_at": "t", "updated_at": "t",
			"creator_id": "u", "dataset_id": "ds-1", "workspace_id": "w",
			"status": "COMPLETED",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetDatasetTask(context.Background(), "ds-1", "task-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != observability.TaskStatusCompleted {
		t.Errorf("got status %q", resp.Status)
	}
}
