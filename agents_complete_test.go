package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/agents"
	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

func TestAgentsComplete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/agents/completions" {
			t.Errorf("expected /v1/agents/completions, got %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["agent_id"] != "agent-123" {
			t.Errorf("expected agent_id=agent-123, got %v", body["agent_id"])
		}
		msgs := body["messages"].([]any)
		if len(msgs) != 1 {
			t.Errorf("expected 1 message, got %d", len(msgs))
		}
		if body["stream"] != false {
			t.Errorf("expected stream=false")
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "agent-resp-1", "object": "chat.completion",
			"model": "mistral-large-latest", "created": 1234567890,
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]any{"role": "assistant", "content": "Agent response"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 20, "completion_tokens": 5, "total_tokens": 25},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.AgentsComplete(context.Background(), &agents.CompletionRequest{
		AgentID: "agent-123",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Hello agent")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "agent-resp-1" {
		t.Errorf("got id %q", resp.ID)
	}
	if resp.Choices[0].Message.Content.String() != "Agent response" {
		t.Errorf("got content %q", resp.Choices[0].Message.Content.String())
	}
}

func TestAgentsComplete_WithTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		tools := body["tools"].([]any)
		if len(tools) != 1 {
			t.Errorf("expected 1 tool, got %d", len(tools))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "a2", "object": "chat.completion",
			"model": "m", "created": 0,
			"choices": []map[string]any{{
				"index": 0,
				"message": map[string]any{
					"role": "assistant", "content": nil,
					"tool_calls": []map[string]any{{
						"id": "tc1", "type": "function",
						"function": map[string]any{"name": "search", "arguments": `{"q":"test"}`},
					}},
				},
				"finish_reason": "tool_calls",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.AgentsComplete(context.Background(), &agents.CompletionRequest{
		AgentID: "agent-456",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Search for test")},
		},
		Tools: []chat.Tool{{
			Type: "function",
			Function: chat.Function{
				Name:       "search",
				Parameters: map[string]any{"type": "object", "properties": map[string]any{"q": map[string]any{"type": "string"}}},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Choices[0].Message.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call")
	}
	if resp.Choices[0].Message.ToolCalls[0].Function.Name != "search" {
		t.Errorf("got function %q", resp.Choices[0].Message.ToolCalls[0].Function.Name)
	}
}

func TestAgentsComplete_ReasoningEffort(t *testing.T) {
	effort := chat.ReasoningEffortHigh

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["reasoning_effort"] != "high" {
			t.Errorf("expected reasoning_effort=high, got %v", body["reasoning_effort"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "a-re", "object": "chat.completion",
			"model": "m", "created": 0,
			"choices": []map[string]any{{
				"index": 0, "message": map[string]any{"role": "assistant", "content": "ok"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.AgentsComplete(context.Background(), &agents.CompletionRequest{
		AgentID:         "agent-1",
		Messages:        []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
		ReasoningEffort: &effort,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAgentsCompleteStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		for _, word := range []string{"Hello", " from", " agent"} {
			chunk := chat.CompletionChunk{
				ID: "ac", Model: "m",
				Choices: []chat.CompletionStreamChoice{{
					Index: 0,
					Delta: chat.DeltaMessage{Content: chat.TextContent(word)},
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
	stream, err := client.AgentsCompleteStream(context.Background(), &agents.CompletionRequest{
		AgentID:  "agent-789",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
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
