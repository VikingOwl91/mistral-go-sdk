package mistral

import (
	"context"
	"encoding/json"
	"io"

	"somegit.dev/vikingowl/mistral-go-sdk/audio"
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
