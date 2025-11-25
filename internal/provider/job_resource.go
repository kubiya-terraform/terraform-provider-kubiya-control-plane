package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
	"terraform-provider-kubiya-control-plane/internal/entities"
)

var _ resource.Resource = (*jobResource)(nil)
var _ resource.ResourceWithImportState = (*jobResource)(nil)

func NewJobResource() resource.Resource {
	return &jobResource{}
}

type jobResource struct {
	client *clients.Client
}

type jobResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Enabled               types.Bool   `tfsdk:"enabled"`
	Status                types.String `tfsdk:"status"`
	TriggerType           types.String `tfsdk:"trigger_type"`
	CronSchedule          types.String `tfsdk:"cron_schedule"`
	CronTimezone          types.String `tfsdk:"cron_timezone"`
	WebhookURL            types.String `tfsdk:"webhook_url"`
	WebhookSecret         types.String `tfsdk:"webhook_secret"`
	PlanningMode          types.String `tfsdk:"planning_mode"`
	EntityType            types.String `tfsdk:"entity_type"`
	EntityID              types.String `tfsdk:"entity_id"`
	PromptTemplate        types.String `tfsdk:"prompt_template"`
	SystemPrompt          types.String `tfsdk:"system_prompt"`
	ExecutorType          types.String `tfsdk:"executor_type"`
	WorkerQueueName       types.String `tfsdk:"worker_queue_name"`
	EnvironmentName       types.String `tfsdk:"environment_name"`
	Config                types.String `tfsdk:"config"`
	ExecutionEnvVars      types.Map    `tfsdk:"execution_env_vars"`
	ExecutionSecrets      types.List   `tfsdk:"execution_secrets"`
	ExecutionIntegrations types.List   `tfsdk:"execution_integrations"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *jobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_job"
}

func (r *jobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Job in the Control Plane. Jobs can be triggered by cron schedules, webhooks, or manually.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Job ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Job name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Job description",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the job is enabled",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"status": schema.StringAttribute{
				Description: "Job status",
				Computed:    true,
			},
			"trigger_type": schema.StringAttribute{
				Description: "Trigger type: 'cron', 'webhook', or 'manual'",
				Required:    true,
			},
			"cron_schedule": schema.StringAttribute{
				Description: "Cron expression in 5-field format: 'minute hour day month weekday' (e.g., '0 17 * * *' for daily at 5pm, '0 9 * * 1-5' for weekdays at 9am). Required when trigger_type is 'cron'. Day of week: 0=Sunday, 1=Monday, ..., 6=Saturday",
				Optional:    true,
			},
			"cron_timezone": schema.StringAttribute{
				Description: "Timezone for cron schedule (e.g., 'America/New_York')",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("UTC"),
			},
			"webhook_url": schema.StringAttribute{
				Description: "Full webhook URL (generated for webhook triggers)",
				Computed:    true,
			},
			"webhook_secret": schema.StringAttribute{
				Description: "Webhook HMAC secret for signature verification",
				Computed:    true,
				Sensitive:   true,
			},
			"planning_mode": schema.StringAttribute{
				Description: "Planning mode: 'on_the_fly', 'predefined_agent', 'predefined_team', or 'predefined_workflow'",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("predefined_agent"),
			},
			"entity_type": schema.StringAttribute{
				Description: "Entity type: 'agent', 'team', or 'workflow'. Required when planning_mode is not 'on_the_fly'",
				Optional:    true,
			},
			"entity_id": schema.StringAttribute{
				Description: "Entity ID (agent_id, team_id, or workflow_id). Required when planning_mode is not 'on_the_fly'",
				Optional:    true,
			},
			"prompt_template": schema.StringAttribute{
				Description: "Prompt template (can include {{variables}} for dynamic params)",
				Required:    true,
			},
			"system_prompt": schema.StringAttribute{
				Description: "Optional system prompt",
				Optional:    true,
			},
			"executor_type": schema.StringAttribute{
				Description: "Executor routing: 'auto', 'specific_queue', or 'environment'",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("auto"),
			},
			"worker_queue_name": schema.StringAttribute{
				Description: "Worker queue name for 'specific_queue' executor type",
				Optional:    true,
			},
			"environment_name": schema.StringAttribute{
				Description: "Environment name for 'environment' executor type",
				Optional:    true,
			},
			"config": schema.StringAttribute{
				Description: "Additional execution config as JSON string (timeout, retry, etc.)",
				Optional:    true,
			},
			"execution_env_vars": schema.MapAttribute{
				Description: "Execution environment variables",
				Optional:    true,
				ElementType: types.StringType,
			},
			"execution_secrets": schema.ListAttribute{
				Description: "List of secret names to inject into execution environment",
				Optional:    true,
				ElementType: types.StringType,
			},
			"execution_integrations": schema.ListAttribute{
				Description: "List of integration IDs to inject into execution environment",
				Optional:    true,
				ElementType: types.StringType,
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

func (r *jobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*clients.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *clients.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *jobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan jobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.JobCreateRequest{
		Name:           plan.Name.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
		TriggerType:    plan.TriggerType.ValueString(),
		PlanningMode:   plan.PlanningMode.ValueString(),
		PromptTemplate: plan.PromptTemplate.ValueString(),
		ExecutorType:   plan.ExecutorType.ValueString(),
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.CronSchedule.IsNull() {
		schedule := plan.CronSchedule.ValueString()
		createReq.CronSchedule = &schedule
	}

	if !plan.CronTimezone.IsNull() {
		tz := plan.CronTimezone.ValueString()
		createReq.CronTimezone = &tz
	}

	if !plan.EntityType.IsNull() {
		et := plan.EntityType.ValueString()
		createReq.EntityType = &et
	}

	if !plan.EntityID.IsNull() {
		eid := plan.EntityID.ValueString()
		createReq.EntityID = &eid
	}

	if !plan.SystemPrompt.IsNull() {
		sp := plan.SystemPrompt.ValueString()
		createReq.SystemPrompt = &sp
	}

	if !plan.WorkerQueueName.IsNull() {
		wqn := plan.WorkerQueueName.ValueString()
		createReq.WorkerQueueName = &wqn
	}

	if !plan.EnvironmentName.IsNull() {
		en := plan.EnvironmentName.ValueString()
		createReq.EnvironmentName = &en
	}

	if !plan.Config.IsNull() {
		config, err := parseJSON(plan.Config.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Config", fmt.Sprintf("Failed to parse config JSON: %s", err))
			return
		}
		createReq.Config = config
	}

	// Build execution environment
	if !plan.ExecutionEnvVars.IsNull() || !plan.ExecutionSecrets.IsNull() || !plan.ExecutionIntegrations.IsNull() {
		createReq.ExecutionEnv = &entities.ExecutionEnvironment{}

		if !plan.ExecutionEnvVars.IsNull() {
			envVars := make(map[string]string)
			diags := plan.ExecutionEnvVars.ElementsAs(ctx, &envVars, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			createReq.ExecutionEnv.EnvVars = envVars
		}

		if !plan.ExecutionSecrets.IsNull() {
			var secrets []string
			diags := plan.ExecutionSecrets.ElementsAs(ctx, &secrets, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			createReq.ExecutionEnv.Secrets = secrets
		}

		if !plan.ExecutionIntegrations.IsNull() {
			var integrations []string
			diags := plan.ExecutionIntegrations.ElementsAs(ctx, &integrations, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			createReq.ExecutionEnv.IntegrationIDs = integrations
		}
	}

	job, err := r.client.CreateJob(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating job", err.Error())
		return
	}

	// Update state from response
	r.updateModelFromJob(&plan, job)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *jobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state jobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	job, err := r.client.GetJob(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading job", err.Error())
		return
	}

	r.updateModelFromJob(&state, job)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *jobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan jobResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &entities.JobUpdateRequest{}

	name := plan.Name.ValueString()
	updateReq.Name = &name

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	enabled := plan.Enabled.ValueBool()
	updateReq.Enabled = &enabled

	triggerType := plan.TriggerType.ValueString()
	updateReq.TriggerType = &triggerType

	if !plan.CronSchedule.IsNull() {
		schedule := plan.CronSchedule.ValueString()
		updateReq.CronSchedule = &schedule
	}

	if !plan.CronTimezone.IsNull() {
		tz := plan.CronTimezone.ValueString()
		updateReq.CronTimezone = &tz
	}

	planningMode := plan.PlanningMode.ValueString()
	updateReq.PlanningMode = &planningMode

	if !plan.EntityType.IsNull() {
		et := plan.EntityType.ValueString()
		updateReq.EntityType = &et
	}

	if !plan.EntityID.IsNull() {
		eid := plan.EntityID.ValueString()
		updateReq.EntityID = &eid
	}

	promptTemplate := plan.PromptTemplate.ValueString()
	updateReq.PromptTemplate = &promptTemplate

	if !plan.SystemPrompt.IsNull() {
		sp := plan.SystemPrompt.ValueString()
		updateReq.SystemPrompt = &sp
	}

	executorType := plan.ExecutorType.ValueString()
	updateReq.ExecutorType = &executorType

	if !plan.WorkerQueueName.IsNull() {
		wqn := plan.WorkerQueueName.ValueString()
		updateReq.WorkerQueueName = &wqn
	}

	if !plan.EnvironmentName.IsNull() {
		en := plan.EnvironmentName.ValueString()
		updateReq.EnvironmentName = &en
	}

	if !plan.Config.IsNull() {
		config, err := parseJSON(plan.Config.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Config", fmt.Sprintf("Failed to parse config JSON: %s", err))
			return
		}
		updateReq.Config = config
	}

	// Build execution environment
	if !plan.ExecutionEnvVars.IsNull() || !plan.ExecutionSecrets.IsNull() || !plan.ExecutionIntegrations.IsNull() {
		updateReq.ExecutionEnv = &entities.ExecutionEnvironment{}

		if !plan.ExecutionEnvVars.IsNull() {
			envVars := make(map[string]string)
			diags := plan.ExecutionEnvVars.ElementsAs(ctx, &envVars, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			updateReq.ExecutionEnv.EnvVars = envVars
		}

		if !plan.ExecutionSecrets.IsNull() {
			var secrets []string
			diags := plan.ExecutionSecrets.ElementsAs(ctx, &secrets, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			updateReq.ExecutionEnv.Secrets = secrets
		}

		if !plan.ExecutionIntegrations.IsNull() {
			var integrations []string
			diags := plan.ExecutionIntegrations.ElementsAs(ctx, &integrations, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
			updateReq.ExecutionEnv.IntegrationIDs = integrations
		}
	}

	job, err := r.client.UpdateJob(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating job", err.Error())
		return
	}

	r.updateModelFromJob(&plan, job)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *jobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state jobResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteJob(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting job", err.Error())
		return
	}
}

func (r *jobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *jobResource) updateModelFromJob(model *jobResourceModel, job *entities.Job) {
	model.ID = types.StringValue(job.ID)
	model.Name = types.StringValue(job.Name)
	model.Enabled = types.BoolValue(job.Enabled)
	model.TriggerType = types.StringValue(job.TriggerType)
	model.PlanningMode = types.StringValue(job.PlanningMode)
	model.PromptTemplate = types.StringValue(job.PromptTemplate)
	model.ExecutorType = types.StringValue(job.ExecutorType)

	if job.Description != nil {
		model.Description = types.StringValue(*job.Description)
	} else {
		model.Description = types.StringNull()
	}

	if job.Status != "" {
		model.Status = types.StringValue(job.Status)
	}

	if job.CronSchedule != nil {
		model.CronSchedule = types.StringValue(*job.CronSchedule)
	} else {
		model.CronSchedule = types.StringNull()
	}

	if job.CronTimezone != nil {
		model.CronTimezone = types.StringValue(*job.CronTimezone)
	} else {
		model.CronTimezone = types.StringNull()
	}

	if job.WebhookURL != nil {
		model.WebhookURL = types.StringValue(*job.WebhookURL)
	} else {
		model.WebhookURL = types.StringNull()
	}

	if job.WebhookSecret != nil {
		model.WebhookSecret = types.StringValue(*job.WebhookSecret)
	} else {
		model.WebhookSecret = types.StringNull()
	}

	if job.EntityType != nil {
		model.EntityType = types.StringValue(*job.EntityType)
	} else {
		model.EntityType = types.StringNull()
	}

	if job.EntityID != nil {
		model.EntityID = types.StringValue(*job.EntityID)
	} else {
		model.EntityID = types.StringNull()
	}

	if job.SystemPrompt != nil {
		model.SystemPrompt = types.StringValue(*job.SystemPrompt)
	} else {
		model.SystemPrompt = types.StringNull()
	}

	if job.WorkerQueueName != nil {
		model.WorkerQueueName = types.StringValue(*job.WorkerQueueName)
	} else {
		model.WorkerQueueName = types.StringNull()
	}

	if job.EnvironmentName != nil {
		model.EnvironmentName = types.StringValue(*job.EnvironmentName)
	} else {
		model.EnvironmentName = types.StringNull()
	}

	if job.CreatedAt != nil {
		model.CreatedAt = types.StringValue(job.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}

	if job.UpdatedAt != nil {
		model.UpdatedAt = types.StringValue(job.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	}
}
