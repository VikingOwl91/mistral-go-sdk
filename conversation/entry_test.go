package conversation

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEntry_MessageInput(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e1","type":"message.input","created_at":"2024-01-01T00:00:00Z","role":"user","content":"Hello"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*MessageInputEntry)
	if !ok {
		t.Fatalf("expected *MessageInputEntry, got %T", entry)
	}
	if e.ID != "e1" {
		t.Errorf("got id %q", e.ID)
	}
	if e.Role != "user" {
		t.Errorf("got role %q", e.Role)
	}
	if TextContent(e.Content) != "Hello" {
		t.Errorf("got content %q", TextContent(e.Content))
	}
}

func TestUnmarshalEntry_MessageOutput(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e2","type":"message.output","created_at":"2024-01-01T00:00:00Z","role":"assistant","content":"Hi there!"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*MessageOutputEntry)
	if !ok {
		t.Fatalf("expected *MessageOutputEntry, got %T", entry)
	}
	if TextContent(e.Content) != "Hi there!" {
		t.Errorf("got content %q", TextContent(e.Content))
	}
}

func TestUnmarshalEntry_FunctionCall(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e3","type":"function.call","created_at":"2024-01-01T00:00:00Z","tool_call_id":"tc1","name":"get_weather","arguments":"{\"city\":\"Paris\"}"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*FunctionCallEntry)
	if !ok {
		t.Fatalf("expected *FunctionCallEntry, got %T", entry)
	}
	if e.Name != "get_weather" {
		t.Errorf("got name %q", e.Name)
	}
	if e.ToolCallID != "tc1" {
		t.Errorf("got tool_call_id %q", e.ToolCallID)
	}
}

func TestUnmarshalEntry_FunctionResult(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e4","type":"function.result","created_at":"2024-01-01T00:00:00Z","tool_call_id":"tc1","result":"{\"temp\":22}"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*FunctionResultEntry)
	if !ok {
		t.Fatalf("expected *FunctionResultEntry, got %T", entry)
	}
	if e.Result != `{"temp":22}` {
		t.Errorf("got result %q", e.Result)
	}
}

func TestUnmarshalEntry_ToolExecution(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e5","type":"tool.execution","created_at":"2024-01-01T00:00:00Z","name":"web_search","arguments":"query"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*ToolExecutionEntry)
	if !ok {
		t.Fatalf("expected *ToolExecutionEntry, got %T", entry)
	}
	if e.Name != "web_search" {
		t.Errorf("got name %q", e.Name)
	}
}

func TestUnmarshalEntry_AgentHandoff(t *testing.T) {
	data := []byte(`{"object":"entry","id":"e6","type":"agent.handoff","created_at":"2024-01-01T00:00:00Z","previous_agent_id":"a1","previous_agent_name":"Agent A","next_agent_id":"a2","next_agent_name":"Agent B"}`)
	entry, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := entry.(*AgentHandoffEntry)
	if !ok {
		t.Fatalf("expected *AgentHandoffEntry, got %T", entry)
	}
	if e.PreviousAgentName != "Agent A" {
		t.Errorf("got prev %q", e.PreviousAgentName)
	}
	if e.NextAgentName != "Agent B" {
		t.Errorf("got next %q", e.NextAgentName)
	}
}

func TestUnmarshalEntry_Unknown(t *testing.T) {
	_, err := UnmarshalEntry([]byte(`{"type":"unknown.type"}`))
	if err == nil {
		t.Error("expected error for unknown type")
	}
}

func TestTextContent_String(t *testing.T) {
	raw := json.RawMessage(`"Hello world"`)
	if TextContent(raw) != "Hello world" {
		t.Errorf("got %q", TextContent(raw))
	}
}

func TestTextContent_ChunkArray(t *testing.T) {
	raw := json.RawMessage(`[{"type":"text","text":"Hello "},{"type":"text","text":"world"}]`)
	if TextContent(raw) != "Hello world" {
		t.Errorf("got %q", TextContent(raw))
	}
}

func TestTextContent_Empty(t *testing.T) {
	if TextContent(nil) != "" {
		t.Error("expected empty for nil")
	}
	if TextContent(json.RawMessage{}) != "" {
		t.Error("expected empty for empty")
	}
}

func TestTextContent_MixedChunks(t *testing.T) {
	raw := json.RawMessage(`[{"type":"text","text":"Hello"},{"type":"tool_reference","tool":"web_search","title":"Result"},{"type":"text","text":" world"}]`)
	if TextContent(raw) != "Hello world" {
		t.Errorf("got %q", TextContent(raw))
	}
}
