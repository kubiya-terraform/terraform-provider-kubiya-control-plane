package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
	"kubiya-control-plane/internal/entities"
)

var _ resource.Resource = (*teamResource)(nil)
var _ resource.ResourceWithImportState = (*teamResource)(nil)

func NewTeamResource() resource.Resource {
	return &teamResource{}
}

type teamResource struct {
	client *clients.Client
}

type teamResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	Configuration        types.String `tfsdk:"configuration"`
	ToolsetIDs           types.List   `tfsdk:"toolset_ids"`
	ExecutionEnvironment types.String `tfsdk:"execution_environment"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (r *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (r *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Team in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Team ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Team name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Team description",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Team status (active, inactive, archived)",
				Computed:    true,
				Optional:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "Team configuration as JSON string",
				Optional:    true,
			},
			"toolset_ids": schema.ListAttribute{
				Description: "List of toolset IDs associated with the team",
				Optional:    true,
				ElementType: types.StringType,
			},
			"execution_environment": schema.StringAttribute{
				Description: "Execution environment configuration as JSON string",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the team was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the team was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build create request
	createReq := &entities.TeamCreateRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.Configuration.IsNull() {
		config, err := parseJSON(plan.Configuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to parse configuration JSON: %s", err))
			return
		}
		createReq.Configuration = config
	}

	if !plan.ToolsetIDs.IsNull() {
		var toolsetIDs []string
		diags = plan.ToolsetIDs.ElementsAs(ctx, &toolsetIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.ToolsetIDs = toolsetIDs
	}

	if !plan.ExecutionEnvironment.IsNull() {
		execEnv, err := parseJSON(plan.ExecutionEnvironment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Execution Environment", fmt.Sprintf("Failed to parse execution_environment JSON: %s", err))
			return
		}
		createReq.ExecutionEnvironment = execEnv
	}

	// Create team
	team, err := r.client.CreateTeam(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating team", err.Error())
		return
	}

	// Map response to state
	plan.ID = types.StringValue(team.ID)
	plan.Name = types.StringValue(team.Name)

	if team.Description != nil {
		plan.Description = types.StringValue(*team.Description)
	}

	if team.Status != "" {
		plan.Status = types.StringValue(string(team.Status))
	} else {
		plan.Status = types.StringNull()
	}

	if team.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(team.CreatedAt.String())
	}

	if team.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(team.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.client.GetTeam(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading team", err.Error())
		return
	}

	// Update state
	state.Name = types.StringValue(team.Name)

	if team.Description != nil {
		state.Description = types.StringValue(*team.Description)
	}

	if team.Status != "" {
		state.Status = types.StringValue(string(team.Status))
	} else {
		state.Status = types.StringNull()
	}

	if team.CreatedAt != nil {
		state.CreatedAt = types.StringValue(team.CreatedAt.String())
	}

	if team.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(team.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state teamResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request
	updateReq := &entities.TeamUpdateRequest{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	if !plan.Status.Equal(state.Status) && !plan.Status.IsNull() {
		status := entities.TeamStatus(plan.Status.ValueString())
		updateReq.Status = &status
	}

	// Update team
	team, err := r.client.UpdateTeam(state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating team", err.Error())
		return
	}

	// Update all computed fields from response
	plan.ID = types.StringValue(team.ID)
	plan.Name = types.StringValue(team.Name)

	if team.Description != nil {
		plan.Description = types.StringValue(*team.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if team.Status != "" {
		plan.Status = types.StringValue(string(team.Status))
	} else {
		plan.Status = types.StringNull()
	}

	if team.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(team.CreatedAt.String())
	}

	if team.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(team.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state teamResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTeam(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting team", err.Error())
		return
	}
}

func (r *teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
