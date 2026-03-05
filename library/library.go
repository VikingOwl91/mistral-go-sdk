package library

// CreateRequest creates a new library.
type CreateRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	ChunkSize   *int    `json:"chunk_size,omitempty"`
}

// UpdateRequest updates a library.
type UpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Library represents a document library.
type Library struct {
	ID                            string  `json:"id"`
	Name                          string  `json:"name"`
	CreatedAt                     string  `json:"created_at"`
	UpdatedAt                     string  `json:"updated_at"`
	OwnerID                       *string `json:"owner_id"`
	OwnerType                     string  `json:"owner_type"`
	TotalSize                     int     `json:"total_size"`
	NbDocuments                   int     `json:"nb_documents"`
	ChunkSize                     *int    `json:"chunk_size"`
	Emoji                         *string `json:"emoji,omitempty"`
	Description                   *string `json:"description,omitempty"`
	GeneratedDescription          *string `json:"generated_description,omitempty"`
	ExplicitUserMembersCount      *int    `json:"explicit_user_members_count,omitempty"`
	ExplicitWorkspaceMembersCount *int    `json:"explicit_workspace_members_count,omitempty"`
	OrgSharingRole                *string `json:"org_sharing_role,omitempty"`
	GeneratedName                 *string `json:"generated_name,omitempty"`
}

// ListLibraryOut is the response from listing libraries.
type ListLibraryOut struct {
	Data []Library `json:"data"`
}

// Document represents a document in a library.
type Document struct {
	ID                           string         `json:"id"`
	LibraryID                    string         `json:"library_id"`
	Hash                         *string        `json:"hash"`
	MimeType                     *string        `json:"mime_type"`
	Extension                    *string        `json:"extension"`
	Size                         *int           `json:"size"`
	Name                         string         `json:"name"`
	Summary                      *string        `json:"summary,omitempty"`
	CreatedAt                    string         `json:"created_at"`
	LastProcessedAt              *string        `json:"last_processed_at,omitempty"`
	NumberOfPages                *int           `json:"number_of_pages,omitempty"`
	ProcessStatus                string         `json:"process_status"`
	UploadedByID                 *string        `json:"uploaded_by_id"`
	UploadedByType               string         `json:"uploaded_by_type"`
	TokensProcessingMainContent  *int           `json:"tokens_processing_main_content,omitempty"`
	TokensProcessingSummary      *int           `json:"tokens_processing_summary,omitempty"`
	URL                          *string        `json:"url,omitempty"`
	Attributes                   map[string]any `json:"attributes,omitempty"`
	ProcessingStatus             string         `json:"processing_status"`
	TokensProcessingTotal        int            `json:"tokens_processing_total"`
}

// DocumentUpdateRequest updates a document's metadata.
type DocumentUpdateRequest struct {
	Name       *string        `json:"name,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

// DocumentTextContent holds the extracted text of a document.
type DocumentTextContent struct {
	Text string `json:"text"`
}

// ListDocumentOut is a paginated list of documents.
type ListDocumentOut struct {
	Pagination PaginationInfo `json:"pagination"`
	Data       []Document     `json:"data"`
}

// PaginationInfo holds pagination metadata.
type PaginationInfo struct {
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	HasMore     bool `json:"has_more"`
}

// ProcessingStatusOut holds document processing status.
type ProcessingStatusOut struct {
	DocumentID       string `json:"document_id"`
	ProcessStatus    string `json:"process_status"`
	ProcessingStatus string `json:"processing_status"`
}

// ListDocumentParams holds query parameters for listing documents.
type ListDocumentParams struct {
	Search           *string
	PageSize         *int
	Page             *int
	FiltersAttributes *string
	SortBy           *string
	SortOrder        *string
}

// SharingRequest creates or updates library sharing.
type SharingRequest struct {
	OrgID         *string `json:"org_id,omitempty"`
	Level         string  `json:"level"`
	ShareWithUUID string  `json:"share_with_uuid"`
	ShareWithType string  `json:"share_with_type"`
}

// SharingDeleteRequest removes library sharing.
type SharingDeleteRequest struct {
	OrgID         *string `json:"org_id,omitempty"`
	ShareWithUUID string  `json:"share_with_uuid"`
	ShareWithType string  `json:"share_with_type"`
}

// SharingOut represents a sharing entry.
type SharingOut struct {
	LibraryID     string  `json:"library_id"`
	UserID        *string `json:"user_id,omitempty"`
	OrgID         string  `json:"org_id"`
	Role          string  `json:"role"`
	ShareWithType string  `json:"share_with_type"`
	ShareWithUUID *string `json:"share_with_uuid"`
}

// ListSharingOut is the response from listing sharing entries.
type ListSharingOut struct {
	Data []SharingOut `json:"data"`
}
