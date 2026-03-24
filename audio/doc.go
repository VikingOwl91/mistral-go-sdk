// Package audio provides types for the Mistral audio APIs.
//
// # Transcription
//
// [TranscriptionRequest] and [TranscriptionResponse] handle speech-to-text.
// Streaming transcription returns typed [StreamEvent] values via a sealed
// interface dispatched by the "type" field.
//
// # Speech (TTS)
//
// [SpeechRequest] and [SpeechResponse] handle text-to-speech.
// Streaming speech returns typed [SpeechStreamEvent] values
// ([SpeechAudioDelta] and [SpeechDone]).
//
// # Voices
//
// [VoiceResponse], [VoiceCreateRequest], and [VoiceUpdateRequest] manage
// custom voices for speech synthesis.
//
// # Realtime
//
// Realtime transcription types ([AudioEncoding], [AudioFormat],
// [RealtimeSession], and WebSocket message types) are defined here.
// The WebSocket client is not yet implemented.
package audio
