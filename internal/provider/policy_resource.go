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

var _ resource.Resource = (*policyResource)(nil)
var _ resource.ResourceWithImportState = (*policyResource)(nil)

func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

type policyResource struct {
	client *clients.Client
}

type policyResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	PolicyContent types.String `tfsdk:"policy_content"`
	PolicyType    types.String `tfsdk:"policy_type"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Tags          types.List   `tfsdk:"tags"`
	Version       types.Int64  `tfsdk:"version"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

func (r *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OPA Policy in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Policy ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Policy name",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Policy description",
				Optional:    true,
			},
			"policy_content": schema.StringAttribute{
				Description: "OPA Rego policy content",
				Required:    true,
			},
			"policy_type": schema.StringAttribute{
				Description: "Policy type (rego, json)",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the policy is enabled",
				Optional:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Policy tags",
				Optional:    true,
				ElementType: types.StringType,
			},
			"version": schema.Int64Attribute{
				Description: "Policy version",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the policy was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the policy was last updated",
				Computed:    true,
			},
		},
	}
}

func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.PolicyCreateRequest{
		Name:          plan.Name.ValueString(),
		PolicyContent: plan.PolicyContent.ValueString(),
		Enabled:       plan.Enabled.ValueBool(),
		PolicyType:    entities.PolicyTypeRego,
	}

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		createReq.Description = &desc
	}

	if !plan.PolicyType.IsNull() {
		createReq.PolicyType = entities.PolicyType(plan.PolicyType.ValueString())
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

	policy, err := r.client.CreatePolicy(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating policy", err.Error())
		return
	}

	plan.ID = types.StringValue(policy.ID)
	plan.Name = types.StringValue(policy.Name)
	plan.PolicyContent = types.StringValue(policy.PolicyContent)
	plan.PolicyType = types.StringValue(string(policy.PolicyType))
	plan.Enabled = types.BoolValue(policy.Enabled)
	plan.Version = types.Int64Value(policy.Version)

	if policy.Description != nil {
		plan.Description = types.StringValue(*policy.Description)
	}

	if policy.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(policy.CreatedAt.String())
	}

	if policy.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(policy.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading policy", err.Error())
		return
	}

	state.ID = types.StringValue(policy.ID)
	state.Name = types.StringValue(policy.Name)
	state.PolicyContent = types.StringValue(policy.PolicyContent)
	state.PolicyType = types.StringValue(string(policy.PolicyType))
	state.Enabled = types.BoolValue(policy.Enabled)
	state.Version = types.Int64Value(policy.Version)

	if policy.Description != nil {
		state.Description = types.StringValue(*policy.Description)
	}

	if policy.CreatedAt != nil {
		state.CreatedAt = types.StringValue(policy.CreatedAt.String())
	}

	if policy.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(policy.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &entities.PolicyUpdateRequest{}

	name := plan.Name.ValueString()
	updateReq.Name = &name

	content := plan.PolicyContent.ValueString()
	updateReq.PolicyContent = &content

	enabled := plan.Enabled.ValueBool()
	updateReq.Enabled = &enabled

	if !plan.Description.IsNull() {
		desc := plan.Description.ValueString()
		updateReq.Description = &desc
	}

	if !plan.Tags.IsNull() {
		var tags []string
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Tags = tags
	}

	policy, err := r.client.UpdatePolicy(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating policy", err.Error())
		return
	}

	// Update all computed fields from response
	plan.ID = types.StringValue(policy.ID)
	plan.Name = types.StringValue(policy.Name)
	plan.PolicyContent = types.StringValue(policy.PolicyContent)
	plan.PolicyType = types.StringValue(string(policy.PolicyType))
	plan.Enabled = types.BoolValue(policy.Enabled)
	plan.Version = types.Int64Value(policy.Version)

	if policy.Description != nil {
		plan.Description = types.StringValue(*policy.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if policy.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(policy.CreatedAt.String())
	}

	if policy.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(policy.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting policy", err.Error())
		return
	}
}

func (r *policyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
