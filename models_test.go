package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/model"
)

func TestListModels_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/models" {
			t.Errorf("expected /v1/models, got %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"data": []map[string]any{
				{
					"id": "mistral-small-latest", "object": "model",
					"created": 1700000000, "owned_by": "mistralai",
					"type": "base", "max_context_length": 32768,
					"aliases": []string{"mistral-small"},
					"capabilities": map[string]any{
						"completion_chat": true, "function_calling": true,
						"vision": false, "fine_tuning": false,
					},
				},
				{
					"id": "ft:mistral-small:my-org:custom", "object": "model",
					"created": 1700000001, "owned_by": "my-org",
					"type": "fine-tuned", "max_context_length": 32768,
					"job": "job-123", "root": "mistral-small-latest",
					"archived": false, "aliases": []string{},
					"capabilities": map[string]any{
						"completion_chat": true, "fine_tuning": true,
					},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	list, err := client.ListModels(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if list.Object != "list" {
		t.Errorf("got object %q", list.Object)
	}
	if len(list.Data) != 2 {
		t.Fatalf("got %d models, want 2", len(list.Data))
	}

	base := list.Data[0]
	if base.ID != "mistral-small-latest" {
		t.Errorf("got id %q", base.ID)
	}
	if base.Type != "base" {
		t.Errorf("got type %q", base.Type)
	}
	if !base.Capabilities.CompletionChat {
		t.Error("expected completion_chat=true")
	}
	if base.MaxContextLength != 32768 {
		t.Errorf("got max_context_length=%d", base.MaxContextLength)
	}

	ft := list.Data[1]
	if ft.Type != "fine-tuned" {
		t.Errorf("got type %q", ft.Type)
	}
	if ft.Job != "job-123" {
		t.Errorf("got job %q", ft.Job)
	}
	if ft.Root != "mistral-small-latest" {
		t.Errorf("got root %q", ft.Root)
	}
}

func TestListModels_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("provider") != "mistralai" {
			t.Errorf("expected provider=mistralai, got %q", r.URL.Query().Get("provider"))
		}
		if r.URL.Query().Get("model") != "mistral-small" {
			t.Errorf("expected model=mistral-small, got %q", r.URL.Query().Get("model"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"data":   []map[string]any{},
		})
	}))
	defer server.Close()

	provider := "mistralai"
	modelName := "mistral-small"
	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.ListModels(context.Background(), &model.ListParams{
		Provider: &provider,
		Model:    &modelName,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetModel_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models/mistral-small-latest" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "mistral-small-latest", "object": "model",
			"created": 1700000000, "owned_by": "mistralai",
			"type": "base", "max_context_length": 32768,
			"capabilities": map[string]any{"completion_chat": true},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	card, err := client.GetModel(context.Background(), "mistral-small-latest")
	if err != nil {
		t.Fatal(err)
	}
	if card.ID != "mistral-small-latest" {
		t.Errorf("got id %q", card.ID)
	}
}

func TestGetModel_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(map[string]any{"message": "not found"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.GetModel(context.Background(), "nope")
	if !IsNotFound(err) {
		t.Errorf("expected not found, got: %v", err)
	}
}

func TestDeleteModel_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/v1/models/ft:model:123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "ft:model:123", "object": "model", "deleted": true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.DeleteModel(context.Background(), "ft:model:123")
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Deleted {
		t.Error("expected deleted=true")
	}
	if resp.ID != "ft:model:123" {
		t.Errorf("got id %q", resp.ID)
	}
}
