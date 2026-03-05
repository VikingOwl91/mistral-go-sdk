package mistral

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/chat"
)

func TestChatComplete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("expected /v1/chat/completions, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["model"] != "mistral-small-latest" {
			t.Errorf("expected model mistral-small-latest, got %v", body["model"])
		}
		if body["stream"] != false {
			t.Errorf("expected stream=false, got %v", body["stream"])
		}
		msgs, ok := body["messages"].([]any)
		if !ok || len(msgs) != 1 {
			t.Errorf("expected 1 message, got %v", body["messages"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":      "chat-123",
			"object":  "chat.completion",
			"model":   "mistral-small-latest",
			"created": 1234567890,
			"choices": []map[string]any{
				{
					"index": 0,
					"message": map[string]any{
						"role":    "assistant",
						"content": "Hello!",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]any{
				"prompt_tokens":     10,
				"completion_tokens": 5,
				"total_tokens":      15,
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("Hi")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "chat-123" {
		t.Errorf("expected id chat-123, got %s", resp.ID)
	}
	if resp.Model != "mistral-small-latest" {
		t.Errorf("expected model mistral-small-latest, got %s", resp.Model)
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content.String() != "Hello!" {
		t.Errorf("expected content 'Hello!', got %q", resp.Choices[0].Message.Content.String())
	}
	if resp.Choices[0].FinishReason != chat.FinishReasonStop {
		t.Errorf("expected finish reason stop, got %s", resp.Choices[0].FinishReason)
	}
	if resp.Usage.PromptTokens != 10 {
		t.Errorf("expected prompt_tokens=10, got %d", resp.Usage.PromptTokens)
	}
	if resp.Usage.CompletionTokens != 5 {
		t.Errorf("expected completion_tokens=5, got %d", resp.Usage.CompletionTokens)
	}
	if resp.Usage.TotalTokens != 15 {
		t.Errorf("expected total_tokens=15, got %d", resp.Usage.TotalTokens)
	}
}

func TestChatComplete_WithSystemMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		msgs := body["messages"].([]any)
		if len(msgs) != 2 {
			t.Errorf("expected 2 messages, got %d", len(msgs))
		}
		first := msgs[0].(map[string]any)
		if first["role"] != "system" {
			t.Errorf("expected first message role=system, got %v", first["role"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "chat-456", "object": "chat.completion",
			"model": "mistral-small-latest", "created": 1234567890,
			"choices": []map[string]any{{
				"index":         0,
				"message":       map[string]any{"role": "assistant", "content": "I'm helpful!"},
				"finish_reason": "stop",
			}},
			"usage": map[string]any{"prompt_tokens": 20, "completion_tokens": 10, "total_tokens": 30},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.SystemMessage{Content: chat.TextContent("You are helpful.")},
			&chat.UserMessage{Content: chat.TextContent("Hi")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content.String() != "I'm helpful!" {
		t.Errorf("got %q", resp.Choices[0].Message.Content.String())
	}
}

func TestChatComplete_WithTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		tools, ok := body["tools"].([]any)
		if !ok || len(tools) != 1 {
			t.Errorf("expected 1 tool, got %v", body["tools"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "chat-789", "object": "chat.completion",
			"model": "mistral-small-latest", "created": 1234567890,
			"choices": []map[string]any{{
				"index": 0,
				"message": map[string]any{
					"role":    "assistant",
					"content": nil,
					"tool_calls": []map[string]any{{
						"id":   "call_1",
						"type": "function",
						"function": map[string]any{
							"name":      "get_weather",
							"arguments": `{"city":"Paris"}`,
						},
					}},
				},
				"finish_reason": "tool_calls",
			}},
			"usage": map[string]any{"prompt_tokens": 30, "completion_tokens": 15, "total_tokens": 45},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model: "mistral-small-latest",
		Messages: []chat.Message{
			&chat.UserMessage{Content: chat.TextContent("What's the weather in Paris?")},
		},
		Tools: []chat.Tool{{
			Type: "function",
			Function: chat.Function{
				Name:        "get_weather",
				Description: "Get weather",
				Parameters:  map[string]any{"type": "object", "properties": map[string]any{"city": map[string]any{"type": "string"}}},
			},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].FinishReason != chat.FinishReasonToolCalls {
		t.Errorf("expected finish reason tool_calls, got %s", resp.Choices[0].FinishReason)
	}
	if len(resp.Choices[0].Message.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(resp.Choices[0].Message.ToolCalls))
	}
	tc := resp.Choices[0].Message.ToolCalls[0]
	if tc.Function.Name != "get_weather" {
		t.Errorf("got function name %q", tc.Function.Name)
	}
	if tc.Function.Arguments != `{"city":"Paris"}` {
		t.Errorf("got arguments %q", tc.Function.Arguments)
	}
}

func TestChatComplete_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "invalid api key",
			"type":    "authentication_error",
		})
	}))
	defer server.Close()

	client := NewClient("bad-key", WithBaseURL(server.URL))
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "mistral-small-latest",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsAuth(err) {
		t.Errorf("expected auth error, got: %v", err)
	}
}

func TestChatComplete_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "model not found",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "nonexistent-model",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestChatComplete_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "rate limit exceeded",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:    "mistral-small-latest",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsRateLimit(err) {
		t.Errorf("expected rate limit error, got: %v", err)
	}
}

func TestChatComplete_RequestBody(t *testing.T) {
	temp := 0.7
	topP := 0.9
	maxTokens := 100
	seed := 42
	freqPenalty := 0.5
	presPenalty := 0.3
	n := 2
	parallel := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		var body map[string]any
		json.Unmarshal(bodyBytes, &body)

		if body["temperature"] != 0.7 {
			t.Errorf("expected temperature=0.7, got %v", body["temperature"])
		}
		if body["top_p"] != 0.9 {
			t.Errorf("expected top_p=0.9, got %v", body["top_p"])
		}
		if body["max_tokens"] != float64(100) {
			t.Errorf("expected max_tokens=100, got %v", body["max_tokens"])
		}
		if body["random_seed"] != float64(42) {
			t.Errorf("expected random_seed=42, got %v", body["random_seed"])
		}
		if body["frequency_penalty"] != 0.5 {
			t.Errorf("expected frequency_penalty=0.5, got %v", body["frequency_penalty"])
		}
		if body["presence_penalty"] != 0.3 {
			t.Errorf("expected presence_penalty=0.3, got %v", body["presence_penalty"])
		}
		if body["n"] != float64(2) {
			t.Errorf("expected n=2, got %v", body["n"])
		}
		if body["parallel_tool_calls"] != false {
			t.Errorf("expected parallel_tool_calls=false, got %v", body["parallel_tool_calls"])
		}
		stops, ok := body["stop"].([]any)
		if !ok || len(stops) != 2 {
			t.Errorf("expected 2 stop sequences, got %v", body["stop"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "chat-params", "object": "chat.completion",
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
	_, err := client.ChatComplete(context.Background(), &chat.CompletionRequest{
		Model:             "m",
		Messages:          []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
		Temperature:       &temp,
		TopP:              &topP,
		MaxTokens:         &maxTokens,
		RandomSeed:        &seed,
		FrequencyPenalty:  &freqPenalty,
		PresencePenalty:   &presPenalty,
		N:                 &n,
		ParallelToolCalls: &parallel,
		Stop:              []string{"<stop>", "</end>"},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestChatComplete_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never responds — context should cancel first
		select {}
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately
	_, err := client.ChatComplete(ctx, &chat.CompletionRequest{
		Model:    "m",
		Messages: []chat.Message{&chat.UserMessage{Content: chat.TextContent("Hi")}},
	})
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
}
