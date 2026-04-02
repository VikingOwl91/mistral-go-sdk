package workflow

// ScheduleRequest is the request body for scheduling a workflow.
type ScheduleRequest struct {
	Schedule               ScheduleDefinition `json:"schedule"`
	WorkflowRegistrationID *string            `json:"workflow_registration_id,omitempty"`
	WorkflowIdentifier     *string            `json:"workflow_identifier,omitempty"`
	ScheduleID             *string            `json:"schedule_id,omitempty"`
	DeploymentName         *string            `json:"deployment_name,omitempty"`
}

// ScheduleDefinition describes when and how a workflow should be scheduled.
type ScheduleDefinition struct {
	Input           any                `json:"input"`
	Calendars       []ScheduleCalendar `json:"calendars,omitempty"`
	Intervals       []ScheduleInterval `json:"intervals,omitempty"`
	CronExpressions []string           `json:"cron_expressions,omitempty"`
	Skip            []ScheduleCalendar `json:"skip,omitempty"`
	StartAt         *string            `json:"start_at,omitempty"`
	EndAt           *string            `json:"end_at,omitempty"`
	Jitter          *string            `json:"jitter,omitempty"`
	TimeZoneName    *string            `json:"time_zone_name,omitempty"`
	Policy          *SchedulePolicy    `json:"policy,omitempty"`
}

// ScheduleCalendar defines calendar-based schedule entries.
type ScheduleCalendar struct {
	Second     []ScheduleRange `json:"second,omitempty"`
	Minute     []ScheduleRange `json:"minute,omitempty"`
	Hour       []ScheduleRange `json:"hour,omitempty"`
	DayOfMonth []ScheduleRange `json:"day_of_month,omitempty"`
	Month      []ScheduleRange `json:"month,omitempty"`
	Year       []ScheduleRange `json:"year,omitempty"`
	DayOfWeek  []ScheduleRange `json:"day_of_week,omitempty"`
	Comment    *string         `json:"comment,omitempty"`
}

// ScheduleRange defines a numeric range for calendar schedules.
type ScheduleRange struct {
	Start int `json:"start"`
	End   int `json:"end,omitempty"`
	Step  int `json:"step,omitempty"`
}

// ScheduleInterval defines an interval-based schedule.
type ScheduleInterval struct {
	Every  string  `json:"every"`
	Offset *string `json:"offset,omitempty"`
}

// SchedulePolicy controls schedule overlap and failure behavior.
type SchedulePolicy struct {
	CatchupWindowSeconds int  `json:"catchup_window_seconds,omitempty"`
	Overlap              *int `json:"overlap,omitempty"`
	PauseOnFailure       bool `json:"pause_on_failure,omitempty"`
}

// ScheduleResponse is the response from creating a workflow schedule.
type ScheduleResponse struct {
	ScheduleID string `json:"schedule_id"`
}

// ScheduleListResponse is the response from listing workflow schedules.
type ScheduleListResponse struct {
	Schedules []Schedule `json:"schedules"`
}

// Schedule represents a workflow schedule.
type Schedule struct {
	ScheduleID   string             `json:"schedule_id"`
	Definition   ScheduleDefinition `json:"definition"`
	WorkflowName string             `json:"workflow_name"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}
