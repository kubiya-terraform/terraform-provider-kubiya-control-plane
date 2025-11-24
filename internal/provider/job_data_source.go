package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*jobDataSource)(nil)

func NewJobDataSource() datasource.DataSource {
	return &jobDataSource{}
}

type jobDataSource struct {
	client *clients.Client
}

type jobDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Status          types.String `tfsdk:"status"`
	TriggerType     types.String `tfsdk:"trigger_type"`
	CronSchedule    types.String `tfsdk:"cron_schedule"`
	CronTimezone    types.String `tfsdk:"cron_timezone"`
	WebhookURL      types.String `tfsdk:"webhook_url"`
	PlanningMode    types.String `tfsdk:"planning_mode"`
	EntityType      types.String `tfsdk:"entity_type"`
	EntityID        types.String `tfsdk:"entity_id"`
	PromptTemplate  types.String `tfsdk:"prompt_template"`
	SystemPrompt    types.String `tfsdk:"system_prompt"`
	ExecutorType    types.String `tfsdk:"executor_type"`
	WorkerQueueName types.String `tfsdk:"worker_queue_name"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func (d *jobDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (d *jobDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Job from the Control Plane by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Job ID",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Job name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Job description",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the job is enabled",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Job status",
				Computed:    true,
			},
			"trigger_type": schema.StringAttribute{
				Description: "Trigger type: 'cron', 'webhook', or 'manual'",
				Computed:    true,
			},
			"cron_schedule": schema.StringAttribute{
				Description: "Cron expression",
				Computed:    true,
			},
			"cron_timezone": schema.StringAttribute{
				Description: "Timezone for cron schedule",
				Computed:    true,
			},
			"webhook_url": schema.StringAttribute{
				Description: "Full webhook URL",
				Computed:    true,
			},
			"planning_mode": schema.StringAttribute{
				Description: "Planning mode",
				Computed:    true,
			},
			"entity_type": schema.StringAttribute{
				Description: "Entity type: 'agent', 'team', or 'workflow'",
				Computed:    true,
			},
			"entity_id": schema.StringAttribute{
				Description: "Entity ID",
				Computed:    true,
			},
			"prompt_template": schema.StringAttribute{
				Description: "Prompt template",
				Computed:    true,
			},
			"system_prompt": schema.StringAttribute{
				Description: "System prompt",
				Computed:    true,
			},
			"executor_type": schema.StringAttribute{
				Description: "Executor routing type",
				Computed:    true,
			},
			"worker_queue_name": schema.StringAttribute{
				Description: "Worker queue name",
				Computed:    true,
			},
			"environment_name": schema.StringAttribute{
				Description: "Environment name",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the job was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the job was last updated",
				Computed:    true,
			},
		},
	}
}

func (d *jobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *jobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config jobDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	job, err := d.client.GetJob(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading job", err.Error())
		return
	}

	config.ID = types.StringValue(job.ID)
	config.Name = types.StringValue(job.Name)
	config.Enabled = types.BoolValue(job.Enabled)
	config.TriggerType = types.StringValue(job.TriggerType)
	config.PlanningMode = types.StringValue(job.PlanningMode)
	config.PromptTemplate = types.StringValue(job.PromptTemplate)
	config.ExecutorType = types.StringValue(job.ExecutorType)

	if job.Description != nil {
		config.Description = types.StringValue(*job.Description)
	}

	if job.Status != "" {
		config.Status = types.StringValue(job.Status)
	}

	if job.CronSchedule != nil {
		config.CronSchedule = types.StringValue(*job.CronSchedule)
	}

	if job.CronTimezone != nil {
		config.CronTimezone = types.StringValue(*job.CronTimezone)
	}

	if job.WebhookURL != nil {
		config.WebhookURL = types.StringValue(*job.WebhookURL)
	}

	if job.EntityType != nil {
		config.EntityType = types.StringValue(*job.EntityType)
	}

	if job.EntityID != nil {
		config.EntityID = types.StringValue(*job.EntityID)
	}

	if job.SystemPrompt != nil {
		config.SystemPrompt = types.StringValue(*job.SystemPrompt)
	}

	if job.WorkerQueueName != nil {
		config.WorkerQueueName = types.StringValue(*job.WorkerQueueName)
	}

	if job.EnvironmentName != nil {
		config.EnvironmentName = types.StringValue(*job.EnvironmentName)
	}

	if job.CreatedAt != nil {
		config.CreatedAt = types.StringValue(job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	if job.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(job.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
