package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*agentDataSource)(nil)

func NewAgentDataSource() datasource.DataSource {
	return &agentDataSource{}
}

type agentDataSource struct {
	client *clients.Client
}

type agentDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	Status               types.String `tfsdk:"status"`
	Capabilities         types.List   `tfsdk:"capabilities"`
	Configuration        types.String `tfsdk:"configuration"`
	ModelID              types.String `tfsdk:"model_id"`
	LLMConfig            types.String `tfsdk:"llm_config"`
	Runtime              types.String `tfsdk:"runtime"`
	TeamID               types.String `tfsdk:"team_id"`
	SystemPrompt         types.String `tfsdk:"system_prompt"`
	Skills               types.List   `tfsdk:"skills"`
	ExecutionEnvironment types.String `tfsdk:"execution_environment"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (d *agentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (d *agentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Agent from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Agent ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Agent name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Agent description",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Agent status (idle, running, paused, completed, failed, stopped)",
				Computed:    true,
			},
			"capabilities": schema.ListAttribute{
				Description: "List of agent capabilities",
				Computed:    true,
				ElementType: types.StringType,
			},
			"configuration": schema.StringAttribute{
				Description: "Agent configuration as JSON string",
				Computed:    true,
			},
			"model_id": schema.StringAttribute{
				Description: "LiteLLM model identifier",
				Computed:    true,
			},
			"llm_config": schema.StringAttribute{
				Description: "LLM configuration as JSON string",
				Computed:    true,
			},
			"runtime": schema.StringAttribute{
				Description: "Runtime type (default or claude_code)",
				Computed:    true,
			},
			"team_id": schema.StringAttribute{
				Description: "Team ID this agent belongs to",
				Computed:    true,
			},
			"system_prompt": schema.StringAttribute{
				Description: "System prompt for the agent",
				Computed:    true,
			},
			"skills": schema.ListAttribute{
				Description: "List of skills available to the agent",
				Computed:    true,
				ElementType: types.StringType,
			},
			"execution_environment": schema.StringAttribute{
				Description: "Execution environment configuration as JSON string",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Timestamp when the agent was created",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Timestamp when the agent was last updated",
				Computed:    true,
			},
		},
	}
}

func (d *agentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *agentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config agentDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	agent, err := d.client.GetAgent(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading agent", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(agent.Name)

	if agent.Description != nil {
		config.Description = types.StringValue(*agent.Description)
	} else {
		config.Description = types.StringNull()
	}

	config.Status = types.StringValue(string(agent.Status))

	// Convert capabilities to list
	if len(agent.Capabilities) > 0 {
		capList := make([]types.String, len(agent.Capabilities))
		for i, cap := range agent.Capabilities {
			capList[i] = types.StringValue(cap)
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, capList)
		resp.Diagnostics.Append(diags...)
		config.Capabilities = listVal
	} else {
		config.Capabilities = types.ListNull(types.StringType)
	}

	// Convert configuration to JSON string
	if len(agent.Configuration) > 0 {
		configJSON, err := toJSONString(agent.Configuration)
		if err != nil {
			resp.Diagnostics.AddError("Error converting configuration", err.Error())
			return
		}
		config.Configuration = types.StringValue(configJSON)
	} else {
		config.Configuration = types.StringNull()
	}

	if agent.ModelID != nil {
		config.ModelID = types.StringValue(*agent.ModelID)
	} else {
		config.ModelID = types.StringNull()
	}

	// Convert LLM config to JSON string
	if len(agent.LLMConfig) > 0 {
		llmConfigJSON, err := toJSONString(agent.LLMConfig)
		if err != nil {
			resp.Diagnostics.AddError("Error converting llm_config", err.Error())
			return
		}
		config.LLMConfig = types.StringValue(llmConfigJSON)
	} else {
		config.LLMConfig = types.StringNull()
	}

	config.Runtime = types.StringValue(string(agent.Runtime))

	if agent.TeamID != nil {
		config.TeamID = types.StringValue(*agent.TeamID)
	} else {
		config.TeamID = types.StringNull()
	}

	if agent.SystemPrompt != nil {
		config.SystemPrompt = types.StringValue(*agent.SystemPrompt)
	} else {
		config.SystemPrompt = types.StringNull()
	}

	// Convert skills from object to list
	// API returns: {"skill1": {}, "skill2": {}} -> convert to ["skill1", "skill2"]
	if len(agent.Skills) > 0 {
		skillsList := make([]types.String, 0, len(agent.Skills))
		for skillName := range agent.Skills {
			skillsList = append(skillsList, types.StringValue(skillName))
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, skillsList)
		resp.Diagnostics.Append(diags...)
		config.Skills = listVal
	} else {
		config.Skills = types.ListNull(types.StringType)
	}

	// Convert execution environment to JSON string
	if len(agent.ExecutionEnvironment) > 0 {
		execEnvJSON, err := toJSONString(agent.ExecutionEnvironment)
		if err != nil {
			resp.Diagnostics.AddError("Error converting execution_environment", err.Error())
			return
		}
		config.ExecutionEnvironment = types.StringValue(execEnvJSON)
	} else {
		config.ExecutionEnvironment = types.StringNull()
	}

	if agent.CreatedAt != nil {
		config.CreatedAt = types.StringValue(agent.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if agent.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(agent.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
