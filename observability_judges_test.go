package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

func judgeJSON() map[string]any {
	return map[string]any{
		"id": "j1", "created_at": "t", "updated_at": "t",
		"owner_id": "o", "workspace_id": "w", "name": "quality",
		"description": "d", "model_name": "m", "instructions": "i",
		"output": map[string]any{"type": "CLASSIFICATION", "options": []any{}},
	}
}

func TestCreateJudge_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/judges" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(judgeJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateJudge(context.Background(), &observability.CreateJudgeRequest{
		Name:         "quality",
		Description:  "d",
		ModelName:    "m",
		Instructions: "i",
		Tools:        []string{},
		Output:       json.RawMessage(`{"type":"CLASSIFICATION","options":[]}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "j1" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestListJudges_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"count":   1,
			"results": []any{judgeJSON()},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListJudges(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 1 {
		t.Errorf("got count %d", resp.Count)
	}
}

func TestGetJudge_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/judges/j1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(judgeJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetJudge(context.Background(), "j1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "quality" {
		t.Errorf("got name %q", resp.Name)
	}
}

func TestUpdateJudge_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(judgeJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.UpdateJudge(context.Background(), "j1", &observability.UpdateJudgeRequest{
		Name:         "quality",
		Description:  "d",
		ModelName:    "m",
		Instructions: "i",
		Tools:        []string{},
		Output:       json.RawMessage(`{"type":"CLASSIFICATION","options":[]}`),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJudge_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	if err := client.DeleteJudge(context.Background(), "j1"); err != nil {
		t.Fatal(err)
	}
}
