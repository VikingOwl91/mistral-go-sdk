package mistral

import (
	"encoding/json"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/conversation"
	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestNetworkEncodedInput_EncodingOptions(t *testing.T) {
	in := workflow.NetworkEncodedInput{
		B64Payload:      "eyJrIjoidiJ9",
		EncodingOptions: []workflow.EncodedPayloadOption{workflow.EncodedPayloadOffloaded, workflow.EncodedPayloadEncrypted},
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	opts, ok := got["encoding_options"].([]any)
	if !ok || len(opts) != 2 {
		t.Fatalf("unexpected encoding_options: %v", got["encoding_options"])
	}
	if opts[0] != "offloaded" || opts[1] != "encrypted" {
		t.Errorf("got %v, want [offloaded encrypted]", opts)
	}
}

func TestBuildConnectorExtensions_WireShape(t *testing.T) {
	creds := "work-account"
	ext := workflow.BuildConnectorExtensions(
		workflow.ConnectorSlot{ConnectorName: "gmail"},
		workflow.ConnectorSlot{ConnectorName: "notion", CredentialsName: &creds},
	)
	b, err := json.Marshal(ext)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"mistralai":{"connectors":{"bindings":[{"connector_name":"gmail"},{"connector_name":"notion","credentials_name":"work-account"}]}}}`
	if string(b) != want {
		t.Errorf("\nwant %s\ngot  %s", want, string(b))
	}
}

func TestExecutionRequest_Extensions(t *testing.T) {
	req := workflow.ExecutionRequest{
		Extensions: workflow.BuildConnectorExtensions(workflow.ConnectorSlot{ConnectorName: "gmail"}),
	}
	b, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["extensions"]; !ok {
		t.Fatalf("expected extensions key in marshalled request: %s", string(b))
	}
}

func TestConnectorAuthTaskState_Roundtrip(t *testing.T) {
	authURL := "https://oauth.example.com/authorize"
	in := workflow.ConnectorAuthTaskState{
		ConnectorName: "gmail",
		ConnectorID:   "conn-1",
		Status:        workflow.ConnectorAuthWaitingForAuth,
		AuthURL:       &authURL,
	}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out workflow.ConnectorAuthTaskState
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if out.Status != workflow.ConnectorAuthWaitingForAuth {
		t.Errorf("got status %q", out.Status)
	}
	if out.AuthURL == nil || *out.AuthURL != authURL {
		t.Errorf("auth_url roundtrip failed")
	}
}

func TestConfirmationConstants_WireValues(t *testing.T) {
	// Reply constants.
	c := conversation.ToolCallConfirmation{
		ToolCallID:   "call_1",
		Confirmation: string(conversation.ConfirmationAllow),
	}
	b, _ := json.Marshal(c)
	if string(b) != `{"tool_call_id":"call_1","confirmation":"allow"}` {
		t.Errorf("got %s", string(b))
	}
	// Inbound status constants.
	if conversation.ConfirmationStatusPending != "pending" ||
		conversation.ConfirmationStatusAllowed != "allowed" ||
		conversation.ConfirmationStatusDenied != "denied" {
		t.Errorf("unexpected confirmation status constants")
	}
}
