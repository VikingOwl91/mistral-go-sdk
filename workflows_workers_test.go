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
		if r.URL.Path != "/v1/workflows/workers/whoami" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"scheduler_url": "https://scheduler.mistral.ai",
			"namespace":     "default",
			"tls":           true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	info, err := client.GetWorkflowWorkerInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Namespace != "default" {
		t.Errorf("got namespace %q", info.Namespace)
	}
	if !info.TLS {
		t.Error("expected tls=true")
	}
}
