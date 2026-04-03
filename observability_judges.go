package mistral

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/VikingOwl91/mistral-go-sdk/observability"
)

// CreateJudge creates a new observability judge.
func (c *Client) CreateJudge(ctx context.Context, req *observability.CreateJudgeRequest) (*observability.Judge, error) {
	var resp observability.Judge
	if err := c.doJSON(ctx, "POST", "/v1/observability/judges", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListJudges lists observability judges.
func (c *Client) ListJudges(ctx context.Context, params *observability.SearchParams) (*observability.ListJudgesResponse, error) {
	path := "/v1/observability/judges"
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
	var resp observability.ListJudgesResponse
	if err := c.doJSON(ctx, "GET", path, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetJudge retrieves a judge by ID.
func (c *Client) GetJudge(ctx context.Context, judgeID string) (*observability.Judge, error) {
	var resp observability.Judge
	if err := c.doJSON(ctx, "GET", fmt.Sprintf("/v1/observability/judges/%s", judgeID), nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateJudge updates a judge.
func (c *Client) UpdateJudge(ctx context.Context, judgeID string, req *observability.UpdateJudgeRequest) (*observability.Judge, error) {
	var resp observability.Judge
	if err := c.doJSON(ctx, "PUT", fmt.Sprintf("/v1/observability/judges/%s", judgeID), req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteJudge deletes a judge.
func (c *Client) DeleteJudge(ctx context.Context, judgeID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/observability/judges/%s", judgeID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}

// JudgeConversation performs live judging on a conversation.
func (c *Client) JudgeConversation(ctx context.Context, judgeID string, req *observability.JudgeConversationRequest) (json.RawMessage, error) {
	var resp json.RawMessage
	if err := c.doJSON(ctx, "POST", fmt.Sprintf("/v1/observability/judges/%s/live-judging", judgeID), req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
