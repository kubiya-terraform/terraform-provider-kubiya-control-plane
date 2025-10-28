package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
	"kubiya-control-plane/internal/entities"
)

var _ resource.Resource = (*environmentResource)(nil)
var _ resource.ResourceWithImportState = (*environmentResource)(nil)

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

type environmentResource struct {
	client *clients.Client
}

type environmentResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	DisplayName          types.String `tfsdk:"display_name"`
	Description          types.String `tfsdk:"description"`
	Tags                 types.List   `tfsdk:"tags"`
	Settings             types.String `tfsdk:"settings"`
	Status               types.String `tfsdk:"status"`
	ExecutionEnvironment types.String `tfsdk:"execution_environment"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Environment in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Environment ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Environment name (e.g., default, production)",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "User-friendly display name",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Environment description",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags for categorization",
				Optional:    true,
				ElementType: types.StringType,
			},
			"settings": schema.StringAttribute{
				Description: "Environment settings as JSON string",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Environment status (active, inactive)",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("active"),
			},
			"execution_environment": schema.StringAttribute{
				Description: "Execution environment configuration as JSON string (env_vars, secrets, integration_ids)",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the environment was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the environment was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build create request
	createReq := &entities.EnvironmentCreateRequest{
		Name: plan.Name.ValueString(),
	}

	if !plan.DisplayName.IsNull() {
		displayName := plan.DisplayName.ValueString()
		createReq.DisplayName = &displayName
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.Tags.IsNull() {
		var tags []string
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Tags = tags
	}

	if !plan.Settings.IsNull() {
		settings, err := parseJSON(plan.Settings.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Settings", fmt.Sprintf("Failed to parse settings JSON: %s", err))
			return
		}
		createReq.Settings = settings
	}

	if !plan.ExecutionEnvironment.IsNull() {
		execEnv, err := parseJSON(plan.ExecutionEnvironment.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Execution Environment", fmt.Sprintf("Failed to parse execution_environment JSON: %s", err))
			return
		}
		createReq.ExecutionEnvironment = execEnv
	}

	// Create environment
	environment, err := r.client.CreateEnvironment(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating environment", err.Error())
		return
	}

	// Map response to state
	plan.ID = types.StringValue(environment.ID)
	plan.Name = types.StringValue(environment.Name)

	if environment.DisplayName != nil {
		plan.DisplayName = types.StringValue(*environment.DisplayName)
	}

	if environment.Description != nil {
		plan.Description = types.StringValue(*environment.Description)
	}

	plan.Status = types.StringValue(string(environment.Status))

	if environment.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(environment.CreatedAt.String())
	}

	if environment.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(environment.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment, err := r.client.GetEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	// Update state
	state.Name = types.StringValue(environment.Name)

	if environment.DisplayName != nil {
		state.DisplayName = types.StringValue(*environment.DisplayName)
	}

	if environment.Description != nil {
		state.Description = types.StringValue(*environment.Description)
	}

	state.Status = types.StringValue(string(environment.Status))

	if environment.CreatedAt != nil {
		state.CreatedAt = types.StringValue(environment.CreatedAt.String())
	}

	if environment.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(environment.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan environmentResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state environmentResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build update request
	updateReq := &entities.EnvironmentUpdateRequest{}

	if !plan.Name.Equal(state.Name) {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		displayName := plan.DisplayName.ValueString()
		updateReq.DisplayName = &displayName
	}

	if !plan.Description.Equal(state.Description) {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	if !plan.Status.Equal(state.Status) && !plan.Status.IsNull() {
		status := entities.EnvironmentStatus(plan.Status.ValueString())
		updateReq.Status = &status
	}

	// Update environment
	environment, err := r.client.UpdateEnvironment(state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating environment", err.Error())
		return
	}

	// Update state
	plan.UpdatedAt = types.StringValue(environment.UpdatedAt.String())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEnvironment(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting environment", err.Error())
		return
	}
}

func (r *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
