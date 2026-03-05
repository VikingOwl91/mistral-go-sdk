package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/conversation"
)

func TestStartConversation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/conversations" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["inputs"] != "Hello" {
			t.Errorf("expected inputs=Hello, got %v", body["inputs"])
		}
		if body["model"] != "mistral-small-latest" {
			t.Errorf("expected model, got %v", body["model"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"object":          "conversation.response",
			"conversation_id": "conv-123",
			"outputs": []map[string]any{{
				"object": "entry", "id": "e1", "type": "message.output",
				"created_at": "2024-01-01T00:00:00Z",
				"role": "assistant", "content": "Hello! How can I help?",
			}},
			"usage": map[string]any{
				"prompt_tokens": 10, "completion_tokens": 8, "total_tokens": 18,
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.StartConversation(context.Background(), &conversation.StartRequest{
		Inputs: conversation.TextInputs("Hello"),
		Model:  "mistral-small-latest",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ConversationID != "conv-123" {
		t.Errorf("got conv id %q", resp.ConversationID)
	}
	if len(resp.Outputs) != 1 {
		t.Fatalf("got %d outputs", len(resp.Outputs))
	}
	out, ok := resp.Outputs[0].(*conversation.MessageOutputEntry)
	if !ok {
		t.Fatalf("expected *MessageOutputEntry, got %T", resp.Outputs[0])
	}
	if conversation.TextContent(out.Content) != "Hello! How can I help?" {
		t.Errorf("got %q", conversation.TextContent(out.Content))
	}
	if resp.Usage.TotalTokens != 18 {
		t.Errorf("got total_tokens %d", resp.Usage.TotalTokens)
	}
}

func TestAppendConversation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/conversations/conv-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "conversation.response", "conversation_id": "conv-123",
			"outputs": []map[string]any{{
				"object": "entry", "id": "e2", "type": "message.output",
				"created_at": "2024-01-01T00:00:00Z",
				"role": "assistant", "content": "The weather is sunny.",
			}},
			"usage": map[string]any{"prompt_tokens": 20, "completion_tokens": 5, "total_tokens": 25},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.AppendConversation(context.Background(), "conv-123", &conversation.AppendRequest{
		Inputs: conversation.TextInputs("What's the weather?"),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ConversationID != "conv-123" {
		t.Errorf("got %q", resp.ConversationID)
	}
}

func TestGetConversation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/conversations/conv-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "conversation", "id": "conv-123",
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:01:00Z",
			"model": "mistral-small-latest",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	conv, err := client.GetConversation(context.Background(), "conv-123")
	if err != nil {
		t.Fatal(err)
	}
	if conv.ID != "conv-123" {
		t.Errorf("got id %q", conv.ID)
	}
	if conv.Model != "mistral-small-latest" {
		t.Errorf("got model %q", conv.Model)
	}
}

func TestDeleteConversation_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteConversation(context.Background(), "conv-123")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetConversationHistory_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/conversations/conv-123/history" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"object": "conversation.history", "conversation_id": "conv-123",
			"entries": []map[string]any{
				{"object": "entry", "id": "e1", "type": "message.input", "created_at": "2024-01-01T00:00:00Z", "role": "user", "content": "Hi"},
				{"object": "entry", "id": "e2", "type": "message.output", "created_at": "2024-01-01T00:00:01Z", "role": "assistant", "content": "Hello!"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	history, err := client.GetConversationHistory(context.Background(), "conv-123")
	if err != nil {
		t.Fatal(err)
	}
	if len(history.Entries) != 2 {
		t.Fatalf("got %d entries", len(history.Entries))
	}
	if _, ok := history.Entries[0].(*conversation.MessageInputEntry); !ok {
		t.Errorf("expected *MessageInputEntry, got %T", history.Entries[0])
	}
	if _, ok := history.Entries[1].(*conversation.MessageOutputEntry); !ok {
		t.Errorf("expected *MessageOutputEntry, got %T", history.Entries[1])
	}
}

func TestStartConversationStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		events := []map[string]any{
			{"type": "conversation.response.started", "created_at": "2024-01-01T00:00:00Z", "conversation_id": "conv-456"},
			{"type": "message.output.delta", "created_at": "2024-01-01T00:00:00Z", "output_index": 0, "id": "m1", "content_index": 0, "content": "Hello", "role": "assistant"},
			{"type": "message.output.delta", "created_at": "2024-01-01T00:00:00Z", "output_index": 0, "id": "m1", "content_index": 0, "content": " world!", "role": "assistant"},
			{"type": "conversation.response.done", "created_at": "2024-01-01T00:00:01Z", "usage": map[string]any{"prompt_tokens": 5, "completion_tokens": 2, "total_tokens": 7}},
		}
		for _, ev := range events {
			data, _ := json.Marshal(ev)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.StartConversationStream(context.Background(), &conversation.StartRequest{
		Inputs: conversation.TextInputs("Hi"),
		Model:  "mistral-small-latest",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var events []conversation.Event
	for stream.Next() {
		events = append(events, stream.Current())
	}
	if stream.Err() != nil {
		t.Fatal(stream.Err())
	}
	if len(events) != 4 {
		t.Fatalf("got %d events, want 4", len(events))
	}

	started, ok := events[0].(*conversation.ResponseStartedEvent)
	if !ok {
		t.Fatalf("expected *ResponseStartedEvent, got %T", events[0])
	}
	if started.ConversationID != "conv-456" {
		t.Errorf("got conv id %q", started.ConversationID)
	}

	msg, ok := events[1].(*conversation.MessageOutputEvent)
	if !ok {
		t.Fatalf("expected *MessageOutputEvent, got %T", events[1])
	}
	if conversation.TextContent(msg.Content) != "Hello" {
		t.Errorf("got %q", conversation.TextContent(msg.Content))
	}

	done, ok := events[3].(*conversation.ResponseDoneEvent)
	if !ok {
		t.Fatalf("expected *ResponseDoneEvent, got %T", events[3])
	}
	if done.Usage.TotalTokens != 7 {
		t.Errorf("got total_tokens %d", done.Usage.TotalTokens)
	}
}
