package audio

// AudioEncoding is the encoding format for realtime audio streams.
type AudioEncoding string

const (
	EncodingPCMS16LE AudioEncoding = "pcm_s16le"
	EncodingPCMS32LE AudioEncoding = "pcm_s32le"
	EncodingPCMF16LE AudioEncoding = "pcm_f16le"
	EncodingPCMF32LE AudioEncoding = "pcm_f32le"
	EncodingPCMMulaw AudioEncoding = "pcm_mulaw"
	EncodingPCMAlaw  AudioEncoding = "pcm_alaw"
)

// AudioFormat describes the encoding and sample rate for realtime audio.
type AudioFormat struct {
	Encoding   AudioEncoding `json:"encoding"`
	SampleRate int           `json:"sample_rate"`
}

// RealtimeSession describes a realtime transcription session.
type RealtimeSession struct {
	RequestID              string       `json:"request_id"`
	Model                  string       `json:"model"`
	AudioFormat            AudioFormat  `json:"audio_format"`
	TargetStreamingDelayMs *int         `json:"target_streaming_delay_ms,omitempty"`
}

// RealtimeSessionUpdate is sent to update session parameters.
// Parameters can only be changed before audio transmission starts.
type RealtimeSessionUpdate struct {
	AudioFormat            *AudioFormat `json:"audio_format,omitempty"`
	TargetStreamingDelayMs *int         `json:"target_streaming_delay_ms,omitempty"`
}

// InputAudioAppend sends a chunk of audio data.
// Audio is base64-encoded (max 262144 bytes decoded).
type InputAudioAppend struct {
	Type  string `json:"type"` // "input_audio.append"
	Audio string `json:"audio"`
}

// InputAudioFlush flushes the audio buffer.
type InputAudioFlush struct {
	Type string `json:"type"` // "input_audio.flush"
}

// InputAudioEnd signals the end of audio input.
type InputAudioEnd struct {
	Type string `json:"type"` // "input_audio.end"
}

// RealtimeSessionCreated is received when a session is created.
type RealtimeSessionCreated struct {
	Type    string          `json:"type"` // "session.created"
	Session RealtimeSession `json:"session"`
}

// RealtimeSessionUpdated is received when a session is updated.
type RealtimeSessionUpdated struct {
	Type    string          `json:"type"` // "session.updated"
	Session RealtimeSession `json:"session"`
}

// RealtimeErrorDetail describes a realtime error.
type RealtimeErrorDetail struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// RealtimeError is received on error.
type RealtimeError struct {
	Type  string              `json:"type"` // "error"
	Error RealtimeErrorDetail `json:"error"`
}
