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

var _ resource.Resource = (*toolsetResource)(nil)
var _ resource.ResourceWithImportState = (*toolsetResource)(nil)

func NewToolSetResource() resource.Resource {
	return &toolsetResource{}
}

type toolsetResource struct {
	client *clients.Client
}

type toolsetResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	Icon          types.String `tfsdk:"icon"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Configuration types.String `tfsdk:"configuration"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (r *toolsetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_toolset"
}

func (r *toolsetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a ToolSet in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ToolSet ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "ToolSet name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "ToolSet type (file_system, shell, docker, python, file_generation, custom)",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "ToolSet description",
				Optional:    true,
			},
			"icon": schema.StringAttribute{
				Description: "Icon name",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the toolset is enabled",
				Optional:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "ToolSet configuration as JSON string",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the toolset was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the toolset was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *toolsetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *toolsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan toolsetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.ToolSetCreateRequest{
		Name:    plan.Name.ValueString(),
		Type:    entities.ToolSetType(plan.Type.ValueString()),
		Enabled: plan.Enabled.ValueBool(),
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.Icon.IsNull() {
		createReq.Icon = plan.Icon.ValueString()
	}

	if !plan.Configuration.IsNull() {
		config, err := parseJSON(plan.Configuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to parse configuration JSON: %s", err))
			return
		}
		createReq.Configuration = config
	}

	toolset, err := r.client.CreateToolSet(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating toolset", err.Error())
		return
	}

	plan.ID = types.StringValue(toolset.ID)
	plan.Name = types.StringValue(toolset.Name)
	plan.Type = types.StringValue(string(toolset.Type))

	if toolset.Description != nil {
		plan.Description = types.StringValue(*toolset.Description)
	}

	if toolset.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(toolset.CreatedAt.String())
	}

	if toolset.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(toolset.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *toolsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state toolsetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toolset, err := r.client.GetToolSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading toolset", err.Error())
		return
	}

	state.ID = types.StringValue(toolset.ID)
	state.Name = types.StringValue(toolset.Name)
	state.Type = types.StringValue(string(toolset.Type))
	state.Enabled = types.BoolValue(toolset.Enabled)

	if toolset.Description != nil {
		state.Description = types.StringValue(*toolset.Description)
	}

	if toolset.CreatedAt != nil {
		state.CreatedAt = types.StringValue(toolset.CreatedAt.String())
	}

	if toolset.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(toolset.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *toolsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan toolsetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &entities.ToolSetUpdateRequest{}

	name := plan.Name.ValueString()
	updateReq.Name = &name

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	enabled := plan.Enabled.ValueBool()
	updateReq.Enabled = &enabled

	if !plan.Configuration.IsNull() {
		config, err := parseJSON(plan.Configuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("Failed to parse configuration JSON: %s", err))
			return
		}
		updateReq.Configuration = config
	}

	toolset, err := r.client.UpdateToolSet(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating toolset", err.Error())
		return
	}

	// Update all computed fields from response
	plan.ID = types.StringValue(toolset.ID)
	plan.Name = types.StringValue(toolset.Name)
	plan.Type = types.StringValue(string(toolset.Type))
	plan.Enabled = types.BoolValue(toolset.Enabled)

	if toolset.Description != nil {
		plan.Description = types.StringValue(*toolset.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if toolset.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(toolset.CreatedAt.String())
	}

	if toolset.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(toolset.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *toolsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state toolsetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteToolSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting toolset", err.Error())
		return
	}
}

func (r *toolsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
