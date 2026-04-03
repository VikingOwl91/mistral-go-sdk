package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VikingOwl91/mistral-go-sdk/audio"
)

// Transcribe sends an audio file for transcription.
// If file is non-nil, it's uploaded via multipart. Otherwise, FileURL or FileID must be set.
func (c *Client) Transcribe(ctx context.Context, req *audio.TranscriptionRequest, filename string, file io.Reader) (*audio.TranscriptionResponse, error) {
	if file != nil {
		fields := map[string]string{"model": req.Model}
		if req.Language != nil {
			fields["language"] = *req.Language
		}
		if req.FileID != nil {
			fields["file_id"] = *req.FileID
		}
		if req.Diarize {
			fields["diarize"] = "true"
		}
		fields["stream"] = "false"

		var resp audio.TranscriptionResponse
		if err := c.doMultipart(ctx, "/v1/audio/transcriptions", filename, file, fields, &resp); err != nil {
			return nil, err
		}
		return &resp, nil
	}

	// JSON body request (file_url or file_id)
	var resp audio.TranscriptionResponse
	if err := c.doJSON(ctx, "POST", "/v1/audio/transcriptions", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TranscribeStream sends an audio file for streaming transcription.
func (c *Client) TranscribeStream(ctx context.Context, req *audio.TranscriptionRequest, filename string, file io.Reader) (*AudioStream, error) {
	fields := map[string]string{"model": req.Model, "stream": "true"}
	if req.Language != nil {
		fields["language"] = *req.Language
	}
	if req.Diarize {
		fields["diarize"] = "true"
	}

	resp, err := c.doMultipartStream(ctx, "/v1/audio/transcriptions", filename, file, fields)
	if err != nil {
		return nil, err
	}
	return newAudioStream(resp.Body), nil
}

// AudioStream wraps the generic Stream to provide typed audio events.
type AudioStream struct {
	stream *Stream[json.RawMessage]
	event  audio.StreamEvent
	err    error
}

func newAudioStream(body readCloser) *AudioStream {
	return &AudioStream{
		stream: newStream[json.RawMessage](body),
	}
}

// Next advances to the next event. Returns false when done or on error.
func (s *AudioStream) Next() bool {
	if s.err != nil {
		return false
	}
	if !s.stream.Next() {
		s.err = s.stream.Err()
		return false
	}
	event, err := audio.UnmarshalStreamEvent(s.stream.Current())
	if err != nil {
		s.err = err
		return false
	}
	s.event = event
	return true
}

// Current returns the most recently read event.
func (s *AudioStream) Current() audio.StreamEvent { return s.event }

// Err returns any error encountered during streaming.
func (s *AudioStream) Err() error { return s.err }

// Close releases the underlying connection.
func (s *AudioStream) Close() error { return s.stream.Close() }

// Speech sends a text-to-speech request and returns the full response.
func (c *Client) Speech(ctx context.Context, req *audio.SpeechRequest) (*audio.SpeechResponse, error) {
	var resp audio.SpeechResponse
	if err := c.doJSON(ctx, "POST", "/v1/audio/speech", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SpeechStream sends a text-to-speech request and returns a stream of audio events.
func (c *Client) SpeechStream(ctx context.Context, req *audio.SpeechRequest) (*SpeechStream, error) {
	req.EnableStream()
	resp, err := c.doStream(ctx, "POST", "/v1/audio/speech", req)
	if err != nil {
		return nil, err
	}
	return newSpeechStream(resp.Body), nil
}

// SpeechStream wraps the generic Stream for speech streaming events.
type SpeechStream struct {
	stream *Stream[json.RawMessage]
	event  audio.SpeechStreamEvent
	err    error
}

func newSpeechStream(body readCloser) *SpeechStream {
	return &SpeechStream{
		stream: newStream[json.RawMessage](body),
	}
}

// Next advances to the next event. Returns false when done or on error.
func (s *SpeechStream) Next() bool {
	if s.err != nil {
		return false
	}
	if !s.stream.Next() {
		s.err = s.stream.Err()
		return false
	}
	event, err := audio.UnmarshalSpeechStreamEvent(s.stream.Current())
	if err != nil {
		s.err = err
		return false
	}
	s.event = event
	return true
}

// Current returns the most recently read event.
func (s *SpeechStream) Current() audio.SpeechStreamEvent { return s.event }

// Err returns any error encountered during streaming.
func (s *SpeechStream) Err() error { return s.err }

// Close releases the underlying connection.
func (s *SpeechStream) Close() error { return s.stream.Close() }

// ListVoices returns available voices.
func (c *Client) ListVoices(ctx context.Context) (*audio.VoiceListResponse, error) {
	var resp audio.VoiceListResponse
	if err := c.doJSON(ctx, "GET", "/v1/audio/voices", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateVoice creates a custom voice.
func (c *Client) CreateVoice(ctx context.Context, req *audio.VoiceCreateRequest) (*audio.VoiceResponse, error) {
	var resp audio.VoiceResponse
	if err := c.doJSON(ctx, "POST", "/v1/audio/voices", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVoice retrieves a voice by ID.
func (c *Client) GetVoice(ctx context.Context, voiceID string) (*audio.VoiceResponse, error) {
	var resp audio.VoiceResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/audio/voices/%s", voiceID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateVoice updates a voice.
func (c *Client) UpdateVoice(ctx context.Context, voiceID string, req *audio.VoiceUpdateRequest) (*audio.VoiceResponse, error) {
	var resp audio.VoiceResponse
	if err := c.doJSON(ctx, "PATCH", fmt.Sprintf("/v1/audio/voices/%s", voiceID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteVoice deletes a voice.
func (c *Client) DeleteVoice(ctx context.Context, voiceID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/audio/voices/%s", voiceID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// GetVoiceSampleAudio retrieves the sample audio for a voice.
// Returns the raw HTTP response; the caller must close the body.
func (c *Client) GetVoiceSampleAudio(ctx context.Context, voiceID string) (*http.Response, error) {
	resp, err := c.do(ctx, "GET", fmt.Sprintf("/v1/audio/voices/%s/sample", voiceID), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, parseAPIError(resp)
	}
	return resp, nil
}
