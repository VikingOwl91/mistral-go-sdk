package chat

import (
	"encoding/json"
	"testing"
)

func TestSystemMessage_MarshalJSON(t *testing.T) {
	msg := &SystemMessage{Content: TextContent("You are helpful.")}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["role"] != "system" {
		t.Errorf("expected role=system, got %v", m["role"])
	}
	if m["content"] != "You are helpful." {
		t.Errorf("expected content='You are helpful.', got %v", m["content"])
	}
}

func TestUserMessage_MarshalJSON(t *testing.T) {
	msg := &UserMessage{Content: TextContent("Hello")}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["role"] != "user" {
		t.Errorf("expected role=user, got %v", m["role"])
	}
	if m["content"] != "Hello" {
		t.Errorf("expected content='Hello', got %v", m["content"])
	}
}

func TestUserMessage_WithChunks(t *testing.T) {
	msg := &UserMessage{
		Content: ChunksContent(
			&TextChunk{Text: "Look at this:"},
			&ImageURLChunk{ImageURL: ImageURL{URL: "https://example.com/img.png"}},
		),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	parts, ok := m["content"].([]any)
	if !ok {
		t.Fatalf("expected content to be array, got %T", m["content"])
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
}

func TestAssistantMessage_MarshalJSON(t *testing.T) {
	msg := &AssistantMessage{Content: TextContent("Hi there!")}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["role"] != "assistant" {
		t.Errorf("expected role=assistant, got %v", m["role"])
	}
	if m["content"] != "Hi there!" {
		t.Errorf("expected content='Hi there!', got %v", m["content"])
	}
	if _, exists := m["tool_calls"]; exists {
		t.Error("expected tool_calls to be omitted")
	}
}

func TestAssistantMessage_WithToolCalls(t *testing.T) {
	msg := &AssistantMessage{
		ToolCalls: []ToolCall{
			{
				ID:   "call_1",
				Type: "function",
				Function: FunctionCall{
					Name:      "get_weather",
					Arguments: `{"location":"Paris"}`,
				},
			},
		},
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	am, ok := parsed.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", parsed)
	}
	if len(am.ToolCalls) != 1 {
		t.Fatalf("got %d tool calls, want 1", len(am.ToolCalls))
	}
	if am.ToolCalls[0].Function.Name != "get_weather" {
		t.Errorf("got function %q, want %q", am.ToolCalls[0].Function.Name, "get_weather")
	}
	if am.ToolCalls[0].Function.Arguments != `{"location":"Paris"}` {
		t.Errorf("got args %q", am.ToolCalls[0].Function.Arguments)
	}
}

func TestToolMessage_MarshalJSON(t *testing.T) {
	msg := &ToolMessage{
		Content:    TextContent(`{"temp": 22}`),
		ToolCallID: "call_1",
		Name:       "get_weather",
	}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["role"] != "tool" {
		t.Errorf("expected role=tool, got %v", m["role"])
	}
	if m["tool_call_id"] != "call_1" {
		t.Errorf("expected tool_call_id=call_1, got %v", m["tool_call_id"])
	}
	if m["name"] != "get_weather" {
		t.Errorf("expected name=get_weather, got %v", m["name"])
	}
}

func TestUnmarshalMessage_System(t *testing.T) {
	data := []byte(`{"role":"system","content":"Hello"}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	sm, ok := msg.(*SystemMessage)
	if !ok {
		t.Fatalf("expected *SystemMessage, got %T", msg)
	}
	if sm.Content.String() != "Hello" {
		t.Errorf("got %q, want %q", sm.Content.String(), "Hello")
	}
	if sm.MessageRole() != "system" {
		t.Errorf("expected role system, got %s", sm.MessageRole())
	}
}

func TestUnmarshalMessage_User(t *testing.T) {
	data := []byte(`{"role":"user","content":"Hi"}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	um, ok := msg.(*UserMessage)
	if !ok {
		t.Fatalf("expected *UserMessage, got %T", msg)
	}
	if um.Content.String() != "Hi" {
		t.Errorf("got %q, want %q", um.Content.String(), "Hi")
	}
}

func TestUnmarshalMessage_Assistant(t *testing.T) {
	data := []byte(`{"role":"assistant","content":"Hello!","prefix":true}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	am, ok := msg.(*AssistantMessage)
	if !ok {
		t.Fatalf("expected *AssistantMessage, got %T", msg)
	}
	if am.Content.String() != "Hello!" {
		t.Errorf("got %q", am.Content.String())
	}
	if !am.Prefix {
		t.Error("expected prefix=true")
	}
}

func TestUnmarshalMessage_Tool(t *testing.T) {
	data := []byte(`{"role":"tool","content":"result","tool_call_id":"c1","name":"fn"}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	tm, ok := msg.(*ToolMessage)
	if !ok {
		t.Fatalf("expected *ToolMessage, got %T", msg)
	}
	if tm.Content.String() != "result" {
		t.Errorf("got %q", tm.Content.String())
	}
	if tm.ToolCallID != "c1" {
		t.Errorf("got tool_call_id %q", tm.ToolCallID)
	}
	if tm.Name != "fn" {
		t.Errorf("got name %q", tm.Name)
	}
}

func TestUnmarshalMessage_UnknownRole(t *testing.T) {
	data := []byte(`{"role":"developer","content":"test"}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	u, ok := msg.(*UnknownMessage)
	if !ok {
		t.Fatalf("expected *UnknownMessage, got %T", msg)
	}
	if u.MessageRole() != "developer" {
		t.Errorf("got role %q", u.MessageRole())
	}
	marshaled, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(marshaled) != string(data) {
		t.Errorf("round-trip failed: got %s", marshaled)
	}
}

func TestMessage_RoundTrip_AllTypes(t *testing.T) {
	messages := []Message{
		&SystemMessage{Content: TextContent("system prompt")},
		&UserMessage{Content: TextContent("user input")},
		&AssistantMessage{Content: TextContent("assistant reply")},
		&ToolMessage{Content: TextContent("tool result"), ToolCallID: "c1"},
	}
	for _, msg := range messages {
		data, err := json.Marshal(msg)
		if err != nil {
			t.Fatalf("marshal %T: %v", msg, err)
		}
		parsed, err := UnmarshalMessage(data)
		if err != nil {
			t.Fatalf("unmarshal %T: %v", msg, err)
		}
		if parsed.MessageRole() != msg.MessageRole() {
			t.Errorf("role mismatch: got %s, want %s", parsed.MessageRole(), msg.MessageRole())
		}
	}
}

func TestUserMessage_NullContent(t *testing.T) {
	data := []byte(`{"role":"user","content":null}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	um := msg.(*UserMessage)
	if !um.Content.IsNull() {
		t.Error("expected null content")
	}
}

func TestUserMessage_ArrayContent(t *testing.T) {
	data := []byte(`{"role":"user","content":[{"type":"text","text":"hello"},{"type":"image_url","image_url":{"url":"https://example.com/img.png"}}]}`)
	msg, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatal(err)
	}
	um := msg.(*UserMessage)
	if len(um.Content.Parts) != 2 {
		t.Fatalf("got %d parts, want 2", len(um.Content.Parts))
	}
	if _, ok := um.Content.Parts[0].(*TextChunk); !ok {
		t.Errorf("expected parts[0] to be *TextChunk, got %T", um.Content.Parts[0])
	}
	if _, ok := um.Content.Parts[1].(*ImageURLChunk); !ok {
		t.Errorf("expected parts[1] to be *ImageURLChunk, got %T", um.Content.Parts[1])
	}
}
