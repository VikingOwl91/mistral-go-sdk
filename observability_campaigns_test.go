package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/observability"
)

func TestCreateCampaign_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/campaigns" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "test-campaign" {
			t.Errorf("got name %v", body["name"])
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": "camp-1", "name": "test-campaign", "description": "d",
			"created_at": "t", "updated_at": "t", "owner_id": "o",
			"workspace_id": "w", "max_nb_events": 100,
			"search_params": map[string]any{}, "judge": map[string]any{"id": "j1"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateCampaign(context.Background(), &observability.CreateCampaignRequest{
		Name:        "test-campaign",
		Description: "d",
		JudgeID:     "j1",
		MaxNbEvents: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "camp-1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestListCampaigns_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/campaigns" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"count":   1,
			"results": []map[string]any{{"id": "c1", "name": "c", "description": "d", "created_at": "t", "updated_at": "t", "owner_id": "o", "workspace_id": "w", "max_nb_events": 10, "search_params": map[string]any{}, "judge": map[string]any{"id": "j"}}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListCampaigns(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 1 {
		t.Errorf("got count %d", resp.Count)
	}
}

func TestGetCampaign_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/campaigns/camp-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "camp-1", "name": "c", "description": "d",
			"created_at": "t", "updated_at": "t", "owner_id": "o",
			"workspace_id": "w", "max_nb_events": 10,
			"search_params": map[string]any{}, "judge": map[string]any{"id": "j"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetCampaign(context.Background(), "camp-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "camp-1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestDeleteCampaign_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	if err := client.DeleteCampaign(context.Background(), "camp-1"); err != nil {
		t.Fatal(err)
	}
}

func TestGetCampaignStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/campaigns/camp-1/status" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"status": "COMPLETED"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetCampaignStatus(context.Background(), "camp-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != observability.TaskStatusCompleted {
		t.Errorf("got status %q", resp.Status)
	}
}

func TestListCampaignEvents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/campaigns/camp-1/selected-events" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"count":   0,
			"results": []any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListCampaignEvents(context.Background(), "camp-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Errorf("got count %d", resp.Count)
	}
}
