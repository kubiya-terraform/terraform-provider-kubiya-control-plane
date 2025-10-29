package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*toolsetDataSource)(nil)

func NewToolSetDataSource() datasource.DataSource {
	return &toolsetDataSource{}
}

type toolsetDataSource struct {
	client *clients.Client
}

type toolsetDataSourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Type          types.String `tfsdk:"type"`
	Configuration types.String `tfsdk:"configuration"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
}

func (d *toolsetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_toolset"
}

func (d *toolsetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing ToolSet from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ToolSet ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "ToolSet name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "ToolSet description",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "ToolSet type",
				Computed:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "ToolSet configuration as JSON string",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the toolset is enabled",
				Computed:    true,
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

func (d *toolsetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *toolsetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config toolsetDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	toolset, err := d.client.GetToolSet(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading toolset", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(toolset.Name)

	if toolset.Description != nil {
		config.Description = types.StringValue(*toolset.Description)
	} else {
		config.Description = types.StringNull()
	}

	config.Type = types.StringValue(string(toolset.Type))

	// Convert configuration to JSON string
	if toolset.Configuration != nil && len(toolset.Configuration) > 0 {
		configJSON, err := toJSONString(toolset.Configuration)
		if err != nil {
			resp.Diagnostics.AddError("Error converting configuration", err.Error())
			return
		}
		config.Configuration = types.StringValue(configJSON)
	} else {
		config.Configuration = types.StringNull()
	}

	config.Enabled = types.BoolValue(toolset.Enabled)

	if toolset.CreatedAt != nil {
		config.CreatedAt = types.StringValue(toolset.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if toolset.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(toolset.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
