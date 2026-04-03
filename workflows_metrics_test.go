package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

func TestGetWorkflowMetrics_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/workflows/my-flow/metrics" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("start_time") != "2026-01-01T00:00:00Z" {
			t.Errorf("got start_time %q", r.URL.Query().Get("start_time"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"execution_count":    map[string]any{"value": 100},
			"success_count":      map[string]any{"value": 95},
			"error_count":        map[string]any{"value": 5},
			"average_latency_ms": map[string]any{"value": 1234.5},
			"latency_over_time":  map[string]any{"value": [][]float64{{1711929600, 1200}, {1711929660, 1300}}},
			"retry_rate":         map[string]any{"value": 0.02},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	start := "2026-01-01T00:00:00Z"
	resp, err := client.GetWorkflowMetrics(context.Background(), "my-flow", &workflow.MetricsParams{StartTime: &start})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ExecutionCount.Value != 100 {
		t.Errorf("got execution_count %v", resp.ExecutionCount.Value)
	}
	if resp.AverageLatencyMs.Value != 1234.5 {
		t.Errorf("got average_latency_ms %v", resp.AverageLatencyMs.Value)
	}
}
