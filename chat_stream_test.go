package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

func TestChatCompleteStream_Success(t *testing.T) {
	chunks := []chat.CompletionChunk{
		{
			ID:    "chunk-1",
			Model: "mistral-small-latest",
			Choices: []chat.CompletionStreamChoice{{
				Index: 0,
				Delta: chat.DeltaMessage{Role: "assistant"},
			}},
		},
		{
			ID:    "chunk-2",
			Model: "mistral-small-latest",
			Choices: []chat.CompletionStreamChoice{{
				Index: 0,
				Delta: chat.DeltaMessage{Content: chat.TextContent("Hello")},
			}},
		},
		{
			ID:    "chunk-3",
			Model: "mistral-small-latest",
			Choices: []chat.CompletionStreamChoice{{
				Index: 0,
				Delta: chat.DeltaMessage{Content: chat.TextContent(" world!")},
			}},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true, got %v", body["stream"])
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		for _, chunk := range chunks {
			data, _ := json.Marshal(chunk)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Hi")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var received []chat.CompletionChunk
	for stream.Next() {
		received = append(received, stream.Current())
	}
	if stream.Err() != nil {
		t.Fatal(stream.Err())
	}
	if len(received) != 3 {
		t.Fatalf("got %d chunks, want 3", len(received))
	}
	if received[0].Choices[0].Delta.Role != "assistant" {
		t.Errorf("expected first chunk role=assistant")
	}
	if received[1].Choices[0].Delta.Content.String() != "Hello" {
		t.Errorf("got %q", received[1].Choices[0].Delta.Content.String())
	}
	if received[2].Choices[0].Delta.Content.String() != " world!" {
		t.Errorf("got %q", received[2].Choices[0].Delta.Content.String())
	}
}

func TestChatCompleteStream_CollectContent(t *testing.T) {
	words := []string{"The", " quick", " brown", " fox"}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)
		for _, word := range words {
			chunk := chat.CompletionChunk{
				ID:    "c",
				Model: "m",
				Choices: []chat.CompletionStreamChoice{{
					Index: 0,
					Delta: chat.DeltaMessage{Content: chat.TextContent(word)},
				}},
			}
			data, _ := json.Marshal(chunk)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
		stop := chat.FinishReasonStop
		final := chat.CompletionChunk{
			ID:    "c",
			Model: "m",
			Choices: []chat.CompletionStreamChoice{{
				Index:        0,
				Delta:        chat.DeltaMessage{},
				FinishReason: &stop,
			}},
			Usage: &chat.UsageInfo{PromptTokens: 5, CompletionTokens: 4, TotalTokens: 9},
		}
		data, _ := json.Marshal(final)
		fmt.Fprintf(w, "data: %s\n\n", data)
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var sb strings.Builder
	var lastChunk chat.CompletionChunk
	for stream.Next() {
		lastChunk = stream.Current()
		if len(lastChunk.Choices) > 0 {
			sb.WriteString(lastChunk.Choices[0].Delta.Content.String())
		}
	}
	if stream.Err() != nil {
		t.Fatal(stream.Err())
	}
	if sb.String() != "The quick brown fox" {
		t.Errorf("got %q", sb.String())
	}
	if lastChunk.Usage == nil {
		t.Fatal("expected usage in final chunk")
	}
	if lastChunk.Usage.TotalTokens != 9 {
		t.Errorf("got total_tokens=%d", lastChunk.Usage.TotalTokens)
	}
}

func TestChatCompleteStream_WithToolCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		chunk := chat.CompletionChunk{
			ID:    "c",
			Model: "m",
			Choices: []chat.CompletionStreamChoice{{
				Index: 0,
				Delta: chat.DeltaMessage{
					ToolCalls: []chat.ToolCall{{
						ID:       "call_1",
						Type:     "function",
						Function: chat.FunctionCall{Name: "get_weather", Arguments: `{"city":"Paris"}`},
					}},
				},
			}},
		}
		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Weather?")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	if !stream.Next() {
		t.Fatalf("expected chunk, err: %v", stream.Err())
	}
	tc := stream.Current().Choices[0].Delta.ToolCalls
	if len(tc) != 1 || tc[0].Function.Name != "get_weather" {
		t.Errorf("got tool calls %+v", tc)
	}
}

func TestChatCompleteStream_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "invalid key",
			"type":    "auth_error",
		})
	}))
	defer server.Close()

	client := NewClient("bad", WithBaseURL(server.URL))
	_, err := client.ChatCompleteStream(context.Background(), &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsAuth(err) {
		t.Errorf("expected auth error, got: %v", err)
	}
}
