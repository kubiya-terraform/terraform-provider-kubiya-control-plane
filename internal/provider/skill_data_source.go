package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*skillDataSource)(nil)

func NewSkillDataSource() datasource.DataSource {
	return &skillDataSource{}
}

type skillDataSource struct {
	client *clients.Client
}

type skillDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
	Configuration types.String `tfsdk:"configuration"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (d *skillDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (d *skillDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Skill from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Skill ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Skill name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Skill description",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Skill type",
				Computed:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "Skill configuration as JSON string",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the skill is enabled",
				Computed:    true,
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

func (d *skillDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *skillDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config skillDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := d.client.GetSkill(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading skill", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(skill.Name)

	if skill.Description != nil {
		config.Description = types.StringValue(*skill.Description)
	} else {
		config.Description = types.StringNull()
	}

	config.Type = types.StringValue(string(skill.Type))

	// Convert configuration to JSON string
	if len(skill.Configuration) > 0 {
		configJSON, err := toJSONString(skill.Configuration)
		if err != nil {
			resp.Diagnostics.AddError("Error converting configuration", err.Error())
			return
		}
		config.Configuration = types.StringValue(configJSON)
	} else {
		config.Configuration = types.StringNull()
	}

	config.Enabled = types.BoolValue(skill.Enabled)

	if skill.CreatedAt != nil {
		config.CreatedAt = types.StringValue(skill.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if skill.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(skill.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
