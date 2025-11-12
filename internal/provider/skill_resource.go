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

var _ resource.Resource = (*skillResource)(nil)
var _ resource.ResourceWithImportState = (*skillResource)(nil)

func NewSkillResource() resource.Resource {
	return &skillResource{}
}

type skillResource struct {
	client *clients.Client
}

type skillResourceModel struct {
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

func (r *skillResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (r *skillResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Skill in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Skill ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Skill name",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Skill type (file_system, shell, docker, python, file_generation, custom)",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Skill description",
				Optional:    true,
			},
			"icon": schema.StringAttribute{
				Description: "Icon name",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the skill is enabled",
				Optional:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "Skill configuration as JSON string",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the skill was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the skill was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *skillResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *skillResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.SkillCreateRequest{
		Name:    plan.Name.ValueString(),
		Type:    entities.SkillType(plan.Type.ValueString()),
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

	skill, err := r.client.CreateSkill(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating skill", err.Error())
		return
	}

	plan.ID = types.StringValue(skill.ID)
	plan.Name = types.StringValue(skill.Name)
	plan.Type = types.StringValue(string(skill.Type))

	if skill.Description != nil {
		plan.Description = types.StringValue(*skill.Description)
	}

	if skill.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(skill.CreatedAt.String())
	}

	if skill.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(skill.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skillResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := r.client.GetSkill(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill", err.Error())
		return
	}

	state.ID = types.StringValue(skill.ID)
	state.Name = types.StringValue(skill.Name)
	state.Type = types.StringValue(string(skill.Type))
	state.Enabled = types.BoolValue(skill.Enabled)

	if skill.Description != nil {
		state.Description = types.StringValue(*skill.Description)
	}

	if skill.CreatedAt != nil {
		state.CreatedAt = types.StringValue(skill.CreatedAt.String())
	}

	if skill.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(skill.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *skillResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &entities.SkillUpdateRequest{}

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

	skill, err := r.client.UpdateSkill(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating skill", err.Error())
		return
	}

	// Update all computed fields from response
	plan.ID = types.StringValue(skill.ID)
	plan.Name = types.StringValue(skill.Name)
	plan.Type = types.StringValue(string(skill.Type))
	plan.Enabled = types.BoolValue(skill.Enabled)

	if skill.Description != nil {
		plan.Description = types.StringValue(*skill.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if skill.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(skill.CreatedAt.String())
	}

	if skill.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(skill.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *skillResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSkill(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting skill", err.Error())
		return
	}
}

func (r *skillResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
