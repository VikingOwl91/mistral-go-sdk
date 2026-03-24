package audio

import (
	"encoding/json"
	"fmt"
)

// SpeechOutputFormat is the output audio format for speech synthesis.
type SpeechOutputFormat string

const (
	SpeechFormatPCM  SpeechOutputFormat = "pcm"
	SpeechFormatWAV  SpeechOutputFormat = "wav"
	SpeechFormatMP3  SpeechOutputFormat = "mp3"
	SpeechFormatFLAC SpeechOutputFormat = "flac"
	SpeechFormatOpus SpeechOutputFormat = "opus"
)

// SpeechRequest represents a text-to-speech request.
type SpeechRequest struct {
	Input          string              `json:"input"`
	Model          string              `json:"model"`
	Metadata       map[string]any      `json:"metadata,omitempty"`
	VoiceID        *string             `json:"voice_id,omitempty"`
	RefAudio       *string             `json:"ref_audio,omitempty"`
	ResponseFormat *SpeechOutputFormat `json:"response_format,omitempty"`
	stream         bool
}

// EnableStream is used internally to enable streaming.
func (r *SpeechRequest) EnableStream() { r.stream = true }

func (r *SpeechRequest) MarshalJSON() ([]byte, error) {
	type Alias SpeechRequest
	return json.Marshal(&struct {
		Stream bool `json:"stream"`
		*Alias
	}{
		Stream: r.stream,
		Alias:  (*Alias)(r),
	})
}

// SpeechResponse is the response from a non-streaming speech request.
type SpeechResponse struct {
	AudioData string `json:"audio_data"`
}

// SpeechStreamEvent is a sealed interface for speech streaming events.
type SpeechStreamEvent interface {
	speechStreamEvent()
}

// SpeechAudioDelta contains a chunk of audio data during streaming.
type SpeechAudioDelta struct {
	Type      string `json:"type"`
	AudioData string `json:"audio_data"`
}

func (*SpeechAudioDelta) speechStreamEvent() {}

// SpeechDone is emitted when speech synthesis is complete.
type SpeechDone struct {
	Type  string    `json:"type"`
	Usage UsageInfo `json:"usage"`
}

func (*SpeechDone) speechStreamEvent() {}

// UnmarshalSpeechStreamEvent dispatches a raw JSON event to the correct type.
func UnmarshalSpeechStreamEvent(data []byte) (SpeechStreamEvent, error) {
	var probe struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}
	switch probe.Type {
	case "speech.audio.delta":
		var e SpeechAudioDelta
		return &e, json.Unmarshal(data, &e)
	case "speech.audio.done":
		var e SpeechDone
		return &e, json.Unmarshal(data, &e)
	default:
		return nil, fmt.Errorf("unknown speech stream event type: %q", probe.Type)
	}
}
