package agents

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalAgentTool_Function(t *testing.T) {
	data := []byte(`{"type":"function","function":{"name":"test"}}`)
	tool, err := UnmarshalAgentTool(data)
	if err != nil {
		t.Fatal(err)
	}
	ft, ok := tool.(*FunctionTool)
	if !ok {
		t.Fatalf("expected *FunctionTool, got %T", tool)
	}
	if ft.Type != "function" {
		t.Errorf("got type %q", ft.Type)
	}
}

func TestUnmarshalAgentTool_WebSearch(t *testing.T) {
	data := []byte(`{"type":"web_search"}`)
	tool, err := UnmarshalAgentTool(data)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := tool.(*WebSearchTool); !ok {
		t.Fatalf("expected *WebSearchTool, got %T", tool)
	}
}

func TestUnmarshalAgentTool_Unknown(t *testing.T) {
	data := []byte(`{"type":"future_tool","config":{}}`)
	tool, err := UnmarshalAgentTool(data)
	if err != nil {
		t.Fatal(err)
	}
	u, ok := tool.(*UnknownAgentTool)
	if !ok {
		t.Fatalf("expected *UnknownAgentTool, got %T", tool)
	}
	if u.Type != "future_tool" {
		t.Errorf("got type %q", u.Type)
	}
}

func TestAgentTools_RoundTrip(t *testing.T) {
	input := `[{"type":"web_search"},{"type":"code_interpreter"},{"type":"future_tool","data":"x"}]`
	var tools AgentTools
	if err := json.Unmarshal([]byte(input), &tools); err != nil {
		t.Fatal(err)
	}
	if len(tools) != 3 {
		t.Fatalf("got %d tools, want 3", len(tools))
	}
	if _, ok := tools[0].(*WebSearchTool); !ok {
		t.Errorf("tools[0]: expected *WebSearchTool, got %T", tools[0])
	}
	if _, ok := tools[1].(*CodeInterpreterTool); !ok {
		t.Errorf("tools[1]: expected *CodeInterpreterTool, got %T", tools[1])
	}
	if _, ok := tools[2].(*UnknownAgentTool); !ok {
		t.Errorf("tools[2]: expected *UnknownAgentTool, got %T", tools[2])
	}
}

func TestUnmarshalAgentTool_Connector(t *testing.T) {
	data := []byte(`{"type":"connector","connector_id":"my-connector","authorization":{"type":"api-key","value":"sk-test"}}`)
	tool, err := UnmarshalAgentTool(data)
	if err != nil {
		t.Fatal(err)
	}
	ct, ok := tool.(*ConnectorTool)
	if !ok {
		t.Fatalf("expected *ConnectorTool, got %T", tool)
	}
	if ct.ConnectorID != "my-connector" {
		t.Errorf("got connector_id %q", ct.ConnectorID)
	}
	if ct.Authorization == nil {
		t.Fatal("expected authorization")
	}
	if ct.Authorization.Type != "api-key" {
		t.Errorf("got auth type %q", ct.Authorization.Type)
	}
	if ct.Authorization.Value != "sk-test" {
		t.Errorf("got auth value %q", ct.Authorization.Value)
	}
}

func TestAgent_UnmarshalWithTools(t *testing.T) {
	data := []byte(`{
		"id":"ag-1","object":"agent","name":"A","model":"m",
		"version":1,"versions":[1],"created_at":"2024-01-01T00:00:00Z",
		"updated_at":"2024-01-01T00:00:00Z","deployment_chat":false,"source":"api",
		"tools":[{"type":"web_search"},{"type":"function","function":{"name":"test"}}]
	}`)
	var agent Agent
	if err := json.Unmarshal(data, &agent); err != nil {
		t.Fatal(err)
	}
	if len(agent.Tools) != 2 {
		t.Fatalf("got %d tools, want 2", len(agent.Tools))
	}
	if _, ok := agent.Tools[0].(*WebSearchTool); !ok {
		t.Errorf("tools[0]: expected *WebSearchTool, got %T", agent.Tools[0])
	}
	if _, ok := agent.Tools[1].(*FunctionTool); !ok {
		t.Errorf("tools[1]: expected *FunctionTool, got %T", agent.Tools[1])
	}
}
