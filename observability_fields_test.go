package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

func TestGetChatCompletionFields_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/observability/chat-completion-fields" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"field_definitions": []map[string]any{
				{
					"name":                "model",
					"label":               "Model",
					"type":                "ENUM",
					"supported_operators": []string{"eq", "neq", "includes"},
				},
			},
			"field_groups": []map[string]any{
				{"name": "request", "label": "Request"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetChatCompletionFields(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.FieldDefinitions) != 1 {
		t.Fatalf("got %d field definitions", len(resp.FieldDefinitions))
	}
	def := resp.FieldDefinitions[0]
	if def.Name != "model" || def.Type != observability.FieldTypeEnum {
		t.Errorf("unexpected field def: %+v", def)
	}
	if len(def.SupportedOperators) != 3 {
		t.Errorf("got %d operators", len(def.SupportedOperators))
	}
	if len(resp.FieldGroups) != 1 || resp.FieldGroups[0].Name != "request" {
		t.Errorf("unexpected groups: %+v", resp.FieldGroups)
	}
}

func TestGetChatCompletionFieldOptions_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/chat-completion-fields/model/options" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("operator"); got != "eq" {
			t.Errorf("got operator=%q want eq", got)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"options": []any{"mistral-small-latest", "mistral-large-latest", nil, true},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetChatCompletionFieldOptions(context.Background(), "model", observability.FieldOperatorEq)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Options) != 4 {
		t.Fatalf("got %d options", len(resp.Options))
	}
}

func TestGetChatCompletionFieldOptionsCounts_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/chat-completion-fields/model/options-counts" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"counts": []map[string]any{
				{"value": "mistral-small-latest", "count": 42},
				{"value": "mistral-large-latest", "count": 17},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetChatCompletionFieldOptionsCounts(context.Background(), "model", &observability.FieldOptionCountsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Counts) != 2 || resp.Counts[0].Count != 42 {
		t.Errorf("unexpected counts: %+v", resp.Counts)
	}
}
