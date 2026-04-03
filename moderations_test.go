package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/classification"
	"github.com/VikingOwl91/mistral-go-sdk/moderation"
)

func TestModerate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/moderations" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["model"] != "mistral-moderation-latest" {
			t.Errorf("got model %v", body["model"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "mod-123",
			"model": "mistral-moderation-latest",
			"results": []map[string]any{{
				"categories":      map[string]any{"sexual": false, "violence_and_threats": false},
				"category_scores": map[string]any{"sexual": 0.001, "violence_and_threats": 0.0003},
			}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.Moderate(context.Background(), &moderation.Request{
		Model: "mistral-moderation-latest",
		Input: moderation.TextInput("Hello world"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "mod-123" {
		t.Errorf("got id %q", resp.ID)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("got %d results", len(resp.Results))
	}
	if resp.Results[0].Categories["sexual"] != false {
		t.Error("expected sexual=false")
	}
	if resp.Results[0].CategoryScores["sexual"] != 0.001 {
		t.Errorf("got score %f", resp.Results[0].CategoryScores["sexual"])
	}
}

func TestModerateChat_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/moderations" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "mod-456",
			"model": "mistral-moderation-latest",
			"results": []map[string]any{{
				"categories":      map[string]any{"pii": false},
				"category_scores": map[string]any{"pii": 0.002},
			}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ModerateChat(context.Background(), &moderation.ChatRequest{
		Model: "mistral-moderation-latest",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "mod-456" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestClassify_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/classifications" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":    "cls-123",
			"model": "mistral-moderation-latest",
			"results": []map[string]any{{
				"default": map[string]any{
					"scores": map[string]any{"safe": 0.99, "unsafe": 0.01},
				},
			}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.Classify(context.Background(), &classification.Request{
		Model: "mistral-moderation-latest",
		Input: classification.TextInput("Hello"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "cls-123" {
		t.Errorf("got id %q", resp.ID)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("got %d results", len(resp.Results))
	}
}

func TestClassifyChat_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/classifications" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id":      "cls-456",
			"model":   "mistral-moderation-latest",
			"results": []map[string]any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ClassifyChat(context.Background(), &classification.ChatRequest{
		Model: "mistral-moderation-latest",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "cls-456" {
		t.Errorf("got id %q", resp.ID)
	}
}
