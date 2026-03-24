package audio

// VoiceResponse represents a voice entity.
type VoiceResponse struct {
	Name            string   `json:"name"`
	ID              string   `json:"id"`
	CreatedAt       string   `json:"created_at"`
	UserID          *string  `json:"user_id,omitempty"`
	Slug            *string  `json:"slug,omitempty"`
	Languages       []string `json:"languages,omitempty"`
	Gender          *string  `json:"gender,omitempty"`
	Age             *int     `json:"age,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Color           *string  `json:"color,omitempty"`
	RetentionNotice *int     `json:"retention_notice,omitempty"`
}

// VoiceCreateRequest creates a custom voice.
type VoiceCreateRequest struct {
	Name            string   `json:"name"`
	SampleAudio     string   `json:"sample_audio"`
	Slug            *string  `json:"slug,omitempty"`
	Languages       []string `json:"languages,omitempty"`
	Gender          *string  `json:"gender,omitempty"`
	Age             *int     `json:"age,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Color           *string  `json:"color,omitempty"`
	RetentionNotice *int     `json:"retention_notice,omitempty"`
	SampleFilename  *string  `json:"sample_filename,omitempty"`
}

// VoiceUpdateRequest updates a voice.
type VoiceUpdateRequest struct {
	Name      *string  `json:"name,omitempty"`
	Languages []string `json:"languages,omitempty"`
	Gender    *string  `json:"gender,omitempty"`
	Age       *int     `json:"age,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// VoiceListResponse is the response from listing voices.
type VoiceListResponse struct {
	Items      []VoiceResponse `json:"items"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}
