package mistral

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VikingOwl91/mistral-go-sdk/file"
)

func TestUploadFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/files" {
			t.Errorf("got path %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart, got %s", ct)
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatal(err)
		}
		f, header, err := r.FormFile("file")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if header.Filename != "train.jsonl" {
			t.Errorf("got filename %q", header.Filename)
		}
		content, _ := io.ReadAll(f)
		if string(content) != `{"text":"hello"}` {
			t.Errorf("got content %q", content)
		}
		if r.FormValue("purpose") != "fine-tune" {
			t.Errorf("got purpose %q", r.FormValue("purpose"))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "file-abc123", "object": "file",
			"bytes": 16, "created_at": 1700000000,
			"filename": "train.jsonl", "purpose": "fine-tune",
			"sample_type": "instruct", "source": "upload",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.UploadFile(
		context.Background(),
		"train.jsonl",
		strings.NewReader(`{"text":"hello"}`),
		&file.UploadParams{Purpose: file.PurposeFineTune},
	)
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "file-abc123" {
		t.Errorf("got id %q", resp.ID)
	}
	if resp.Filename != "train.jsonl" {
		t.Errorf("got filename %q", resp.Filename)
	}
	if resp.Purpose != file.PurposeFineTune {
		t.Errorf("got purpose %q", resp.Purpose)
	}
}

func TestUploadFile_WithExpiryAndVisibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Fatal(err)
		}
		if r.FormValue("purpose") != "fine-tune" {
			t.Errorf("got purpose %q", r.FormValue("purpose"))
		}
		if r.FormValue("expiry") != "48" {
			t.Errorf("expected expiry=48, got %q", r.FormValue("expiry"))
		}
		if r.FormValue("visibility") != "private" {
			t.Errorf("expected visibility=private, got %q", r.FormValue("visibility"))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "file-ev", "object": "file", "bytes": 10,
			"created_at": 1, "filename": "data.jsonl",
			"purpose": "fine-tune", "sample_type": "instruct",
			"source": "upload",
		})
	}))
	defer server.Close()

	expiry := 48
	vis := file.VisibilityPrivate
	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.UploadFile(context.Background(), "data.jsonl", strings.NewReader("{}"), &file.UploadParams{
		Purpose:    file.PurposeFineTune,
		Expiry:     &expiry,
		Visibility: &vis,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListFiles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/files" {
			t.Errorf("got path %s", r.URL.Path)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"object": "list",
			"total":  2,
			"data": []map[string]any{
				{"id": "f1", "object": "file", "bytes": 100, "created_at": 1, "filename": "a.jsonl", "purpose": "fine-tune", "sample_type": "instruct", "source": "upload"},
				{"id": "f2", "object": "file", "bytes": 200, "created_at": 2, "filename": "b.jsonl", "purpose": "batch", "sample_type": "batch_request", "source": "upload"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListFiles(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("got %d files, want 2", len(resp.Data))
	}
	if resp.Data[0].ID != "f1" {
		t.Errorf("got id %q", resp.Data[0].ID)
	}
	if resp.Total == nil || *resp.Total != 2 {
		t.Errorf("got total %v", resp.Total)
	}
}

func TestListFiles_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("page") != "1" {
			t.Errorf("expected page=1, got %q", q.Get("page"))
		}
		if q.Get("page_size") != "10" {
			t.Errorf("expected page_size=10, got %q", q.Get("page_size"))
		}
		if q.Get("purpose") != "fine-tune" {
			t.Errorf("expected purpose=fine-tune, got %q", q.Get("purpose"))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"object": "list", "data": []any{},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	page, pageSize := 1, 10
	purpose := file.PurposeFineTune
	_, err := client.ListFiles(context.Background(), &file.ListParams{
		Page:     &page,
		PageSize: &pageSize,
		Purpose:  &purpose,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/files/file-123" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "file-123", "object": "file", "bytes": 100,
			"created_at": 1, "filename": "test.jsonl",
			"purpose": "fine-tune", "sample_type": "instruct",
			"source": "upload", "deleted": false,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	f, err := client.GetFile(context.Background(), "file-123")
	if err != nil {
		t.Fatal(err)
	}
	if f.ID != "file-123" {
		t.Errorf("got id %q", f.ID)
	}
	if f.Deleted {
		t.Error("expected deleted=false")
	}
}

func TestDeleteFile_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "file-123", "object": "file", "deleted": true,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.DeleteFile(context.Background(), "file-123")
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Deleted {
		t.Error("expected deleted=true")
	}
}

func TestGetFileContent_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/files/file-123/content" {
			t.Errorf("got path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte(`{"text":"training data"}`))
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	body, err := client.GetFileContent(context.Background(), "file-123")
	if err != nil {
		t.Fatal(err)
	}
	defer body.Close()
	data, _ := io.ReadAll(body)
	if string(data) != `{"text":"training data"}` {
		t.Errorf("got %q", data)
	}
}

func TestGetFileURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/files/file-123/url" {
			t.Errorf("got path %s", r.URL.Path)
		}
		if r.URL.Query().Get("expiry") != "48" {
			t.Errorf("expected expiry=48, got %q", r.URL.Query().Get("expiry"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"url": "https://storage.example.com/file-123?token=abc",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetFileURL(context.Background(), "file-123", 48)
	if err != nil {
		t.Fatal(err)
	}
	if resp.URL != "https://storage.example.com/file-123?token=abc" {
		t.Errorf("got url %q", resp.URL)
	}
}

func TestGetFileURL_DefaultExpiry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("expiry") != "" {
			t.Errorf("expected no expiry param, got %q", r.URL.Query().Get("expiry"))
		}
		json.NewEncoder(w).Encode(map[string]any{"url": "https://example.com/f"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.GetFileURL(context.Background(), "file-123", 0)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]any{"message": "invalid file"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.UploadFile(context.Background(), "bad.txt", strings.NewReader(""), &file.UploadParams{Purpose: file.PurposeFineTune})
	if err == nil {
		t.Fatal("expected error")
	}
}
