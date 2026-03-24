package observability

// ChatCompletionEvent is a full chat completion event.
type ChatCompletionEvent struct {
	EventID                  string                    `json:"event_id"`
	CorrelationID            string                    `json:"correlation_id"`
	CreatedAt                string                    `json:"created_at"`
	ExtraFields              map[string]any            `json:"extra_fields,omitempty"`
	NbInputTokens            int                       `json:"nb_input_tokens"`
	NbOutputTokens           int                       `json:"nb_output_tokens"`
	EnabledTools             []map[string]any          `json:"enabled_tools,omitempty"`
	RequestMessages          []map[string]any          `json:"request_messages,omitempty"`
	ResponseMessages         []map[string]any          `json:"response_messages,omitempty"`
	NbMessages               int                       `json:"nb_messages"`
	ChatTranscriptionEvents  []ChatTranscriptionEvent  `json:"chat_transcription_events,omitempty"`
}

// ChatCompletionEventPreview is a summary of a chat completion event.
type ChatCompletionEventPreview struct {
	EventID        string         `json:"event_id"`
	CorrelationID  string         `json:"correlation_id"`
	CreatedAt      string         `json:"created_at"`
	ExtraFields    map[string]any `json:"extra_fields,omitempty"`
	NbInputTokens  int            `json:"nb_input_tokens"`
	NbOutputTokens int            `json:"nb_output_tokens"`
}

// ChatTranscriptionEvent is an audio transcription within a chat event.
type ChatTranscriptionEvent struct {
	AudioURL        string         `json:"audio_url"`
	Model           string         `json:"model"`
	ResponseMessage map[string]any `json:"response_message"`
}

// SearchEventsRequest is the request body for searching chat completion events.
type SearchEventsRequest struct {
	SearchParams FilterPayload `json:"search_params"`
	ExtraFields  []string      `json:"extra_fields,omitempty"`
}

// SearchEventsResponse is the response from searching events.
type SearchEventsResponse struct {
	Results []ChatCompletionEventPreview `json:"results,omitempty"`
	Next    *string                      `json:"next,omitempty"`
	Cursor  *string                      `json:"cursor,omitempty"`
}

// SearchEventIDsRequest is the request body for searching event IDs.
type SearchEventIDsRequest struct {
	SearchParams FilterPayload `json:"search_params"`
	ExtraFields  []string      `json:"extra_fields,omitempty"`
}

// SearchEventIDsResponse is the response from searching event IDs.
type SearchEventIDsResponse struct {
	CompletionEventIDs []string `json:"completion_event_ids"`
}

// JudgeEventRequest is the request body for judging a chat completion event.
type JudgeEventRequest struct {
	JudgeDefinition CreateJudgeRequest `json:"judge_definition"`
}

// SimilarEventsResponse is the response from fetching similar events.
type SimilarEventsResponse struct {
	Count    int                          `json:"count"`
	Results  []ChatCompletionEventPreview `json:"results,omitempty"`
	Next     *string                      `json:"next,omitempty"`
	Previous *string                      `json:"previous,omitempty"`
}
