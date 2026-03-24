package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/audio"
)

func TestSpeech_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/audio/speech" {
			t.Errorf("got path %s", r.URL.Path)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["input"] != "Hello world" {
			t.Errorf("got input %v", body["input"])
		}
		if body["stream"] != false {
			t.Errorf("expected stream=false")
		}
		json.NewEncoder(w).Encode(map[string]any{
			"audio_data": "base64audiodata==",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.Speech(context.Background(), &audio.SpeechRequest{
		Input: "Hello world",
		Model: "mistral-speech",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.AudioData != "base64audiodata==" {
		t.Errorf("got audio_data %q", resp.AudioData)
	}
}

func TestSpeechStream_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		flusher, _ := w.(http.Flusher)

		delta, _ := json.Marshal(map[string]any{
			"type": "speech.audio.delta", "audio_data": "chunk1==",
		})
		fmt.Fprintf(w, "data: %s\n\n", delta)
		flusher.Flush()

		done, _ := json.Marshal(map[string]any{
			"type": "speech.audio.done",
			"usage": map[string]any{
				"prompt_tokens": 5, "completion_tokens": 10, "total_tokens": 15,
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", done)
		flusher.Flush()

		fmt.Fprint(w, "data: [DONE]\n\n")
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	stream, err := client.SpeechStream(context.Background(), &audio.SpeechRequest{
		Input: "Hi",
		Model: "mistral-speech",
	})
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var events int
	for stream.Next() {
		events++
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if events != 2 {
		t.Errorf("got %d events, want 2", events)
	}
}

func TestListVoices_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/voices" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{
				{"id": "v1", "name": "Default", "created_at": "2025-01-01"},
			},
			"total": 1, "page": 1, "page_size": 10, "total_pages": 1,
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.ListVoices(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("got %d voices", len(resp.Items))
	}
	if resp.Items[0].ID != "v1" {
		t.Errorf("got id %q", resp.Items[0].ID)
	}
}

func TestCreateVoice_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "MyVoice" {
			t.Errorf("got name %v", body["name"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "v2", "name": "MyVoice", "created_at": "2025-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.CreateVoice(context.Background(), &audio.VoiceCreateRequest{
		Name:        "MyVoice",
		SampleAudio: "base64audio==",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.ID != "v2" {
		t.Errorf("got id %q", resp.ID)
	}
}

func TestGetVoice_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/voices/v1" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "v1", "name": "Default", "created_at": "2025-01-01",
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.GetVoice(context.Background(), "v1")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "Default" {
		t.Errorf("got name %q", resp.Name)
	}
}

func TestUpdateVoice_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"id": "v1", "name": "Renamed", "created_at": "2025-01-01",
		})
	}))
	defer server.Close()

	name := "Renamed"
	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.UpdateVoice(context.Background(), "v1", &audio.VoiceUpdateRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Name != "Renamed" {
		t.Errorf("got name %q", resp.Name)
	}
}

func TestDeleteVoice_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	err := client.DeleteVoice(context.Background(), "v1")
	if err != nil {
		t.Fatal(err)
	}
}
