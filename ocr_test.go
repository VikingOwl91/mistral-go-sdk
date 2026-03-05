package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/ocr"
)

func TestOCR_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/ocr" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"pages": []map[string]any{{
				"index":    0,
				"markdown": "# Hello World",
				"images":   []any{},
				"dimensions": map[string]any{
					"dpi": 200, "height": 2200, "width": 1700,
				},
			}},
			"model": "mistral-ocr-latest",
			"usage_info": map[string]any{
				"pages_processed": 1,
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	model := "mistral-ocr-latest"
	resp, err := client.OCR(context.Background(), &ocr.Request{
		Model:    &model,
		Document: json.RawMessage(`{"type":"document_url","document_url":"https://example.com/doc.pdf"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Model != "mistral-ocr-latest" {
		t.Errorf("got model %q", resp.Model)
	}
	if len(resp.Pages) != 1 {
		t.Fatalf("got %d pages", len(resp.Pages))
	}
	if resp.Pages[0].Markdown != "# Hello World" {
		t.Errorf("got markdown %q", resp.Pages[0].Markdown)
	}
	if resp.UsageInfo.PagesProcessed != 1 {
		t.Errorf("got pages_processed %d", resp.UsageInfo.PagesProcessed)
	}
}
