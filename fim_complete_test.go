package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/chat"
	"github.com/VikingOwl91/mistral-go-sdk/fim"
)

func TestFIMComplete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/fim/completions" {
			t.Errorf("expected /v1/fim/completions, got %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["prompt"] != "def add(a, b):" {
			t.Errorf("expected prompt, got %v", body["prompt"])
		}
		if body["suffix"] != "return result" {
			t.Errorf("expected suffix, got %v", body["suffix"])
		}
		if body["model"] != "codestral-latest" {
			t.Errorf("expected model codestral-latest, got %v", body["model"])
		}
		if body["stream"] != false {
			t.Errorf("expected stream=false, got %v", body["stream"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "fim-1", "object": "chat.completion",
			"model": "codestral-latest", "created": 1234567890,
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]any{"role": "assistant", "content": "\n    result = a + b\n    "},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 15, "completion_tokens": 10, "total_tokens": 25},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	suffix := "return result"
	resp, err := client.FIMComplete(context.Background(), &fim.CompletionRequest{
		Model:  "codestral-latest",
		Prompt: "def add(a, b):",
		Suffix: &suffix,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "fim-1" {
		t.Errorf("got id %q", resp.ID)
	}
	if resp.Choices[0].Message.Content.String() != "\n    result = a + b\n    " {
		t.Errorf("got content %q", resp.Choices[0].Message.Content.String())
	}
}

func TestFIMComplete_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["temperature"] != 0.2 {
			t.Errorf("expected temperature=0.2, got %v", body["temperature"])
		}
		if body["max_tokens"] != float64(50) {
			t.Errorf("expected max_tokens=50, got %v", body["max_tokens"])
		}
		if body["min_tokens"] != float64(10) {
			t.Errorf("expected min_tokens=10, got %v", body["min_tokens"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "fim-2", "object": "chat.completion",
			"model": "codestral-latest", "created": 0,
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "code"},
				"finish_reason": "length",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	temp := 0.2
	maxTok := 50
	minTok := 10
	_, err := client.FIMComplete(context.Background(), &fim.CompletionRequest{
		Model:       "codestral-latest",
		Prompt:      "fn main() {",
		Temperature: &temp,
		MaxTokens:   &maxTok,
		MinTokens:   &minTok,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFIMCompleteStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		for _, content := range []string{"\n    ", "result = a + b", "\n    "} {
			chunk := chat.CompletionChunk{
				ID: "fc", Model: "codestral-latest",
				Choices: []chat.CompletionStreamChoice{{
					Index: 0,
					Delta: chat.DeltaMessage{Content: chat.TextContent(content)},
				}},
			}
			data, _ := json.Marshal(chunk)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.FIMCompleteStream(context.Background(), &fim.CompletionRequest{
		Model:  "codestral-latest",
		Prompt: "def add(a, b):",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var count int
	for stream.Next() {
		count++
	}
	if stream.Err() != nil {
		t.Fatal(stream.Err())
	}
	if count != 3 {
		t.Errorf("got %d chunks, want 3", count)
	}
}

func TestFIMComplete_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{"message": "model not found"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.FIMComplete(context.Background(), &fim.CompletionRequest{
		Model:  "bad-model",
		Prompt: "code",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found, got: %v", err)
	}
}
