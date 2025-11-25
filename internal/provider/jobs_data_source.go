package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*jobsDataSource)(nil)

func NewJobsDataSource() datasource.DataSource {
	return &jobsDataSource{}
}

type jobsDataSource struct {
	client *clients.Client
}

type jobsDataSourceModel struct {
	Jobs []jobDataSourceModel `tfsdk:"jobs"`
}

func (d *jobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobs"
}

func (d *jobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all Jobs from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"jobs": schema.ListNestedAttribute{
				Description: "List of jobs",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Job ID",
							Computed:    true,
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
							Description: "Cron expression in 5-field format: 'minute hour day month weekday'",
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
				},
			},
		},
	}
}

func (d *jobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *jobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data jobsDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	jobs, err := d.client.ListJobs()
	if err != nil {
		resp.Diagnostics.AddError("Error listing jobs", err.Error())
		return
	}

	data.Jobs = make([]jobDataSourceModel, 0, len(jobs))
	for _, job := range jobs {
		jobModel := jobDataSourceModel{
			ID:             types.StringValue(job.ID),
			Name:           types.StringValue(job.Name),
			Enabled:        types.BoolValue(job.Enabled),
			TriggerType:    types.StringValue(job.TriggerType),
			PlanningMode:   types.StringValue(job.PlanningMode),
			PromptTemplate: types.StringValue(job.PromptTemplate),
			ExecutorType:   types.StringValue(job.ExecutorType),
		}

		if job.Description != nil {
			jobModel.Description = types.StringValue(*job.Description)
		} else {
			jobModel.Description = types.StringNull()
		}

		if job.Status != "" {
			jobModel.Status = types.StringValue(job.Status)
		} else {
			jobModel.Status = types.StringNull()
		}

		if job.CronSchedule != nil {
			jobModel.CronSchedule = types.StringValue(*job.CronSchedule)
		} else {
			jobModel.CronSchedule = types.StringNull()
		}

		if job.CronTimezone != nil {
			jobModel.CronTimezone = types.StringValue(*job.CronTimezone)
		} else {
			jobModel.CronTimezone = types.StringNull()
		}

		if job.WebhookURL != nil {
			jobModel.WebhookURL = types.StringValue(*job.WebhookURL)
		} else {
			jobModel.WebhookURL = types.StringNull()
		}

		if job.EntityType != nil {
			jobModel.EntityType = types.StringValue(*job.EntityType)
		} else {
			jobModel.EntityType = types.StringNull()
		}

		if job.EntityID != nil {
			jobModel.EntityID = types.StringValue(*job.EntityID)
		} else {
			jobModel.EntityID = types.StringNull()
		}

		if job.SystemPrompt != nil {
			jobModel.SystemPrompt = types.StringValue(*job.SystemPrompt)
		} else {
			jobModel.SystemPrompt = types.StringNull()
		}

		if job.WorkerQueueName != nil {
			jobModel.WorkerQueueName = types.StringValue(*job.WorkerQueueName)
		} else {
			jobModel.WorkerQueueName = types.StringNull()
		}

		if job.EnvironmentName != nil {
			jobModel.EnvironmentName = types.StringValue(*job.EnvironmentName)
		} else {
			jobModel.EnvironmentName = types.StringNull()
		}

		if job.CreatedAt != nil {
			jobModel.CreatedAt = types.StringValue(job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
		} else {
			jobModel.CreatedAt = types.StringNull()
		}

		if job.UpdatedAt != nil {
			jobModel.UpdatedAt = types.StringValue(job.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
		} else {
			jobModel.UpdatedAt = types.StringNull()
		}

		data.Jobs = append(data.Jobs, jobModel)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
