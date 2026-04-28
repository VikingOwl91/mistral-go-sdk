package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWorkflowWorkerInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/workflows/workers/whoami" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"scheduler_url": "scheduler.example.com:7233",
			"namespace":     "tenant-2",
			"tls":           true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	info, err := client.GetWorkflowWorkerInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.SchedulerURL != "scheduler.example.com:7233" || info.Namespace != "tenant-2" || !info.TLS {
		t.Errorf("unexpected info: %+v", info)
	}
}

func TestGetWorkflowWorkerInfo_TLSDefault(t *testing.T) {
	// Server omits the tls field; the SDK should default it to false.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"scheduler_url": "s",
			"namespace":     "n",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	info, err := client.GetWorkflowWorkerInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.TLS {
		t.Errorf("expected default tls=false, got true")
	}
}
