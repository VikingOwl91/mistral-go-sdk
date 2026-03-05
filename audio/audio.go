package audio

import (
	"encoding/json"
	"fmt"
)

// TranscriptionRequest represents an audio transcription request.
// The file is sent via multipart form, not in this struct.
type TranscriptionRequest struct {
	Model                   string   `json:"model"`
	FileURL                 *string  `json:"file_url,omitempty"`
	FileID                  *string  `json:"file_id,omitempty"`
	Language                *string  `json:"language,omitempty"`
	Temperature             *float64 `json:"temperature,omitempty"`
	Diarize                 bool     `json:"diarize,omitempty"`
	ContextBias             []string `json:"context_bias,omitempty"`
	TimestampGranularities  []string `json:"timestamp_granularities,omitempty"`
}

// TranscriptionResponse is the result of a transcription request.
type TranscriptionResponse struct {
	Model    string    `json:"model"`
	Text     string    `json:"text"`
	Language *string   `json:"language"`
	Segments []Segment `json:"segments,omitempty"`
	Usage    UsageInfo `json:"usage"`
}

// Segment represents a transcription segment with timestamps.
type Segment struct {
	Type      string  `json:"type,omitempty"`
	Text      string  `json:"text"`
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	Score     *float64 `json:"score,omitempty"`
	SpeakerID *string  `json:"speaker_id,omitempty"`
}

// UsageInfo for audio transcription.
type UsageInfo struct {
	PromptTokens       int  `json:"prompt_tokens"`
	CompletionTokens   int  `json:"completion_tokens"`
	TotalTokens        int  `json:"total_tokens"`
	PromptAudioSeconds *int `json:"prompt_audio_seconds,omitempty"`
}

// StreamEvent is a discriminated union for transcription stream events.
type StreamEvent interface {
	streamEvent()
}

// LanguageEvent is emitted when the language is detected.
type LanguageEvent struct {
	Type          string `json:"type"`
	AudioLanguage string `json:"audio_language"`
}

func (LanguageEvent) streamEvent() {}

// TextDeltaEvent contains a text chunk during streaming.
type TextDeltaEvent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (TextDeltaEvent) streamEvent() {}

// SegmentDeltaEvent contains a segment during streaming.
type SegmentDeltaEvent struct {
	Type      string  `json:"type"`
	Text      string  `json:"text"`
	Start     float64 `json:"start"`
	End       float64 `json:"end"`
	SpeakerID *string `json:"speaker_id,omitempty"`
}

func (SegmentDeltaEvent) streamEvent() {}

// DoneEvent is emitted when transcription is complete.
type DoneEvent struct {
	Type     string    `json:"type"`
	Model    string    `json:"model"`
	Text     string    `json:"text"`
	Language *string   `json:"language"`
	Segments []Segment `json:"segments,omitempty"`
	Usage    UsageInfo `json:"usage"`
}

func (DoneEvent) streamEvent() {}

// UnmarshalStreamEvent dispatches a raw JSON event to the correct type.
func UnmarshalStreamEvent(data []byte) (StreamEvent, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	switch probe.Type {
	case "transcription.language":
		var e LanguageEvent
		return &e, json.Unmarshal(data, &e)
	case "transcription.text.delta":
		var e TextDeltaEvent
		return &e, json.Unmarshal(data, &e)
	case "transcription.segment":
		var e SegmentDeltaEvent
		return &e, json.Unmarshal(data, &e)
	case "transcription.done":
		var e DoneEvent
		return &e, json.Unmarshal(data, &e)
	default:
		return nil, fmt.Errorf("unknown transcription stream event type: %q", probe.Type)
	}
}
