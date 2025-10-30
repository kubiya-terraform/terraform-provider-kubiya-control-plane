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

	"kubiya-control-plane/internal/clients"
	"kubiya-control-plane/internal/entities"
)

var _ resource.Resource = (*projectResource)(nil)
var _ resource.ResourceWithImportState = (*projectResource)(nil)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

type projectResource struct {
	client *clients.Client
}

type projectResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Key                   types.String `tfsdk:"key"`
	Description           types.String `tfsdk:"description"`
	Goals                 types.String `tfsdk:"goals"`
	Settings              types.String `tfsdk:"settings"`
	Status                types.String `tfsdk:"status"`
	Visibility            types.String `tfsdk:"visibility"`
	RestrictToEnvironment types.Bool   `tfsdk:"restrict_to_environment"`
	PolicyIDs             types.List   `tfsdk:"policy_ids"`
	DefaultModel          types.String `tfsdk:"default_model"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (r *projectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Project in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Required:    true,
			},
			"key": schema.StringAttribute{
				Description: "Short project key (e.g., JIRA, PROJ)",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Project description",
				Optional:    true,
			},
			"goals": schema.StringAttribute{
				Description: "Project goals and objectives",
				Optional:    true,
			},
			"settings": schema.StringAttribute{
				Description: "Project settings as JSON string",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Project status (active, archived, paused)",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("active"),
			},
			"visibility": schema.StringAttribute{
				Description: "Project visibility (private or org)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("private"),
			},
			"restrict_to_environment": schema.BoolAttribute{
				Description: "Restrict to specific runners/environment",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"policy_ids": schema.ListAttribute{
				Description: "List of OPA policy IDs for access control",
				Optional:    true,
				ElementType: types.StringType,
			},
			"default_model": schema.StringAttribute{
				Description: "Default LLM model for this project",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the project was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the project was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *projectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build create request
	createReq := &entities.ProjectCreateRequest{
		Name:                  plan.Name.ValueString(),
		Key:                   plan.Key.ValueString(),
		Visibility:            plan.Visibility.ValueString(),
		RestrictToEnvironment: plan.RestrictToEnvironment.ValueBool(),
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.Goals.IsNull() {
		goals := plan.Goals.ValueString()
		createReq.Goals = &goals
	}

	if !plan.Settings.IsNull() {
		settings, err := parseJSON(plan.Settings.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Settings", fmt.Sprintf("Failed to parse settings JSON: %s", err))
			return
		}
		createReq.Settings = settings
	}

	if !plan.PolicyIDs.IsNull() {
		var policyIDs []string
		diags = plan.PolicyIDs.ElementsAs(ctx, &policyIDs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.PolicyIDs = policyIDs
	}

	if !plan.DefaultModel.IsNull() {
		model := plan.DefaultModel.ValueString()
		createReq.DefaultModel = &model
	}

	// Create project
	project, err := r.client.CreateProject(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating project", err.Error())
		return
	}

	// Map response to state
	plan.ID = types.StringValue(project.ID)
	plan.Name = types.StringValue(project.Name)
	plan.Key = types.StringValue(project.Key)

	if project.Description != nil {
		plan.Description = types.StringValue(*project.Description)
	}

	if project.Goals != nil {
		plan.Goals = types.StringValue(*project.Goals)
	}

	if project.Status != "" {
		plan.Status = types.StringValue(string(project.Status))
	} else {
		plan.Status = types.StringNull()
	}
	plan.Visibility = types.StringValue(project.Visibility)
	plan.RestrictToEnvironment = types.BoolValue(project.RestrictToEnvironment)

	if project.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(project.CreatedAt.String())
	}

	if project.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	// Update state
	state.Name = types.StringValue(project.Name)
	state.Key = types.StringValue(project.Key)

	if project.Description != nil {
		state.Description = types.StringValue(*project.Description)
	}

	if project.Goals != nil {
		state.Goals = types.StringValue(*project.Goals)
	}

	if project.Status != "" {
		state.Status = types.StringValue(string(project.Status))
	} else {
		state.Status = types.StringNull()
	}
	state.Visibility = types.StringValue(project.Visibility)
	state.RestrictToEnvironment = types.BoolValue(project.RestrictToEnvironment)

	if project.CreatedAt != nil {
		state.CreatedAt = types.StringValue(project.CreatedAt.String())
	}

	if project.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state projectResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request
	updateReq := &entities.ProjectUpdateRequest{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.Key.Equal(state.Key) {
		key := plan.Key.ValueString()
		updateReq.Key = &key
	}

	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	if !plan.Goals.Equal(state.Goals) {
		goals := plan.Goals.ValueString()
		updateReq.Goals = &goals
	}

	if !plan.Status.Equal(state.Status) && !plan.Status.IsNull() {
		status := entities.ProjectStatus(plan.Status.ValueString())
		updateReq.Status = &status
	}

	if !plan.Visibility.Equal(state.Visibility) {
		visibility := plan.Visibility.ValueString()
		updateReq.Visibility = &visibility
	}

	if !plan.RestrictToEnvironment.Equal(state.RestrictToEnvironment) {
		restrict := plan.RestrictToEnvironment.ValueBool()
		updateReq.RestrictToEnvironment = &restrict
	}

	if !plan.DefaultModel.Equal(state.DefaultModel) {
		model := plan.DefaultModel.ValueString()
		updateReq.DefaultModel = &model
	}

	// Update project
	project, err := r.client.UpdateProject(state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating project", err.Error())
		return
	}

	// Update all computed fields from response
	plan.ID = types.StringValue(project.ID)
	plan.Name = types.StringValue(project.Name)
	plan.Key = types.StringValue(project.Key)

	if project.Description != nil {
		plan.Description = types.StringValue(*project.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if project.Goals != nil {
		plan.Goals = types.StringValue(*project.Goals)
	} else {
		plan.Goals = types.StringNull()
	}

	if project.Status != "" {
		plan.Status = types.StringValue(string(project.Status))
	} else {
		plan.Status = types.StringNull()
	}

	plan.Visibility = types.StringValue(project.Visibility)
	plan.RestrictToEnvironment = types.BoolValue(project.RestrictToEnvironment)

	if project.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(project.CreatedAt.String())
	}

	if project.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state projectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting project", err.Error())
		return
	}
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
