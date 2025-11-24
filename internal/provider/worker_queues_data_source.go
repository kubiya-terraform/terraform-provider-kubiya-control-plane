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

var _ datasource.DataSource = (*workerQueuesDataSource)(nil)

func NewWorkerQueuesDataSource() datasource.DataSource {
	return &workerQueuesDataSource{}
}

type workerQueuesDataSource struct {
	client *clients.Client
}

type workerQueuesDataSourceModel struct {
	EnvironmentID types.String                 `tfsdk:"environment_id"`
	Queues        []workerQueueDataSourceModel `tfsdk:"queues"`
}

func (d *workerQueuesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_worker_queues"
}

func (d *workerQueuesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all Worker Queues in an environment.",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Description: "Environment ID to list worker queues from",
				Required:    true,
			},
			"queues": schema.ListNestedAttribute{
				Description: "List of worker queues",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Worker Queue ID",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *workerQueuesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *workerQueuesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data workerQueuesDataSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queues, err := d.client.ListWorkerQueues(data.EnvironmentID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error listing worker queues", err.Error())
		return
	}

	data.Queues = make([]workerQueueDataSourceModel, 0, len(queues))
	for _, queue := range queues {
		queueModel := workerQueueDataSourceModel{
			ID:                types.StringValue(queue.ID),
			EnvironmentID:     types.StringValue(queue.EnvironmentID),
			Name:              types.StringValue(queue.Name),
			HeartbeatInterval: types.Int64Value(int64(queue.HeartbeatInterval)),
			ActiveWorkers:     types.Int64Value(int64(queue.ActiveWorkers)),
			TaskQueueName:     types.StringValue(queue.TaskQueueName),
		}

		if queue.DisplayName != nil {
			queueModel.DisplayName = types.StringValue(*queue.DisplayName)
		} else {
			queueModel.DisplayName = types.StringNull()
		}

		if queue.Description != nil {
			queueModel.Description = types.StringValue(*queue.Description)
		} else {
			queueModel.Description = types.StringNull()
		}

		if queue.Status != "" {
			queueModel.Status = types.StringValue(string(queue.Status))
		} else {
			queueModel.Status = types.StringNull()
		}

		if queue.MaxWorkers != nil {
			queueModel.MaxWorkers = types.Int64Value(int64(*queue.MaxWorkers))
		} else {
			queueModel.MaxWorkers = types.Int64Null()
		}

		// Convert tags
		if queue.Tags != nil && len(queue.Tags) > 0 {
			tagValues := make([]attr.Value, len(queue.Tags))
			for i, tag := range queue.Tags {
				tagValues[i] = types.StringValue(tag)
			}
			queueModel.Tags = types.ListValueMust(types.StringType, tagValues)
		} else {
			queueModel.Tags = types.ListValueMust(types.StringType, []attr.Value{})
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
			queueModel.Settings = types.MapValueMust(types.StringType, settingsMap)
		} else {
			queueModel.Settings = types.MapValueMust(types.StringType, map[string]attr.Value{})
		}

		if queue.CreatedAt != nil {
			queueModel.CreatedAt = types.StringValue(queue.CreatedAt.String())
		} else {
			queueModel.CreatedAt = types.StringNull()
		}

		if queue.UpdatedAt != nil {
			queueModel.UpdatedAt = types.StringValue(queue.UpdatedAt.String())
		} else {
			queueModel.UpdatedAt = types.StringNull()
		}

		data.Queues = append(data.Queues, queueModel)
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
