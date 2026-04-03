package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestListWorkflowEvents_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/events/list" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("got limit %q", r.URL.Query().Get("limit"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"events":      []map[string]any{{"event_type": "WORKFLOW_EXECUTION_STARTED"}},
			"next_cursor": "cur-1",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	limit := 50
	resp, err := client.ListWorkflowEvents(context.Background(), &workflow.EventListParams{Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) != 1 {
		t.Fatalf("got %d events", len(resp.Events))
	}
	if resp.NextCursor == nil || *resp.NextCursor != "cur-1" {
		t.Errorf("got cursor %v", resp.NextCursor)
	}
}
