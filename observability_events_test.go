package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/observability"
)

func TestSearchChatCompletionEvents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/observability/chat-completion-events/search" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"results": []map[string]any{
				{"event_id": "ev-1", "correlation_id": "c1", "created_at": "t", "nb_input_tokens": 10, "nb_output_tokens": 5},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.SearchChatCompletionEvents(context.Background(), &observability.SearchEventsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("got %d results", len(resp.Results))
	}
	if resp.Results[0].EventID != "ev-1" {
		t.Errorf("got event_id %q", resp.Results[0].EventID)
	}
}

func TestSearchChatCompletionEventIDs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/chat-completion-events/search-ids" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"completion_event_ids": []string{"ev-1", "ev-2"},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.SearchChatCompletionEventIDs(context.Background(), &observability.SearchEventIDsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.CompletionEventIDs) != 2 {
		t.Errorf("got %d ids", len(resp.CompletionEventIDs))
	}
}

func TestGetChatCompletionEvent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/chat-completion-events/ev-1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"event_id": "ev-1", "correlation_id": "c1", "created_at": "t",
			"nb_input_tokens": 10, "nb_output_tokens": 5, "nb_messages": 2,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetChatCompletionEvent(context.Background(), "ev-1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.EventID != "ev-1" {
		t.Errorf("got event_id %q", resp.EventID)
	}
}

func TestGetSimilarChatCompletionEvents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/observability/chat-completion-events/ev-1/similar-events" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"count": 0, "results": []any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetSimilarChatCompletionEvents(context.Background(), "ev-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Count != 0 {
		t.Errorf("got count %d", resp.Count)
	}
}
