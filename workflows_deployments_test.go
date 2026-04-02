package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWorkflowDeployments_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/deployments" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"deployments": []map[string]any{
				{"id": "dep-1", "name": "prod", "is_active": true, "created_at": "2026-01-01", "updated_at": "2026-01-01"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListWorkflowDeployments(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Deployments) != 1 {
		t.Fatalf("got %d deployments", len(resp.Deployments))
	}
	if resp.Deployments[0].Name != "prod" {
		t.Errorf("got name %q", resp.Deployments[0].Name)
	}
}

func TestGetWorkflowDeployment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/deployments/dep-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "dep-1", "name": "prod", "is_active": true,
			"created_at": "2026-01-01", "updated_at": "2026-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	dep, err := client.GetWorkflowDeployment(context.Background(), "dep-1")
	if err != nil {
		t.Fatal(err)
	}
	if dep.ID != "dep-1" {
		t.Errorf("got id %q", dep.ID)
	}
}
