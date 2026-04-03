package workflow

import (
	"encoding/json"
	"testing"
)

func TestCodeDefinition_RoundTrip(t *testing.T) {
	raw := `{
		"input_schema": {"type": "object", "properties": {"prompt": {"type": "string"}}},
		"output_schema": {"type": "object", "properties": {"result": {"type": "string"}}},
		"signals": [
			{"name": "cancel", "input_schema": {"type": "object"}, "description": "Cancel the workflow"}
		],
		"queries": [
			{"name": "status", "input_schema": {"type": "object"}, "description": "Get status", "output_schema": {"type": "string"}}
		],
		"updates": [
			{"name": "set_priority", "input_schema": {"type": "object", "properties": {"level": {"type": "integer"}}}, "description": "Set priority", "output_schema": null}
		],
		"enforce_determinism": true,
		"execution_timeout": 3600.5
	}`

	var def CodeDefinition
	if err := json.Unmarshal([]byte(raw), &def); err != nil {
		t.Fatal(err)
	}

	if def.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}
	if def.OutputSchema == nil {
		t.Fatal("OutputSchema is nil")
	}
	if len(def.Signals) != 1 {
		t.Fatalf("expected 1 signal, got %d", len(def.Signals))
	}
	if def.Signals[0].Name != "cancel" {
		t.Errorf("signal name = %q, want cancel", def.Signals[0].Name)
	}
	if def.Signals[0].Description == nil || *def.Signals[0].Description != "Cancel the workflow" {
		t.Errorf("signal description wrong")
	}
	if len(def.Queries) != 1 {
		t.Fatalf("expected 1 query, got %d", len(def.Queries))
	}
	if def.Queries[0].Name != "status" {
		t.Errorf("query name = %q, want status", def.Queries[0].Name)
	}
	if def.Queries[0].OutputSchema == nil {
		t.Error("query OutputSchema is nil, expected non-nil")
	}
	if len(def.Updates) != 1 {
		t.Fatalf("expected 1 update, got %d", len(def.Updates))
	}
	if def.Updates[0].Name != "set_priority" {
		t.Errorf("update name = %q, want set_priority", def.Updates[0].Name)
	}
	if def.EnforceDeterminism != true {
		t.Error("EnforceDeterminism should be true")
	}
	if def.ExecutionTimeout == nil || *def.ExecutionTimeout != 3600.5 {
		t.Errorf("ExecutionTimeout = %v, want 3600.5", def.ExecutionTimeout)
	}

	// Re-marshal and verify round-trip
	out, err := json.Marshal(def)
	if err != nil {
		t.Fatal(err)
	}
	var def2 CodeDefinition
	if err := json.Unmarshal(out, &def2); err != nil {
		t.Fatal(err)
	}
	if len(def2.Signals) != 1 || def2.Signals[0].Name != "cancel" {
		t.Error("round-trip failed for signals")
	}
	if def2.EnforceDeterminism != true {
		t.Error("round-trip failed for enforce_determinism")
	}
}

func TestCodeDefinition_MinimalFields(t *testing.T) {
	raw := `{"input_schema": {"type": "object"}}`

	var def CodeDefinition
	if err := json.Unmarshal([]byte(raw), &def); err != nil {
		t.Fatal(err)
	}
	if def.InputSchema == nil {
		t.Fatal("InputSchema is nil")
	}
	if def.OutputSchema != nil {
		t.Errorf("OutputSchema should be nil, got %v", def.OutputSchema)
	}
	if def.Signals != nil {
		t.Errorf("Signals should be nil, got %v", def.Signals)
	}
	if def.EnforceDeterminism != false {
		t.Error("EnforceDeterminism should default to false")
	}
	if def.ExecutionTimeout != nil {
		t.Errorf("ExecutionTimeout should be nil, got %v", def.ExecutionTimeout)
	}
}

func TestRegistration_NewFields(t *testing.T) {
	raw := `{
		"id": "reg-1",
		"workflow_id": "wf-1",
		"task_queue": "legacy-queue",
		"deployment_id": "dep-abc",
		"compatible_with_chat_assistant": true,
		"definition": {
			"input_schema": {"type": "object"},
			"enforce_determinism": false
		},
		"created_at": "2026-04-01T00:00:00Z",
		"updated_at": "2026-04-02T00:00:00Z"
	}`

	var reg Registration
	if err := json.Unmarshal([]byte(raw), &reg); err != nil {
		t.Fatal(err)
	}
	if reg.ID != "reg-1" {
		t.Errorf("ID = %q", reg.ID)
	}
	if reg.DeploymentID == nil || *reg.DeploymentID != "dep-abc" {
		t.Errorf("DeploymentID = %v, want dep-abc", reg.DeploymentID)
	}
	if reg.CompatibleWithChatAssistant != true {
		t.Error("CompatibleWithChatAssistant should be true")
	}
	if reg.Definition == nil {
		t.Fatal("Definition is nil")
	}
	if reg.Definition.InputSchema == nil {
		t.Error("Definition.InputSchema is nil")
	}
	// TaskQueue still works for backward compat
	if reg.TaskQueue != "legacy-queue" {
		t.Errorf("TaskQueue = %q, want legacy-queue", reg.TaskQueue)
	}
}

func TestRegistration_NullDeploymentID(t *testing.T) {
	raw := `{
		"id": "reg-2",
		"workflow_id": "wf-2",
		"task_queue": "q",
		"definition": {"input_schema": {"type": "object"}}
	}`

	var reg Registration
	if err := json.Unmarshal([]byte(raw), &reg); err != nil {
		t.Fatal(err)
	}
	if reg.DeploymentID != nil {
		t.Errorf("DeploymentID should be nil, got %v", reg.DeploymentID)
	}
	if reg.CompatibleWithChatAssistant != false {
		t.Error("CompatibleWithChatAssistant should default to false")
	}
}
