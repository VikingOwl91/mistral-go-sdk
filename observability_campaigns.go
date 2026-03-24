package mistral

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"somegit.dev/vikingowl/mistral-go-sdk/observability"
)

// CreateCampaign creates a new observability campaign.
func (c *Client) CreateCampaign(ctx context.Context, req *observability.CreateCampaignRequest) (*observability.Campaign, error) {
	var resp observability.Campaign
	if err := c.doJSON(ctx, "POST", "/v1/observability/campaigns", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCampaigns lists observability campaigns.
func (c *Client) ListCampaigns(ctx context.Context, params *observability.SearchParams) (*observability.ListCampaignsResponse, error) {
	path := "/v1/observability/campaigns"
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
	var resp observability.ListCampaignsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetCampaign retrieves a campaign by ID.
func (c *Client) GetCampaign(ctx context.Context, campaignID string) (*observability.Campaign, error) {
	var resp observability.Campaign
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/campaigns/%s", campaignID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCampaign deletes a campaign.
func (c *Client) DeleteCampaign(ctx context.Context, campaignID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/observability/campaigns/%s", campaignID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// GetCampaignStatus retrieves the status of a campaign.
func (c *Client) GetCampaignStatus(ctx context.Context, campaignID string) (*observability.CampaignStatusResponse, error) {
	var resp observability.CampaignStatusResponse
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/campaigns/%s/status", campaignID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListCampaignEvents lists events selected by a campaign.
func (c *Client) ListCampaignEvents(ctx context.Context, campaignID string, params *observability.PaginationParams) (*observability.ListCampaignEventsResponse, error) {
	path := fmt.Sprintf("/v1/observability/campaigns/%s/selected-events", campaignID)
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
	var resp observability.ListCampaignEventsResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
