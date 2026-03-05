package finetune

import "encoding/json"

// JobIn is the request to create a fine-tuning job.
type JobIn struct {
	Model                       string          `json:"model"`
	Hyperparameters             json.RawMessage `json:"hyperparameters"`
	TrainingFiles               []TrainingFile  `json:"training_files,omitempty"`
	ValidationFiles             []string        `json:"validation_files,omitempty"`
	Suffix                      *string         `json:"suffix,omitempty"`
	Integrations                []WandbIntegration `json:"integrations,omitempty"`
	AutoStart                   bool            `json:"auto_start"`
	InvalidSampleSkipPercentage float64         `json:"invalid_sample_skip_percentage,omitempty"`
	JobType                     *string         `json:"job_type,omitempty"`
	Repositories                json.RawMessage `json:"repositories,omitempty"`
	ClassifierTargets           json.RawMessage `json:"classifier_targets,omitempty"`
}

// TrainingFile references an uploaded file with an optional weight.
type TrainingFile struct {
	FileID string  `json:"file_id"`
	Weight float64 `json:"weight,omitempty"`
}

// WandbIntegration configures Weights & Biases integration.
type WandbIntegration struct {
	Type    string  `json:"type"`
	Project string  `json:"project"`
	APIKey  string  `json:"api_key"`
	Name    *string `json:"name,omitempty"`
	RunName *string `json:"run_name,omitempty"`
}

// JobOut represents a fine-tuning job response.
// Fields are a superset across completion and classifier job types.
type JobOut struct {
	ID              string          `json:"id"`
	Object          string          `json:"object,omitempty"`
	Model           string          `json:"model"`
	Status          string          `json:"status"`
	JobType         string          `json:"job_type"`
	CreatedAt       int             `json:"created_at"`
	ModifiedAt      int             `json:"modified_at"`
	TrainingFiles   []string        `json:"training_files"`
	ValidationFiles []string        `json:"validation_files,omitempty"`
	FineTunedModel  *string         `json:"fine_tuned_model,omitempty"`
	Suffix          *string         `json:"suffix,omitempty"`
	Integrations    json.RawMessage `json:"integrations,omitempty"`
	TrainedTokens   *int            `json:"trained_tokens,omitempty"`
	Hyperparameters json.RawMessage `json:"hyperparameters,omitempty"`
	Repositories    json.RawMessage `json:"repositories,omitempty"`
	Metadata        *JobMetadata    `json:"metadata,omitempty"`
	AutoStart       bool            `json:"auto_start"`
	// DetailedJobOut fields (present when fetching individual job)
	Checkpoints json.RawMessage `json:"checkpoints,omitempty"`
	Events      json.RawMessage `json:"events,omitempty"`
}

// JobMetadata holds cost and timing estimates.
type JobMetadata struct {
	ExpectedDurationSeconds *int     `json:"expected_duration_seconds,omitempty"`
	Cost                    *float64 `json:"cost,omitempty"`
	CostCurrency            *string  `json:"cost_currency,omitempty"`
	TrainTokensPerStep      *int     `json:"train_tokens_per_step,omitempty"`
	TrainTokens             *int     `json:"train_tokens,omitempty"`
	DataTokens              *int     `json:"data_tokens,omitempty"`
	EstimatedStartTime      *int     `json:"estimated_start_time,omitempty"`
}

// JobsOut is a paginated list of fine-tuning jobs.
type JobsOut struct {
	Data   []JobOut `json:"data"`
	Object string   `json:"object"`
	Total  int      `json:"total"`
}

// CompletionTrainingParameters configures completion fine-tuning.
type CompletionTrainingParameters struct {
	TrainingSteps  *int     `json:"training_steps,omitempty"`
	LearningRate   float64  `json:"learning_rate,omitempty"`
	WeightDecay    *float64 `json:"weight_decay,omitempty"`
	WarmupFraction *float64 `json:"warmup_fraction,omitempty"`
	Epochs         *float64 `json:"epochs,omitempty"`
	SeqLen         *int     `json:"seq_len,omitempty"`
	FIMRatio       *float64 `json:"fim_ratio,omitempty"`
}

// UpdateFTModelIn updates a fine-tuned model's metadata.
type UpdateFTModelIn struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// FTModelOut represents a fine-tuned model response.
type FTModelOut struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	ModelType string `json:"model_type,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
}

// ArchiveFTModelOut is the response from archiving a model.
type ArchiveFTModelOut struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Archived bool   `json:"archived"`
}

// UnarchiveFTModelOut is the response from unarchiving a model.
type UnarchiveFTModelOut struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Archived bool   `json:"archived"`
}

// ListParams holds query parameters for listing fine-tuning jobs.
type ListParams struct {
	Page          *int
	PageSize      *int
	Model         *string
	CreatedAfter  *string
	CreatedBefore *string
	CreatedByMe   *bool
	Status        *string
	WandbProject  *string
	WandbName     *string
	Suffix        *string
}
