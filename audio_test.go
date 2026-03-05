package mistral

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"somegit.dev/vikingowl/mistral-go-sdk/audio"
)

func TestTranscribe_WithFileURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/transcriptions" {
			t.Errorf("got path %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"model": "voxtral-mini-latest", "text": "Hello world",
			"language": "en", "segments": []any{},
			"usage": map[string]any{
				"prompt_tokens": 4, "completion_tokens": 10,
				"total_tokens": 14, "prompt_audio_seconds": 5,
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	url := "https://example.com/audio.mp3"
	resp, err := client.Transcribe(context.Background(), &audio.TranscriptionRequest{
		Model:   "voxtral-mini-latest",
		FileURL: &url,
	}, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "Hello world" {
		t.Errorf("got text %q", resp.Text)
	}
	if resp.Usage.TotalTokens != 14 {
		t.Errorf("got total_tokens %d", resp.Usage.TotalTokens)
	}
}

func TestTranscribe_WithFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/audio/transcriptions" {
			t.Errorf("got path %s", r.URL.Path)
		}
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, "multipart/form-data") {
			t.Errorf("expected multipart, got %q", ct)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"model": "voxtral-mini-latest", "text": "Test",
			"language": "en", "segments": []any{},
			"usage": map[string]any{
				"prompt_tokens": 2, "completion_tokens": 5, "total_tokens": 7,
			},
		})
	}))
	defer server.Close()

	client := NewClient("key", WithBaseURL(server.URL))
	resp, err := client.Transcribe(context.Background(), &audio.TranscriptionRequest{
		Model: "voxtral-mini-latest",
	}, "test.mp3", strings.NewReader("fake audio data"))
	if err != nil {
		t.Fatal(err)
	}
	if resp.Text != "Test" {
		t.Errorf("got text %q", resp.Text)
	}
}

func TestUnmarshalStreamEvent_Language(t *testing.T) {
	data := []byte(`{"type":"transcription.language","audio_language":"en"}`)
	event, err := audio.UnmarshalStreamEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*audio.LanguageEvent)
	if !ok {
		t.Fatalf("expected *LanguageEvent, got %T", event)
	}
	if e.AudioLanguage != "en" {
		t.Errorf("got %q", e.AudioLanguage)
	}
}

func TestUnmarshalStreamEvent_TextDelta(t *testing.T) {
	data := []byte(`{"type":"transcription.text.delta","text":"Hello"}`)
	event, err := audio.UnmarshalStreamEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*audio.TextDeltaEvent)
	if !ok {
		t.Fatalf("expected *TextDeltaEvent, got %T", event)
	}
	if e.Text != "Hello" {
		t.Errorf("got %q", e.Text)
	}
}

func TestUnmarshalStreamEvent_Segment(t *testing.T) {
	data := []byte(`{"type":"transcription.segment","text":"Hello world","start":0.0,"end":1.5}`)
	event, err := audio.UnmarshalStreamEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*audio.SegmentDeltaEvent)
	if !ok {
		t.Fatalf("expected *SegmentDeltaEvent, got %T", event)
	}
	if e.End != 1.5 {
		t.Errorf("got end %f", e.End)
	}
}

func TestUnmarshalStreamEvent_Done(t *testing.T) {
	data := []byte(`{"type":"transcription.done","model":"voxtral-mini-latest","text":"Full text","language":"en","segments":[],"usage":{"prompt_tokens":4,"completion_tokens":10,"total_tokens":14}}`)
	event, err := audio.UnmarshalStreamEvent(data)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := event.(*audio.DoneEvent)
	if !ok {
		t.Fatalf("expected *DoneEvent, got %T", event)
	}
	if e.Usage.TotalTokens != 14 {
		t.Errorf("got total_tokens %d", e.Usage.TotalTokens)
	}
}

func TestUnmarshalStreamEvent_Unknown(t *testing.T) {
	_, err := audio.UnmarshalStreamEvent([]byte(`{"type":"unknown.event"}`))
	if err == nil {
		t.Error("expected error for unknown type")
	}
}
