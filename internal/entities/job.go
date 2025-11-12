package entities

import (
	"time"
)

// Job represents a job in the control plane
type Job struct {
	ID                 string                 `json:"id,omitempty"`
	OrganizationID     string                 `json:"organization_id,omitempty"`
	Name               string                 `json:"name"`
	Description        *string                `json:"description,omitempty"`
	Enabled            bool                   `json:"enabled"`
	Status             string                 `json:"status,omitempty"`
	TriggerType        string                 `json:"trigger_type"`
	CronSchedule       *string                `json:"cron_schedule,omitempty"`
	CronTimezone       *string                `json:"cron_timezone,omitempty"`
	WebhookURL         *string                `json:"webhook_url,omitempty"`
	WebhookSecret      *string                `json:"webhook_secret,omitempty"`
	TemporalScheduleID *string                `json:"temporal_schedule_id,omitempty"`
	PlanningMode       string                 `json:"planning_mode"`
	EntityType         *string                `json:"entity_type,omitempty"`
	EntityID           *string                `json:"entity_id,omitempty"`
	PromptTemplate     string                 `json:"prompt_template"`
	SystemPrompt       *string                `json:"system_prompt,omitempty"`
	ExecutorType       string                 `json:"executor_type"`
	WorkerQueueName    *string                `json:"worker_queue_name,omitempty"`
	EnvironmentName    *string                `json:"environment_name,omitempty"`
	Config             map[string]interface{} `json:"config,omitempty"`
	ExecutionEnv       *ExecutionEnvironment  `json:"execution_environment,omitempty"`
	CreatedAt          *time.Time             `json:"created_at,omitempty"`
	UpdatedAt          *time.Time             `json:"updated_at,omitempty"`
}

// JobCreateRequest represents the request to create a job
type JobCreateRequest struct {
	Name            string                 `json:"name"`
	Description     *string                `json:"description,omitempty"`
	Enabled         bool                   `json:"enabled"`
	TriggerType     string                 `json:"trigger_type"`
	CronSchedule    *string                `json:"cron_schedule,omitempty"`
	CronTimezone    *string                `json:"cron_timezone,omitempty"`
	PlanningMode    string                 `json:"planning_mode"`
	EntityType      *string                `json:"entity_type,omitempty"`
	EntityID        *string                `json:"entity_id,omitempty"`
	PromptTemplate  string                 `json:"prompt_template"`
	SystemPrompt    *string                `json:"system_prompt,omitempty"`
	ExecutorType    string                 `json:"executor_type"`
	WorkerQueueName *string                `json:"worker_queue_name,omitempty"`
	EnvironmentName *string                `json:"environment_name,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
	ExecutionEnv    *ExecutionEnvironment  `json:"execution_environment,omitempty"`
}

// JobUpdateRequest represents the request to update a job
type JobUpdateRequest struct {
	Name            *string                `json:"name,omitempty"`
	Description     *string                `json:"description,omitempty"`
	Enabled         *bool                  `json:"enabled,omitempty"`
	TriggerType     *string                `json:"trigger_type,omitempty"`
	CronSchedule    *string                `json:"cron_schedule,omitempty"`
	CronTimezone    *string                `json:"cron_timezone,omitempty"`
	PlanningMode    *string                `json:"planning_mode,omitempty"`
	EntityType      *string                `json:"entity_type,omitempty"`
	EntityID        *string                `json:"entity_id,omitempty"`
	PromptTemplate  *string                `json:"prompt_template,omitempty"`
	SystemPrompt    *string                `json:"system_prompt,omitempty"`
	ExecutorType    *string                `json:"executor_type,omitempty"`
	WorkerQueueName *string                `json:"worker_queue_name,omitempty"`
	EnvironmentName *string                `json:"environment_name,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
	ExecutionEnv    *ExecutionEnvironment  `json:"execution_environment,omitempty"`
}
