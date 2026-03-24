package observability

// Campaign represents an observability campaign.
type Campaign struct {
	ID            string        `json:"id"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	DeletedAt     *string       `json:"deleted_at,omitempty"`
	Name          string        `json:"name"`
	OwnerID       string        `json:"owner_id"`
	WorkspaceID   string        `json:"workspace_id"`
	Description   string        `json:"description"`
	MaxNbEvents   int           `json:"max_nb_events"`
	SearchParams  FilterPayload `json:"search_params"`
	Judge         Judge         `json:"judge"`
}

// CreateCampaignRequest creates a new campaign.
type CreateCampaignRequest struct {
	SearchParams FilterPayload `json:"search_params"`
	JudgeID      string        `json:"judge_id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	MaxNbEvents  int           `json:"max_nb_events"`
}

// CampaignStatusResponse is the response for campaign status.
type CampaignStatusResponse struct {
	Status TaskStatus `json:"status"`
}

// ListCampaignsResponse is the response from listing campaigns.
type ListCampaignsResponse struct {
	Count    int        `json:"count"`
	Results  []Campaign `json:"results,omitempty"`
	Next     *string    `json:"next,omitempty"`
	Previous *string    `json:"previous,omitempty"`
}

// ListCampaignEventsResponse is the response from listing campaign events.
type ListCampaignEventsResponse struct {
	Count    int                          `json:"count"`
	Results  []ChatCompletionEventPreview `json:"results,omitempty"`
	Next     *string                      `json:"next,omitempty"`
	Previous *string                      `json:"previous,omitempty"`
}
