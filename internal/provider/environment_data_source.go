package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*environmentDataSource)(nil)

func NewEnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

type environmentDataSource struct {
	client *clients.Client
}

type environmentDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	DisplayName          types.String `tfsdk:"display_name"`
	Description          types.String `tfsdk:"description"`
	Tags                 types.List   `tfsdk:"tags"`
	Settings             types.String `tfsdk:"settings"`
	Status               types.String `tfsdk:"status"`
	ExecutionEnvironment types.String `tfsdk:"execution_environment"`
	WorkerToken          types.String `tfsdk:"worker_token"`
	TemporalNamespaceID  types.String `tfsdk:"temporal_namespace_id"`
	ActiveWorkers        types.Int64  `tfsdk:"active_workers"`
	IdleWorkers          types.Int64  `tfsdk:"idle_workers"`
	BusyWorkers          types.Int64  `tfsdk:"busy_workers"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches an existing Environment from the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Environment ID to lookup",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Environment name",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "User-friendly display name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Environment description",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags for categorization",
				Computed:    true,
				ElementType: types.StringType,
			},
			"settings": schema.StringAttribute{
				Description: "Environment settings as JSON string",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Environment status",
				Computed:    true,
			},
			"execution_environment": schema.StringAttribute{
				Description: "Execution environment configuration as JSON string",
				Computed:    true,
			},
			"worker_token": schema.StringAttribute{
				Description: "Worker registration token",
				Computed:    true,
				Sensitive:   true,
			},
			"temporal_namespace_id": schema.StringAttribute{
				Description: "Temporal namespace ID",
				Computed:    true,
			},
			"active_workers": schema.Int64Attribute{
				Description: "Number of active workers",
				Computed:    true,
			},
			"idle_workers": schema.Int64Attribute{
				Description: "Number of idle workers",
				Computed:    true,
			},
			"busy_workers": schema.Int64Attribute{
				Description: "Number of busy workers",
				Computed:    true,
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

func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config environmentDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment, err := d.client.GetEnvironment(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading environment", err.Error())
		return
	}

	// Map response to state
	config.Name = types.StringValue(environment.Name)

	if environment.DisplayName != nil {
		config.DisplayName = types.StringValue(*environment.DisplayName)
	} else {
		config.DisplayName = types.StringNull()
	}

	if environment.Description != nil {
		config.Description = types.StringValue(*environment.Description)
	} else {
		config.Description = types.StringNull()
	}

	// Convert tags to list
	if len(environment.Tags) > 0 {
		tagList := make([]types.String, len(environment.Tags))
		for i, tag := range environment.Tags {
			tagList[i] = types.StringValue(tag)
		}
		listVal, diags := types.ListValueFrom(ctx, types.StringType, tagList)
		resp.Diagnostics.Append(diags...)
		config.Tags = listVal
	} else {
		config.Tags = types.ListNull(types.StringType)
	}

	// Convert settings to JSON string
	if len(environment.Settings) > 0 {
		settingsJSON, err := toJSONString(environment.Settings)
		if err != nil {
			resp.Diagnostics.AddError("Error converting settings", err.Error())
			return
		}
		config.Settings = types.StringValue(settingsJSON)
	} else {
		config.Settings = types.StringNull()
	}

	config.Status = types.StringValue(string(environment.Status))

	// Convert execution environment to JSON string
	if len(environment.ExecutionEnvironment) > 0 {
		execEnvJSON, err := toJSONString(environment.ExecutionEnvironment)
		if err != nil {
			resp.Diagnostics.AddError("Error converting execution_environment", err.Error())
			return
		}
		config.ExecutionEnvironment = types.StringValue(execEnvJSON)
	} else {
		config.ExecutionEnvironment = types.StringNull()
	}

	if environment.WorkerToken != nil {
		config.WorkerToken = types.StringValue(*environment.WorkerToken)
	} else {
		config.WorkerToken = types.StringNull()
	}

	if environment.TemporalNamespaceID != nil {
		config.TemporalNamespaceID = types.StringValue(*environment.TemporalNamespaceID)
	} else {
		config.TemporalNamespaceID = types.StringNull()
	}

	config.ActiveWorkers = types.Int64Value(int64(environment.ActiveWorkers))
	config.IdleWorkers = types.Int64Value(int64(environment.IdleWorkers))
	config.BusyWorkers = types.Int64Value(int64(environment.BusyWorkers))

	if environment.CreatedAt != nil {
		config.CreatedAt = types.StringValue(environment.CreatedAt.String())
	} else {
		config.CreatedAt = types.StringNull()
	}

	if environment.UpdatedAt != nil {
		config.UpdatedAt = types.StringValue(environment.UpdatedAt.String())
	} else {
		config.UpdatedAt = types.StringNull()
	}

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}
