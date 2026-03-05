package conversation

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEvent_ResponseStarted(t *testing.T) {
	data := []byte(`{"type":"conversation.response.started","created_at":"2024-01-01T00:00:00Z","conversation_id":"conv-123"}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*ResponseStartedEvent)
	if !ok {
		t.Fatalf("expected *ResponseStartedEvent, got %T", event)
	}
	if e.ConversationID != "conv-123" {
		t.Errorf("got %q", e.ConversationID)
	}
}

func TestUnmarshalEvent_ResponseDone(t *testing.T) {
	data := []byte(`{"type":"conversation.response.done","created_at":"2024-01-01T00:00:00Z","usage":{"prompt_tokens":10,"completion_tokens":5,"total_tokens":15}}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*ResponseDoneEvent)
	if !ok {
		t.Fatalf("expected *ResponseDoneEvent, got %T", event)
	}
	if e.Usage.TotalTokens != 15 {
		t.Errorf("got total_tokens %d", e.Usage.TotalTokens)
	}
}

func TestUnmarshalEvent_ResponseError(t *testing.T) {
	data := []byte(`{"type":"conversation.response.error","created_at":"2024-01-01T00:00:00Z","message":"error occurred","code":500}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*ResponseErrorEvent)
	if !ok {
		t.Fatalf("expected *ResponseErrorEvent, got %T", event)
	}
	if e.Message != "error occurred" {
		t.Errorf("got %q", e.Message)
	}
	if e.Code != 500 {
		t.Errorf("got code %d", e.Code)
	}
}

func TestUnmarshalEvent_MessageOutput(t *testing.T) {
	data := []byte(`{"type":"message.output.delta","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"m1","content_index":0,"content":"Hello","role":"assistant"}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*MessageOutputEvent)
	if !ok {
		t.Fatalf("expected *MessageOutputEvent, got %T", event)
	}
	if e.ID != "m1" {
		t.Errorf("got id %q", e.ID)
	}
	var content string
	json.Unmarshal(e.Content, &content)
	if content != "Hello" {
		t.Errorf("got content %q", content)
	}
}

func TestUnmarshalEvent_FunctionCall(t *testing.T) {
	data := []byte(`{"type":"function.call.delta","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"fc1","name":"search","tool_call_id":"tc1","arguments":"{\"q\":\"test\"}"}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*FunctionCallEvent)
	if !ok {
		t.Fatalf("expected *FunctionCallEvent, got %T", event)
	}
	if e.Name != "search" {
		t.Errorf("got name %q", e.Name)
	}
	if e.ToolCallID != "tc1" {
		t.Errorf("got tool_call_id %q", e.ToolCallID)
	}
}

func TestUnmarshalEvent_ToolExecutionStarted(t *testing.T) {
	data := []byte(`{"type":"tool.execution.started","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"te1","name":"web_search","arguments":"query"}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*ToolExecutionStartedEvent)
	if !ok {
		t.Fatalf("expected *ToolExecutionStartedEvent, got %T", event)
	}
	if e.Name != "web_search" {
		t.Errorf("got %q", e.Name)
	}
}

func TestUnmarshalEvent_ToolExecutionDone(t *testing.T) {
	data := []byte(`{"type":"tool.execution.done","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"te1","name":"web_search"}`)
	event, err := UnmarshalEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	_, ok := event.(*ToolExecutionDoneEvent)
	if !ok {
		t.Fatalf("expected *ToolExecutionDoneEvent, got %T", event)
	}
}

func TestUnmarshalEvent_AgentHandoff(t *testing.T) {
	started := []byte(`{"type":"agent.handoff.started","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"h1","previous_agent_id":"a1","previous_agent_name":"A"}`)
	event, err := UnmarshalEvent(started)
	if err != nil {
		t.Fatal(err)
	}
	hs, ok := event.(*AgentHandoffStartedEvent)
	if !ok {
		t.Fatalf("expected *AgentHandoffStartedEvent, got %T", event)
	}
	if hs.PreviousAgentID != "a1" {
		t.Errorf("got %q", hs.PreviousAgentID)
	}

	done := []byte(`{"type":"agent.handoff.done","created_at":"2024-01-01T00:00:00Z","output_index":0,"id":"h1","next_agent_id":"a2","next_agent_name":"B"}`)
	event, err = UnmarshalEvent(done)
	if err != nil {
		t.Fatal(err)
	}
	hd, ok := event.(*AgentHandoffDoneEvent)
	if !ok {
		t.Fatalf("expected *AgentHandoffDoneEvent, got %T", event)
	}
	if hd.NextAgentID != "a2" {
		t.Errorf("got %q", hd.NextAgentID)
	}
}

func TestUnmarshalEvent_Unknown(t *testing.T) {
	event, err := UnmarshalEvent([]byte(`{"type":"future.event","id":"x"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	u, ok := event.(*UnknownEvent)
	if !ok {
		t.Fatalf("expected *UnknownEvent, got %T", event)
	}
	if u.Type != "future.event" {
		t.Errorf("got type %q", u.Type)
	}
	if len(u.Raw) == 0 {
		t.Error("expected raw data")
	}
}

