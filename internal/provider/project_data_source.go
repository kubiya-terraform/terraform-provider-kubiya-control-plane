package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*projectDataSource)(nil)

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

type projectDataSource struct {
	client *clients.Client
}

type projectDataSourceModel struct {
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
	AgentCount            types.Int64  `tfsdk:"agent_count"`
	TeamCount             types.Int64  `tfsdk:"team_count"`
	CreatedAt             types.String `tfsdk:"created_at"`
	UpdatedAt             types.String `tfsdk:"updated_at"`
}

func (d *projectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Project from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Project ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Project name",
				Computed:    true,
			},
			"key": schema.StringAttribute{
				Description: "Short project key",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Project description",
				Computed:    true,
			},
			"goals": schema.StringAttribute{
				Description: "Project goals and objectives",
				Computed:    true,
			},
			"settings": schema.StringAttribute{
				Description: "Project settings as JSON string",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Project status",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Project visibility",
				Computed:    true,
			},
			"restrict_to_environment": schema.BoolAttribute{
				Description: "Whether restricted to specific environment",
				Computed:    true,
			},
			"policy_ids": schema.ListAttribute{
				Description: "List of OPA policy IDs",
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_model": schema.StringAttribute{
				Description: "Default LLM model",
				Computed:    true,
			},
			"agent_count": schema.Int64Attribute{
				Description: "Number of agents in this project",
				Computed:    true,
			},
			"team_count": schema.Int64Attribute{
				Description: "Number of teams in this project",
				Computed:    true,
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

func (d *projectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config projectDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.client.GetProject(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading project", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(project.Name)
	config.Key = types.StringValue(project.Key)

	if project.Description != nil {
		config.Description = types.StringValue(*project.Description)
	} else {
		config.Description = types.StringNull()
	}

	if project.Goals != nil {
		config.Goals = types.StringValue(*project.Goals)
	} else {
		config.Goals = types.StringNull()
	}

	// Convert settings to JSON string
	if project.Settings != nil && len(project.Settings) > 0 {
		settingsJSON, err := toJSONString(project.Settings)
		if err != nil {
			resp.Diagnostics.AddError("Error converting settings", err.Error())
			return
		}
		config.Settings = types.StringValue(settingsJSON)
	} else {
		config.Settings = types.StringNull()
	}

	config.Status = types.StringValue(string(project.Status))
	config.Visibility = types.StringValue(project.Visibility)
	config.RestrictToEnvironment = types.BoolValue(project.RestrictToEnvironment)

	// Convert policy IDs to list
	if len(project.PolicyIDs) > 0 {
		policyList := make([]types.String, len(project.PolicyIDs))
		for i, id := range project.PolicyIDs {
			policyList[i] = types.StringValue(id)
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, policyList)
		resp.Diagnostics.Append(diags...)
		config.PolicyIDs = listVal
	} else {
		config.PolicyIDs = types.ListNull(types.StringType)
	}

	if project.DefaultModel != nil {
		config.DefaultModel = types.StringValue(*project.DefaultModel)
	} else {
		config.DefaultModel = types.StringNull()
	}

	config.AgentCount = types.Int64Value(int64(project.AgentCount))
	config.TeamCount = types.Int64Value(int64(project.TeamCount))

	if project.CreatedAt != nil {
		config.CreatedAt = types.StringValue(project.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if project.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(project.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
