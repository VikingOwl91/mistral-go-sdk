package observability

import "encoding/json"

// ConversationSource indicates how a dataset record was created.
type ConversationSource string

const (
	SourceExplorer     ConversationSource = "EXPLORER"
	SourceUploadedFile ConversationSource = "UPLOADED_FILE"
	SourceDirectInput  ConversationSource = "DIRECT_INPUT"
	SourcePlayground   ConversationSource = "PLAYGROUND"
)

// Dataset represents a dataset entity.
type Dataset struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	OwnerID     string  `json:"owner_id"`
	WorkspaceID string  `json:"workspace_id"`
}

// CreateDatasetRequest creates a new dataset.
type CreateDatasetRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateDatasetRequest updates a dataset.
type UpdateDatasetRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// DatasetRecord is a single record in a dataset.
type DatasetRecord struct {
	ID         string             `json:"id"`
	CreatedAt  string             `json:"created_at"`
	UpdatedAt  string             `json:"updated_at"`
	DeletedAt  *string            `json:"deleted_at,omitempty"`
	DatasetID  string             `json:"dataset_id"`
	Payload    ConversationPayload `json:"payload"`
	Properties map[string]any     `json:"properties,omitempty"`
	Source     ConversationSource `json:"source"`
}

// ConversationPayload holds the messages for a dataset record.
type ConversationPayload struct {
	Messages []map[string]any `json:"messages"`
}

// CreateRecordRequest creates a new dataset record.
type CreateRecordRequest struct {
	Payload    ConversationPayload `json:"payload"`
	Properties map[string]any     `json:"properties"`
}

// UpdateRecordPayloadRequest updates a record's payload.
type UpdateRecordPayloadRequest struct {
	Payload ConversationPayload `json:"payload"`
}

// UpdateRecordPropertiesRequest updates a record's properties.
type UpdateRecordPropertiesRequest struct {
	Properties map[string]any `json:"properties"`
}

// BulkDeleteRecordsRequest deletes multiple records.
type BulkDeleteRecordsRequest struct {
	DatasetRecordIDs []string `json:"dataset_record_ids"`
}

// JudgeRecordRequest judges a dataset record.
type JudgeRecordRequest struct {
	JudgeDefinition CreateJudgeRequest `json:"judge_definition"`
}

// DatasetImportTask tracks an async import operation.
type DatasetImportTask struct {
	ID          string     `json:"id"`
	CreatedAt   string     `json:"created_at"`
	UpdatedAt   string     `json:"updated_at"`
	DeletedAt   *string    `json:"deleted_at,omitempty"`
	CreatorID   string     `json:"creator_id"`
	DatasetID   string     `json:"dataset_id"`
	WorkspaceID string     `json:"workspace_id"`
	Status      TaskStatus `json:"status"`
	Progress    *int       `json:"progress,omitempty"`
	Message     *string    `json:"message,omitempty"`
}

// ExportDatasetResponse is the response from exporting a dataset.
type ExportDatasetResponse struct {
	FileURL string `json:"file_url"`
}

// Import request types.

// ImportFromCampaignRequest imports records from a campaign.
type ImportFromCampaignRequest struct {
	CampaignID string `json:"campaign_id"`
}

// ImportFromExplorerRequest imports records from explorer events.
type ImportFromExplorerRequest struct {
	CompletionEventIDs []string `json:"completion_event_ids"`
}

// ImportFromFileRequest imports records from a file.
type ImportFromFileRequest struct {
	FileID string `json:"file_id"`
}

// ImportFromPlaygroundRequest imports records from playground conversations.
type ImportFromPlaygroundRequest struct {
	ConversationIDs []string `json:"conversation_ids"`
}

// ImportFromDatasetRequest imports records from another dataset.
type ImportFromDatasetRequest struct {
	DatasetRecordIDs []string `json:"dataset_record_ids"`
}

// List response types.

// ListDatasetsResponse is the response from listing datasets.
type ListDatasetsResponse struct {
	Count    int       `json:"count"`
	Results  []Dataset `json:"results,omitempty"`
	Next     *string   `json:"next,omitempty"`
	Previous *string   `json:"previous,omitempty"`
}

// ListRecordsResponse is the response from listing dataset records.
type ListRecordsResponse struct {
	Count    int             `json:"count"`
	Results  []DatasetRecord `json:"results,omitempty"`
	Next     *string         `json:"next,omitempty"`
	Previous *string         `json:"previous,omitempty"`
}

// ListTasksResponse is the response from listing import tasks.
type ListTasksResponse struct {
	Count    int                 `json:"count"`
	Results  []DatasetImportTask `json:"results,omitempty"`
	Next     *string             `json:"next,omitempty"`
	Previous *string             `json:"previous,omitempty"`
}

// JudgeResultResponse is the raw response from judging operations.
// The shape depends on the judge type (classification or regression).
type JudgeResultResponse json.RawMessage
