package provider

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
	"terraform-provider-kubiya-control-plane/internal/entities"
	kubiyasentry "terraform-provider-kubiya-control-plane/internal/sentry"
)

var _ resource.Resource = (*agentResource)(nil)
var _ resource.ResourceWithImportState = (*agentResource)(nil)

func NewAgentResource() resource.Resource {
	return &agentResource{}
}

type agentResource struct {
	client *clients.Client
}

type agentResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Status        types.String `tfsdk:"status"`
	Capabilities  types.List   `tfsdk:"capabilities"`
	Configuration types.String `tfsdk:"configuration"`
	ModelID       types.String `tfsdk:"model_id"`
	LLMConfig     types.String `tfsdk:"llm_config"`
	Runtime       types.String `tfsdk:"runtime"`
	TeamID        types.String `tfsdk:"team_id"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (r *agentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *agentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Agent in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Agent ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Agent name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Agent description",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Agent status (idle, running, paused, completed, failed, stopped)",
				Computed:    true,
				Optional:    true,
			},
			"capabilities": schema.ListAttribute{
				Description: "List of agent capabilities",
				Optional:    true,
				ElementType: types.StringType,
			},
			"configuration": schema.StringAttribute{
				Description: "Agent configuration as JSON string",
				Optional:    true,
			},
			"model_id": schema.StringAttribute{
				Description: "LiteLLM model identifier",
				Optional:    true,
			},
			"llm_config": schema.StringAttribute{
				Description: "LLM configuration as JSON string (temperature, top_p, etc.)",
				Optional:    true,
			},
			"runtime": schema.StringAttribute{
				Description: "Runtime type (default or claude_code)",
				Optional:    true,
			},
			"team_id": schema.StringAttribute{
				Description: "Team ID to assign this agent to",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the agent was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the agent was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *agentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *agentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Start tracing for this operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "agent", "", kubiyasentry.OpResourceCreate)
	defer kubiyasentry.FinishSpan(span)

	logger := kubiyasentry.LoggerFromContext(ctx)
	logger.Info("Creating agent resource")

	var plan agentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInvalidArgument)
		return
	}

	// Build create request
	createReq := &entities.AgentCreateRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	// Handle capabilities
	if !plan.Capabilities.IsNull() {
		var capabilities []string
		diags = plan.Capabilities.ElementsAs(ctx, &capabilities, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Capabilities = capabilities
	}

	if !plan.Configuration.IsNull() {
		config, err := parseJSON(plan.Configuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to parse configuration JSON: %s", err))
			return
		}
		createReq.Configuration = config
	}

	if !plan.ModelID.IsNull() {
		modelID := plan.ModelID.ValueString()
		createReq.ModelID = &modelID
	}

	if !plan.LLMConfig.IsNull() {
		llmConfig, err := parseJSON(plan.LLMConfig.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid LLM Config", fmt.Sprintf("Failed to parse llm_config JSON: %s", err))
			return
		}
		createReq.LLMConfig = llmConfig
	}

	if !plan.Runtime.IsNull() {
		runtime := entities.RuntimeType(plan.Runtime.ValueString())
		createReq.Runtime = &runtime
	}

	if !plan.TeamID.IsNull() {
		teamID := plan.TeamID.ValueString()
		createReq.TeamID = &teamID
	}

	// Create agent
	agent, err := r.client.CreateAgent(createReq)
	if err != nil {
		logger.Error("Failed to create agent", "error", err)
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Error creating agent", err.Error())
		return
	}

	logger.Info("Successfully created agent", "agent_id", agent.ID)

	// Map response to state
	plan.ID = types.StringValue(agent.ID)
	plan.Name = types.StringValue(agent.Name)

	if agent.Description != nil {
		plan.Description = types.StringValue(*agent.Description)
	}

	if agent.Status != "" {
		plan.Status = types.StringValue(string(agent.Status))
	} else {
		plan.Status = types.StringNull()
	}

	// Set optional fields from the response if available
	if agent.ModelID != nil {
		plan.ModelID = types.StringValue(*agent.ModelID)
	}

	if agent.Runtime != "" {
		plan.Runtime = types.StringValue(string(agent.Runtime))
	}

	if agent.TeamID != nil {
		plan.TeamID = types.StringValue(*agent.TeamID)
	}

	// Keep the input values for fields not returned by API
	// Capabilities, Configuration, and LLMConfig are preserved from plan

	if agent.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(agent.CreatedAt.String())
	}

	if agent.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(agent.UpdatedAt.String())
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
	kubiyasentry.SetSpanTag(span, kubiyasentry.TagResourceID, agent.ID)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *agentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state agentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start tracing for this operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "agent", state.ID.ValueString(), kubiyasentry.OpResourceRead)
	defer kubiyasentry.FinishSpan(span)

	logger := kubiyasentry.LoggerFromContext(ctx)
	logger.Debug("Reading agent resource", "agent_id", state.ID.ValueString())

	agent, err := r.client.GetAgent(state.ID.ValueString())
	if err != nil {
		logger.Error("Failed to read agent", "error", err, "agent_id", state.ID.ValueString())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Error reading agent", err.Error())
		return
	}

	// Update state
	state.Name = types.StringValue(agent.Name)

	if agent.Description != nil {
		state.Description = types.StringValue(*agent.Description)
	}

	if agent.Status != "" {
		state.Status = types.StringValue(string(agent.Status))
	} else {
		state.Status = types.StringNull()
	}

	// Set optional fields from the response if available
	if agent.ModelID != nil {
		state.ModelID = types.StringValue(*agent.ModelID)
	}

	if agent.Runtime != "" {
		state.Runtime = types.StringValue(string(agent.Runtime))
	}

	if agent.TeamID != nil {
		state.TeamID = types.StringValue(*agent.TeamID)
	}

	// Keep existing values for fields not returned by API
	// Capabilities, Configuration, and LLMConfig are preserved from existing state

	if agent.CreatedAt != nil {
		state.CreatedAt = types.StringValue(agent.CreatedAt.String())
	}

	if agent.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(agent.UpdatedAt.String())
	}

	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *agentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan agentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state agentResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start tracing for this operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "agent", state.ID.ValueString(), kubiyasentry.OpResourceUpdate)
	defer kubiyasentry.FinishSpan(span)

	logger := kubiyasentry.LoggerFromContext(ctx)
	logger.Info("Updating agent resource", "agent_id", state.ID.ValueString())

	// Build update request
	updateReq := &entities.AgentUpdateRequest{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	if !plan.Status.Equal(state.Status) && !plan.Status.IsNull() {
		status := entities.AgentStatus(plan.Status.ValueString())
		updateReq.Status = &status
	}

	if !plan.ModelID.Equal(state.ModelID) {
		modelID := plan.ModelID.ValueString()
		updateReq.ModelID = &modelID
	}

	if !plan.Runtime.Equal(state.Runtime) && !plan.Runtime.IsNull() {
		runtime := entities.RuntimeType(plan.Runtime.ValueString())
		updateReq.Runtime = &runtime
	}

	if !plan.TeamID.Equal(state.TeamID) {
		teamID := plan.TeamID.ValueString()
		updateReq.TeamID = &teamID
	}

	// Update agent
	agent, err := r.client.UpdateAgent(state.ID.ValueString(), updateReq)
	if err != nil {
		logger.Error("Failed to update agent", "error", err, "agent_id", state.ID.ValueString())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Error updating agent", err.Error())
		return
	}

	logger.Info("Successfully updated agent", "agent_id", state.ID.ValueString())
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)

	// Update all computed fields from response
	plan.ID = types.StringValue(agent.ID)
	plan.Name = types.StringValue(agent.Name)

	if agent.Description != nil {
		plan.Description = types.StringValue(*agent.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if agent.Status != "" {
		plan.Status = types.StringValue(string(agent.Status))
	} else {
		plan.Status = types.StringNull()
	}

	if agent.ModelID != nil {
		plan.ModelID = types.StringValue(*agent.ModelID)
	}

	if agent.Runtime != "" {
		plan.Runtime = types.StringValue(string(agent.Runtime))
	}

	if agent.TeamID != nil {
		plan.TeamID = types.StringValue(*agent.TeamID)
	}

	if agent.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(agent.CreatedAt.String())
	}

	if agent.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(agent.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *agentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state agentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Start tracing for this operation
	ctx, span := kubiyasentry.TraceResourceOperation(ctx, "agent", state.ID.ValueString(), kubiyasentry.OpResourceDelete)
	defer kubiyasentry.FinishSpan(span)

	logger := kubiyasentry.LoggerFromContext(ctx)
	logger.Info("Deleting agent resource", "agent_id", state.ID.ValueString())

	err := r.client.DeleteAgent(state.ID.ValueString())
	if err != nil {
		logger.Error("Failed to delete agent", "error", err, "agent_id", state.ID.ValueString())
		kubiyasentry.RecordError(ctx, err)
		kubiyasentry.SetSpanStatus(span, sentry.SpanStatusInternalError)
		resp.Diagnostics.AddError("Error deleting agent", err.Error())
		return
	}

	logger.Info("Successfully deleted agent", "agent_id", state.ID.ValueString())
	kubiyasentry.SetSpanStatus(span, sentry.SpanStatusOK)
}

func (r *agentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
