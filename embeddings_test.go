package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/embedding"
)

func TestCreateEmbeddings_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("expected /v1/embeddings, got %s", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["model"] != "mistral-embed" {
			t.Errorf("expected model=mistral-embed, got %v", body["model"])
		}
		inputs := body["input"].([]any)
		if len(inputs) != 2 {
			t.Errorf("expected 2 inputs, got %d", len(inputs))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":     "emb-1",
			"object": "list",
			"model":  "mistral-embed",
			"usage":  map[string]any{"prompt_tokens": 12, "completion_tokens": 0, "total_tokens": 12},
			"data": []map[string]any{
				{"object": "embedding", "embedding": []float64{0.1, 0.2, 0.3}, "index": 0},
				{"object": "embedding", "embedding": []float64{0.4, 0.5, 0.6}, "index": 1},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model: "mistral-embed",
		Input: []string{"Hello world", "Goodbye world"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "emb-1" {
		t.Errorf("got id %q", resp.ID)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("got %d embeddings, want 2", len(resp.Data))
	}
	if resp.Data[0].Index != 0 {
		t.Errorf("got index %d", resp.Data[0].Index)
	}
	if len(resp.Data[0].Embedding) != 3 {
		t.Fatalf("got %d dims, want 3", len(resp.Data[0].Embedding))
	}
	if resp.Data[0].Embedding[0] != 0.1 {
		t.Errorf("got embedding[0]=%f", resp.Data[0].Embedding[0])
	}
	if resp.Data[1].Embedding[2] != 0.6 {
		t.Errorf("got embedding[2]=%f", resp.Data[1].Embedding[2])
	}
	if resp.Usage.PromptTokens != 12 {
		t.Errorf("got prompt_tokens=%d", resp.Usage.PromptTokens)
	}
}

func TestCreateEmbeddings_SingleInput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		inputs := body["input"].([]any)
		if len(inputs) != 1 {
			t.Errorf("expected 1 input, got %d", len(inputs))
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "emb-2", "object": "list", "model": "mistral-embed",
			"usage": map[string]any{"prompt_tokens": 5, "completion_tokens": 0, "total_tokens": 5},
			"data": []map[string]any{
				{"object": "embedding", "embedding": []float64{0.1, 0.2}, "index": 0},
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model: "mistral-embed",
		Input: []string{"Just one"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Data) != 1 {
		t.Errorf("got %d, want 1", len(resp.Data))
	}
}

func TestCreateEmbeddings_WithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["output_dimension"] != float64(256) {
			t.Errorf("expected output_dimension=256, got %v", body["output_dimension"])
		}
		if body["output_dtype"] != "int8" {
			t.Errorf("expected output_dtype=int8, got %v", body["output_dtype"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id": "emb-3", "object": "list", "model": "m",
			"usage": map[string]any{"prompt_tokens": 0, "completion_tokens": 0, "total_tokens": 0},
			"data":  []map[string]any{{"object": "embedding", "embedding": []float64{1, 2}, "index": 0}},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	dim := 256
	dtype := embedding.DtypeInt8
	_, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model:           "m",
		Input:           []string{"test"},
		OutputDimension: &dim,
		OutputDtype:     &dtype,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateEmbeddings_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(map[string]any{"message": "rate limited"})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	_, err := client.CreateEmbeddings(context.Background(), &embedding.Request{
		Model: "m",
		Input: []string{"test"},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsRateLimit(err) {
		t.Errorf("expected rate limit, got: %v", err)
	}
}
