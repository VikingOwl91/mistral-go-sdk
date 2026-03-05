package chat

import (
	"encoding/json"
	"testing"
)

func TestFunctionCall_MarshalJSON(t *testing.T) {
	fc := FunctionCall{Name: "get_weather", Arguments: `{"city":"Paris"}`}
	data, err := json.Marshal(fc)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["name"] != "get_weather" {
		t.Errorf("got name %v", m["name"])
	}
	if m["arguments"] != `{"city":"Paris"}` {
		t.Errorf("got arguments %v", m["arguments"])
	}
}

func TestFunctionCall_UnmarshalJSON_StringArgs(t *testing.T) {
	data := []byte(`{"name":"fn","arguments":"{\"key\":\"val\"}"}`)
	var fc FunctionCall
	if err := json.Unmarshal(data, &fc); err != nil {
		t.Fatal(err)
	}
	if fc.Name != "fn" {
		t.Errorf("got name %q", fc.Name)
	}
	if fc.Arguments != `{"key":"val"}` {
		t.Errorf("got arguments %q", fc.Arguments)
	}
}

func TestFunctionCall_UnmarshalJSON_ObjectArgs(t *testing.T) {
	data := []byte(`{"name":"fn","arguments":{"key":"val"}}`)
	var fc FunctionCall
	if err := json.Unmarshal(data, &fc); err != nil {
		t.Fatal(err)
	}
	if fc.Name != "fn" {
		t.Errorf("got name %q", fc.Name)
	}
	if fc.Arguments != `{"key":"val"}` {
		t.Errorf("got arguments %q", fc.Arguments)
	}
}

func TestToolChoice_MarshalJSON_Mode(t *testing.T) {
	tc := ToolChoice{Mode: ToolChoiceAuto}
	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `"auto"` {
		t.Errorf("got %s, want %q", data, "auto")
	}
}

func TestToolChoice_MarshalJSON_Function(t *testing.T) {
	tc := ToolChoice{Function: &FunctionName{Name: "get_weather"}}
	data, err := json.Marshal(tc)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	json.Unmarshal(data, &m)
	if m["type"] != "function" {
		t.Errorf("expected type=function, got %v", m["type"])
	}
	fn, ok := m["function"].(map[string]any)
	if !ok {
		t.Fatalf("expected function object, got %T", m["function"])
	}
	if fn["name"] != "get_weather" {
		t.Errorf("got function name %v", fn["name"])
	}
}

func TestToolChoice_UnmarshalJSON_Mode(t *testing.T) {
	var tc ToolChoice
	if err := json.Unmarshal([]byte(`"none"`), &tc); err != nil {
		t.Fatal(err)
	}
	if tc.Mode != ToolChoiceNone {
		t.Errorf("got mode %q, want %q", tc.Mode, ToolChoiceNone)
	}
	if tc.Function != nil {
		t.Error("expected nil function")
	}
}

func TestToolChoice_UnmarshalJSON_Function(t *testing.T) {
	var tc ToolChoice
	if err := json.Unmarshal([]byte(`{"type":"function","function":{"name":"fn"}}`), &tc); err != nil {
		t.Fatal(err)
	}
	if tc.Function == nil {
		t.Fatal("expected non-nil function")
	}
	if tc.Function.Name != "fn" {
		t.Errorf("got function name %q", tc.Function.Name)
	}
}

func TestToolChoice_RoundTrip(t *testing.T) {
	tests := []ToolChoice{
		{Mode: ToolChoiceAuto},
		{Mode: ToolChoiceNone},
		{Mode: ToolChoiceAny},
		{Mode: ToolChoiceRequired},
		{Function: &FunctionName{Name: "my_func"}},
	}
	for _, tc := range tests {
		data, err := json.Marshal(tc)
		if err != nil {
			t.Fatalf("marshal %+v: %v", tc, err)
		}
		var decoded ToolChoice
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s: %v", data, err)
		}
		if tc.Function != nil {
			if decoded.Function == nil || decoded.Function.Name != tc.Function.Name {
				t.Errorf("function round-trip failed: got %+v", decoded)
			}
		} else {
			if decoded.Mode != tc.Mode {
				t.Errorf("mode round-trip failed: got %q, want %q", decoded.Mode, tc.Mode)
			}
		}
	}
}

func TestTool_MarshalJSON(t *testing.T) {
	tool := Tool{
		Type: "function",
		Function: Function{
			Name:        "get_weather",
			Description: "Get weather for a city",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"city": map[string]any{"type": "string"},
				},
				"required": []any{"city"},
			},
		},
	}
	data, err := json.Marshal(tool)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Tool
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Function.Name != "get_weather" {
		t.Errorf("got name %q", decoded.Function.Name)
	}
	if decoded.Function.Description != "Get weather for a city" {
		t.Errorf("got desc %q", decoded.Function.Description)
	}
}

func TestResponseFormat_JSON(t *testing.T) {
	rf := ResponseFormat{
		Type: ResponseFormatJSONSchema,
		JsonSchema: &JsonSchema{
			Name: "output",
			Schema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"answer": map[string]any{"type": "string"},
				},
			},
			Strict: true,
		},
	}
	data, err := json.Marshal(rf)
	if err != nil {
		t.Fatal(err)
	}
	var decoded ResponseFormat
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Type != ResponseFormatJSONSchema {
		t.Errorf("got type %q", decoded.Type)
	}
	if decoded.JsonSchema == nil {
		t.Fatal("expected non-nil json_schema")
	}
	if decoded.JsonSchema.Name != "output" {
		t.Errorf("got schema name %q", decoded.JsonSchema.Name)
	}
	if !decoded.JsonSchema.Strict {
		t.Error("expected strict=true")
	}
}
