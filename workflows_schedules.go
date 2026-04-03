package mistral

import (
	"context"
	"fmt"

	"github.com/VikingOwl91/mistral-go-sdk/workflow"
)

// ListWorkflowSchedules lists workflow schedules.
func (c *Client) ListWorkflowSchedules(ctx context.Context) (*workflow.ScheduleListResponse, error) {
	var resp workflow.ScheduleListResponse
	if err := c.doJSON(ctx, "GET", "/v1/workflows/schedules", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ScheduleWorkflow creates a workflow schedule.
func (c *Client) ScheduleWorkflow(ctx context.Context, req *workflow.ScheduleRequest) (*workflow.ScheduleResponse, error) {
	var resp workflow.ScheduleResponse
	if err := c.doJSON(ctx, "POST", "/v1/workflows/schedules", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UnscheduleWorkflow removes a workflow schedule.
func (c *Client) UnscheduleWorkflow(ctx context.Context, scheduleID string) error {
	resp, err := c.do(ctx, "DELETE", fmt.Sprintf("/v1/workflows/schedules/%s", scheduleID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return parseAPIError(resp)
	}
	return nil
}
