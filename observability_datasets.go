package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

// CreateDataset creates a new observability dataset.
func (c *Client) CreateDataset(ctx context.Context, req *observability.CreateDatasetRequest) (*observability.Dataset, error) {
	var resp observability.Dataset
	if err := c.doJSON(ctx, "POST", "/v1/observability/datasets", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListDatasets lists observability datasets.
func (c *Client) ListDatasets(ctx context.Context, params *observability.SearchParams) (*observability.ListDatasetsResponse, error) {
	path := "/v1/observability/datasets"
	if params != nil {
		q := url.Values{}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if params.Q != nil {
			q.Set("q", *params.Q)
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp observability.ListDatasetsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDataset retrieves a dataset by ID.
func (c *Client) GetDataset(ctx context.Context, datasetID string) (*observability.Dataset, error) {
	var resp observability.Dataset
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/datasets/%s", datasetID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDataset updates a dataset.
func (c *Client) UpdateDataset(ctx context.Context, datasetID string, req *observability.UpdateDatasetRequest) (*observability.Dataset, error) {
	var resp observability.Dataset
	if err := c.doJSON(ctx, "PATCH", fmt.Sprintf("/v1/observability/datasets/%s", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDataset deletes a dataset.
func (c *Client) DeleteDataset(ctx context.Context, datasetID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/observability/datasets/%s", datasetID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// ExportDatasetToJSONL exports a dataset to JSONL format.
func (c *Client) ExportDatasetToJSONL(ctx context.Context, datasetID string) (*observability.ExportDatasetResponse, error) {
	var resp observability.ExportDatasetResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/datasets/%s/exports/to-jsonl", datasetID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Dataset records

// ListDatasetRecords lists records in a dataset.
func (c *Client) ListDatasetRecords(ctx context.Context, datasetID string, params *observability.PaginationParams) (*observability.ListRecordsResponse, error) {
	path := fmt.Sprintf("/v1/observability/datasets/%s/records", datasetID)
	if params != nil {
		q := url.Values{}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp observability.ListRecordsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateDatasetRecord creates a record in a dataset.
func (c *Client) CreateDatasetRecord(ctx context.Context, datasetID string, req *observability.CreateRecordRequest) (*observability.DatasetRecord, error) {
	var resp observability.DatasetRecord
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/records", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDatasetRecord retrieves a dataset record by ID.
func (c *Client) GetDatasetRecord(ctx context.Context, recordID string) (*observability.DatasetRecord, error) {
	var resp observability.DatasetRecord
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/dataset-records/%s", recordID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDatasetRecordPayload updates a record's payload.
func (c *Client) UpdateDatasetRecordPayload(ctx context.Context, recordID string, req *observability.UpdateRecordPayloadRequest) (*observability.DatasetRecord, error) {
	var resp observability.DatasetRecord
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/observability/dataset-records/%s/payload", recordID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateDatasetRecordProperties updates a record's properties.
func (c *Client) UpdateDatasetRecordProperties(ctx context.Context, recordID string, req *observability.UpdateRecordPropertiesRequest) (*observability.DatasetRecord, error) {
	var resp observability.DatasetRecord
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/observability/dataset-records/%s/properties", recordID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteDatasetRecord deletes a dataset record.
func (c *Client) DeleteDatasetRecord(ctx context.Context, recordID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/observability/dataset-records/%s", recordID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// BulkDeleteDatasetRecords deletes multiple dataset records.
func (c *Client) BulkDeleteDatasetRecords(ctx context.Context, req *observability.BulkDeleteRecordsRequest) error {
	return c.doJSON(ctx, "POST", "/v1/observability/dataset-records/bulk-delete", req, nil)
}

// JudgeDatasetRecord judges a dataset record.
func (c *Client) JudgeDatasetRecord(ctx context.Context, recordID string, req *observability.JudgeRecordRequest) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/dataset-records/%s/live-judging", recordID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Import operations

// ImportDatasetFromCampaign imports records from a campaign.
func (c *Client) ImportDatasetFromCampaign(ctx context.Context, datasetID string, req *observability.ImportFromCampaignRequest) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/imports/from-campaign", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportDatasetFromExplorer imports records from explorer events.
func (c *Client) ImportDatasetFromExplorer(ctx context.Context, datasetID string, req *observability.ImportFromExplorerRequest) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/imports/from-explorer", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportDatasetFromFile imports records from a file.
func (c *Client) ImportDatasetFromFile(ctx context.Context, datasetID string, req *observability.ImportFromFileRequest) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/imports/from-file", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportDatasetFromPlayground imports records from playground conversations.
func (c *Client) ImportDatasetFromPlayground(ctx context.Context, datasetID string, req *observability.ImportFromPlaygroundRequest) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/imports/from-playground", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ImportDatasetFromDataset imports records from another dataset.
func (c *Client) ImportDatasetFromDataset(ctx context.Context, datasetID string, req *observability.ImportFromDatasetRequest) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/datasets/%s/imports/from-dataset", datasetID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Tasks

// ListDatasetTasks lists import tasks for a dataset.
func (c *Client) ListDatasetTasks(ctx context.Context, datasetID string, params *observability.PaginationParams) (*observability.ListTasksResponse, error) {
	path := fmt.Sprintf("/v1/observability/datasets/%s/tasks", datasetID)
	if params != nil {
		q := url.Values{}
		if params.PageSize != nil {
			q.Set("page_size", strconv.Itoa(*params.PageSize))
		}
		if params.Page != nil {
			q.Set("page", strconv.Itoa(*params.Page))
		}
		if encoded := q.Encode(); encoded != "" {
			path += "?" + encoded
		}
	}
	var resp observability.ListTasksResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetDatasetTask retrieves an import task by ID.
func (c *Client) GetDatasetTask(ctx context.Context, datasetID, taskID string) (*observability.DatasetImportTask, error) {
	var resp observability.DatasetImportTask
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/datasets/%s/tasks/%s", datasetID, taskID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
