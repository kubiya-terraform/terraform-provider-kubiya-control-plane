package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*teamDataSource)(nil)

func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

type teamDataSource struct {
	client *clients.Client
}

type teamDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	Configuration        types.String `tfsdk:"configuration"`
	SkillIDs             types.List   `tfsdk:"skill_ids"`
	ExecutionEnvironment types.String `tfsdk:"execution_environment"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Team from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Team ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Team name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Team description",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Team status (active, inactive, archived)",
				Computed:    true,
			},
			"configuration": schema.StringAttribute{
				Description: "Team configuration as JSON string",
				Computed:    true,
			},
			"skill_ids": schema.ListAttribute{
				Description: "List of skill IDs associated with the team",
				Computed:    true,
				ElementType: types.StringType,
			},
			"execution_environment": schema.StringAttribute{
				Description: "Execution environment configuration as JSON string",
				Computed:    true,
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

func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config teamDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := d.client.GetTeam(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading team", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(team.Name)

	if team.Description != nil {
		config.Description = types.StringValue(*team.Description)
	} else {
		config.Description = types.StringNull()
	}

	config.Status = types.StringValue(string(team.Status))

	// Convert configuration to JSON string
	if len(team.Configuration) > 0 {
		configJSON, err := toJSONString(team.Configuration)
		if err != nil {
			resp.Diagnostics.AddError("Error converting configuration", err.Error())
			return
		}
		config.Configuration = types.StringValue(configJSON)
	} else {
		config.Configuration = types.StringNull()
	}

	// Convert skill IDs to list
	if len(team.SkillIDs) > 0 {
		skillList := make([]types.String, len(team.SkillIDs))
		for i, id := range team.SkillIDs {
			skillList[i] = types.StringValue(id)
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, skillList)
		resp.Diagnostics.Append(diags...)
		config.SkillIDs = listVal
	} else {
		config.SkillIDs = types.ListNull(types.StringType)
	}

	// Convert execution environment to JSON string
	if len(team.ExecutionEnvironment) > 0 {
		execEnvJSON, err := toJSONString(team.ExecutionEnvironment)
		if err != nil {
			resp.Diagnostics.AddError("Error converting execution_environment", err.Error())
			return
		}
		config.ExecutionEnvironment = types.StringValue(execEnvJSON)
	} else {
		config.ExecutionEnvironment = types.StringNull()
	}

	if team.CreatedAt != nil {
		config.CreatedAt = types.StringValue(team.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if team.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(team.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
