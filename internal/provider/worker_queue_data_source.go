package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
)

var _ datasource.DataSource = (*workerQueueDataSource)(nil)

func NewWorkerQueueDataSource() datasource.DataSource {
	return &workerQueueDataSource{}
}

type workerQueueDataSource struct {
	client *clients.Client
}

type workerQueueDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	EnvironmentID     types.String `tfsdk:"environment_id"`
	Name              types.String `tfsdk:"name"`
	DisplayName       types.String `tfsdk:"display_name"`
	Description       types.String `tfsdk:"description"`
	Status            types.String `tfsdk:"status"`
	MaxWorkers        types.Int64  `tfsdk:"max_workers"`
	HeartbeatInterval types.Int64  `tfsdk:"heartbeat_interval"`
	Tags              types.List   `tfsdk:"tags"`
	Settings          types.Map    `tfsdk:"settings"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
	ActiveWorkers     types.Int64  `tfsdk:"active_workers"`
	TaskQueueName     types.String `tfsdk:"task_queue_name"`
}

func (d *workerQueueDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_worker_queue"
}

func (d *workerQueueDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a Worker Queue by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Worker Queue ID",
				Required:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Environment ID",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Worker queue name",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "User-friendly display name",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Queue description",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Worker queue status",
				Computed:    true,
			},
			"max_workers": schema.Int64Attribute{
				Description: "Maximum number of workers allowed",
				Computed:    true,
			},
			"heartbeat_interval": schema.Int64Attribute{
				Description: "Seconds between heartbeats",
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "Tags for the worker queue",
				ElementType: types.StringType,
				Computed:    true,
			},
			"settings": schema.MapAttribute{
				Description: "Additional settings as key-value pairs",
				ElementType: types.StringType,
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp",
				Computed:    true,
			},
			"active_workers": schema.Int64Attribute{
				Description: "Number of active workers",
				Computed:    true,
			},
			"task_queue_name": schema.StringAttribute{
				Description: "Task queue name for Temporal",
				Computed:    true,
			},
		},
	}
}

func (d *workerQueueDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *workerQueueDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workerQueueDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue, err := d.client.GetWorkerQueue(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading worker queue", err.Error())
		return
	}

	data.ID = types.StringValue(queue.ID)
	data.EnvironmentID = types.StringValue(queue.EnvironmentID)
	data.Name = types.StringValue(queue.Name)

	if queue.DisplayName != nil {
		data.DisplayName = types.StringValue(*queue.DisplayName)
	} else {
		data.DisplayName = types.StringNull()
	}

	if queue.Description != nil {
		data.Description = types.StringValue(*queue.Description)
	} else {
		data.Description = types.StringNull()
	}

	if queue.Status != "" {
		data.Status = types.StringValue(string(queue.Status))
	} else {
		data.Status = types.StringNull()
	}

	if queue.MaxWorkers != nil {
		data.MaxWorkers = types.Int64Value(int64(*queue.MaxWorkers))
	} else {
		data.MaxWorkers = types.Int64Null()
	}

	data.HeartbeatInterval = types.Int64Value(int64(queue.HeartbeatInterval))

	// Convert tags
	if queue.Tags != nil && len(queue.Tags) > 0 {
		tagValues := make([]attr.Value, len(queue.Tags))
		for i, tag := range queue.Tags {
			tagValues[i] = types.StringValue(tag)
		}
		data.Tags = types.ListValueMust(types.StringType, tagValues)
	} else {
		data.Tags = types.ListValueMust(types.StringType, []attr.Value{})
	}

	// Convert settings
	if queue.Settings != nil && len(queue.Settings) > 0 {
		settingsMap := make(map[string]attr.Value)
		for k, v := range queue.Settings {
			// Convert value to JSON string if it's not already a string
			var strValue string
			if s, ok := v.(string); ok {
				strValue = s
			} else {
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					resp.Diagnostics.AddError("Error marshaling settings value", err.Error())
					return
				}
				strValue = string(jsonBytes)
			}
			settingsMap[k] = types.StringValue(strValue)
		}
		data.Settings = types.MapValueMust(types.StringType, settingsMap)
	} else {
		data.Settings = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if queue.CreatedAt != nil {
		data.CreatedAt = types.StringValue(queue.CreatedAt.String())
	} else {
		data.CreatedAt = types.StringNull()
	}

	if queue.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(queue.UpdatedAt.String())
	} else {
		data.UpdatedAt = types.StringNull()
	}

	data.ActiveWorkers = types.Int64Value(int64(queue.ActiveWorkers))
	data.TaskQueueName = types.StringValue(queue.TaskQueueName)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
