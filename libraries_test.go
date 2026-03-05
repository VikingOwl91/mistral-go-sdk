package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/library"
)

func newLibraryJSON() map[string]any {
	return map[string]any{
		"id": "lib-123", "name": "TestLib",
		"created_at": "2024-01-01T00:00:00Z", "updated_at": "2024-01-01T00:00:00Z",
		"owner_id": "user-1", "owner_type": "user",
		"total_size": 1024, "nb_documents": 5, "chunk_size": 512,
	}
}

func TestCreateLibrary_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newLibraryJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	lib, err := client.CreateLibrary(context.Background(), &library.CreateRequest{
		Name: "TestLib",
	})
	if err != nil {
		t.Fatal(err)
	}
	if lib.ID != "lib-123" {
		t.Errorf("got id %q", lib.ID)
	}
}

func TestListLibraries_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{newLibraryJSON()},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListLibraries(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("got %d libraries", len(resp.Data))
	}
}

func TestGetLibrary_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/libraries/lib-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(newLibraryJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	lib, err := client.GetLibrary(context.Background(), "lib-123")
	if err != nil {
		t.Fatal(err)
	}
	if lib.Name != "TestLib" {
		t.Errorf("got name %q", lib.Name)
	}
}

func TestDeleteLibrary_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		json.NewEncoder(w).Encode(newLibraryJSON())
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	lib, err := client.DeleteLibrary(context.Background(), "lib-123")
	if err != nil {
		t.Fatal(err)
	}
	if lib.ID != "lib-123" {
		t.Errorf("got id %q", lib.ID)
	}
}

func TestUploadDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart, got %q", ct)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id": "doc-1", "library_id": "lib-123", "name": "test.pdf",
			"hash": "abc", "mime_type": "application/pdf", "extension": ".pdf", "size": 1024,
			"created_at": "2024-01-01T00:00:00Z", "process_status": "todo",
			"uploaded_by_id": "user-1", "uploaded_by_type": "user",
			"processing_status": "todo", "tokens_processing_total": 0,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	doc, err := client.UploadDocument(context.Background(), "lib-123", "test.pdf", strings.NewReader("fake pdf"))
	if err != nil {
		t.Fatal(err)
	}
	if doc.ID != "doc-1" {
		t.Errorf("got id %q", doc.ID)
	}
}

func TestListDocuments_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/libraries/lib-123/documents" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"pagination": map[string]any{
				"total_items": 1, "total_pages": 1, "current_page": 0, "page_size": 100, "has_more": false,
			},
			"data": []map[string]any{{
				"id": "doc-1", "library_id": "lib-123", "name": "test.pdf",
				"hash": "abc", "mime_type": "application/pdf", "extension": ".pdf", "size": 1024,
				"created_at": "2024-01-01T00:00:00Z", "process_status": "done",
				"uploaded_by_id": "user-1", "uploaded_by_type": "user",
				"processing_status": "done", "tokens_processing_total": 500,
			}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListDocuments(context.Background(), "lib-123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("got %d documents", len(resp.Data))
	}
	if resp.Pagination.TotalItems != 1 {
		t.Errorf("got total_items %d", resp.Pagination.TotalItems)
	}
}

func TestGetDocumentTextContent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"text": "Hello world"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	content, err := client.GetDocumentTextContent(context.Background(), "lib-123", "doc-1")
	if err != nil {
		t.Fatal(err)
	}
	if content.Text != "Hello world" {
		t.Errorf("got text %q", content.Text)
	}
}

func TestGetDocumentStatus_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"document_id": "doc-1", "process_status": "done", "processing_status": "done",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	status, err := client.GetDocumentStatus(context.Background(), "lib-123", "doc-1")
	if err != nil {
		t.Fatal(err)
	}
	if status.ProcessStatus != "done" {
		t.Errorf("got status %q", status.ProcessStatus)
	}
}

func TestDeleteDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("got method %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteDocument(context.Background(), "lib-123", "doc-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestReprocessDocument_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("got method %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.ReprocessDocument(context.Background(), "lib-123", "doc-1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestListLibrarySharing_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{{
				"library_id": "lib-123", "org_id": "org-1",
				"role": "Viewer", "share_with_type": "User", "share_with_uuid": "user-2",
			}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListLibrarySharing(context.Background(), "lib-123")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("got %d sharing entries", len(resp.Data))
	}
	if resp.Data[0].Role != "Viewer" {
		t.Errorf("got role %q", resp.Data[0].Role)
	}
}
