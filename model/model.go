package model

// ModelCard represents a model (base or fine-tuned).
// For fine-tuned models (Type == "fine-tuned"), Job, Root, and Archived
// are populated.
type ModelCard struct {
	ID                          string            `json:"id"`
	Object                      string            `json:"object"`
	Created                     int64             `json:"created"`
	OwnedBy                     string            `json:"owned_by"`
	Capabilities                ModelCapabilities  `json:"capabilities"`
	Name                        *string           `json:"name,omitempty"`
	Description                 *string           `json:"description,omitempty"`
	MaxContextLength            int               `json:"max_context_length"`
	Aliases                     []string          `json:"aliases"`
	Deprecation                 *string           `json:"deprecation,omitempty"`
	DeprecationReplacementModel *string           `json:"deprecation_replacement_model,omitempty"`
	DefaultModelTemperature     *float64          `json:"default_model_temperature,omitempty"`
	Type                        string            `json:"type"`

	// Fine-tuned model fields (only when Type == "fine-tuned")
	Job      string `json:"job,omitempty"`
	Root     string `json:"root,omitempty"`
	Archived bool   `json:"archived,omitempty"`
}

// ModelCapabilities describes what a model can do.
type ModelCapabilities struct {
	CompletionChat     bool `json:"completion_chat"`
	FunctionCalling    bool `json:"function_calling"`
	CompletionFIM      bool `json:"completion_fim"`
	FineTuning         bool `json:"fine_tuning"`
	Vision             bool `json:"vision"`
	OCR                bool `json:"ocr"`
	Classification     bool `json:"classification"`
	Moderation         bool `json:"moderation"`
	Audio                      bool `json:"audio"`
	AudioTranscription         bool `json:"audio_transcription"`
	AudioTranscriptionRealtime bool `json:"audio_transcription_realtime"`
	AudioSpeech                bool `json:"audio_speech"`
}

// ModelList is the response from listing models.
type ModelList struct {
	Object string      `json:"object"`
	Data   []ModelCard `json:"data"`
}

// DeleteModelOut is the response from deleting a model.
type DeleteModelOut struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// ListParams holds optional parameters for listing models.
type ListParams struct {
	Provider *string
	Model    *string
}
