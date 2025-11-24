package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-kubiya-control-plane/internal/clients"
	"terraform-provider-kubiya-control-plane/internal/entities"
)

var _ resource.Resource = (*workerQueueResource)(nil)
var _ resource.ResourceWithImportState = (*workerQueueResource)(nil)

func NewWorkerQueueResource() resource.Resource {
	return &workerQueueResource{}
}

type workerQueueResource struct {
	client *clients.Client
}

type workerQueueResourceModel struct {
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

func (r *workerQueueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_worker_queue"
}

func (r *workerQueueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Worker Queue in the Control Plane. Worker queues are used to organize and manage workers within an environment.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Worker Queue ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				Description: "Environment ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Worker queue name (lowercase, no spaces, 2-50 characters)",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "User-friendly display name",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"description": schema.StringAttribute{
				Description: "Queue description",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"status": schema.StringAttribute{
				Description: "Worker queue status (active, inactive, paused)",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("active"),
			},
			"max_workers": schema.Int64Attribute{
				Description: "Maximum number of workers allowed (null = unlimited)",
				Optional:    true,
			},
			"heartbeat_interval": schema.Int64Attribute{
				Description: "Seconds between heartbeats (lightweight) (10-300)",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60),
			},
			"tags": schema.ListAttribute{
				Description: "Tags for the worker queue",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, []attr.Value{})),
			},
			"settings": schema.MapAttribute{
				Description: "Additional settings as key-value pairs",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
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
				Description: "Number of active workers (computed)",
				Computed:    true,
			},
			"task_queue_name": schema.StringAttribute{
				Description: "Task queue name for Temporal (computed)",
				Computed:    true,
			},
		},
	}
}

func (r *workerQueueResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workerQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workerQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.WorkerQueueCreateRequest{
		Name:              plan.Name.ValueString(),
		HeartbeatInterval: int(plan.HeartbeatInterval.ValueInt64()),
	}

	if !plan.DisplayName.IsNull() && plan.DisplayName.ValueString() != "" {
		displayName := plan.DisplayName.ValueString()
		createReq.DisplayName = &displayName
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		description := plan.Description.ValueString()
		createReq.Description = &description
	}

	if !plan.MaxWorkers.IsNull() {
		maxWorkers := int(plan.MaxWorkers.ValueInt64())
		createReq.MaxWorkers = &maxWorkers
	}

	// Convert tags
	if !plan.Tags.IsNull() {
		var tags []string
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Tags = tags
	}

	// Convert settings
	if !plan.Settings.IsNull() {
		var settings map[string]string
		diags = plan.Settings.ElementsAs(ctx, &settings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert to map[string]interface{}
		settingsInterface := make(map[string]interface{})
		for k, v := range settings {
			settingsInterface[k] = v
		}
		createReq.Settings = settingsInterface
	}

	queue, err := r.client.CreateWorkerQueue(plan.EnvironmentID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating worker queue", err.Error())
		return
	}

	// Update state
	plan.ID = types.StringValue(queue.ID)
	plan.Name = types.StringValue(queue.Name)

	if queue.DisplayName != nil {
		plan.DisplayName = types.StringValue(*queue.DisplayName)
	} else {
		plan.DisplayName = types.StringNull()
	}

	if queue.Description != nil {
		plan.Description = types.StringValue(*queue.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if queue.Status != "" {
		plan.Status = types.StringValue(string(queue.Status))
	} else {
		plan.Status = types.StringValue("active")
	}

	if queue.MaxWorkers != nil {
		plan.MaxWorkers = types.Int64Value(int64(*queue.MaxWorkers))
	} else {
		plan.MaxWorkers = types.Int64Null()
	}

	plan.HeartbeatInterval = types.Int64Value(int64(queue.HeartbeatInterval))

	// Convert tags
	if queue.Tags != nil && len(queue.Tags) > 0 {
		tagValues := make([]attr.Value, len(queue.Tags))
		for i, tag := range queue.Tags {
			tagValues[i] = types.StringValue(tag)
		}
		plan.Tags = types.ListValueMust(types.StringType, tagValues)
	} else {
		plan.Tags = types.ListValueMust(types.StringType, []attr.Value{})
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
		plan.Settings = types.MapValueMust(types.StringType, settingsMap)
	} else {
		plan.Settings = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if queue.CreatedAt != nil {
		plan.CreatedAt = types.StringValue(queue.CreatedAt.String())
	} else {
		plan.CreatedAt = types.StringNull()
	}

	if queue.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(queue.UpdatedAt.String())
	} else {
		plan.UpdatedAt = types.StringNull()
	}

	plan.ActiveWorkers = types.Int64Value(int64(queue.ActiveWorkers))
	plan.TaskQueueName = types.StringValue(queue.TaskQueueName)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workerQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workerQueueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	queue, err := r.client.GetWorkerQueue(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading worker queue", err.Error())
		return
	}

	state.ID = types.StringValue(queue.ID)
	state.EnvironmentID = types.StringValue(queue.EnvironmentID)
	state.Name = types.StringValue(queue.Name)

	if queue.DisplayName != nil {
		state.DisplayName = types.StringValue(*queue.DisplayName)
	} else {
		state.DisplayName = types.StringNull()
	}

	if queue.Description != nil {
		state.Description = types.StringValue(*queue.Description)
	} else {
		state.Description = types.StringNull()
	}

	if queue.Status != "" {
		state.Status = types.StringValue(string(queue.Status))
	} else {
		state.Status = types.StringValue("active")
	}

	if queue.MaxWorkers != nil {
		state.MaxWorkers = types.Int64Value(int64(*queue.MaxWorkers))
	} else {
		state.MaxWorkers = types.Int64Null()
	}

	state.HeartbeatInterval = types.Int64Value(int64(queue.HeartbeatInterval))

	// Convert tags
	if queue.Tags != nil && len(queue.Tags) > 0 {
		tagValues := make([]attr.Value, len(queue.Tags))
		for i, tag := range queue.Tags {
			tagValues[i] = types.StringValue(tag)
		}
		state.Tags = types.ListValueMust(types.StringType, tagValues)
	} else {
		state.Tags = types.ListValueMust(types.StringType, []attr.Value{})
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
		state.Settings = types.MapValueMust(types.StringType, settingsMap)
	} else {
		state.Settings = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if queue.CreatedAt != nil {
		state.CreatedAt = types.StringValue(queue.CreatedAt.String())
	} else {
		state.CreatedAt = types.StringNull()
	}

	if queue.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(queue.UpdatedAt.String())
	} else {
		state.UpdatedAt = types.StringNull()
	}

	state.ActiveWorkers = types.Int64Value(int64(queue.ActiveWorkers))
	state.TaskQueueName = types.StringValue(queue.TaskQueueName)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *workerQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workerQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &entities.WorkerQueueUpdateRequest{}

	// Only update fields that are set
	if !plan.Name.IsNull() {
		name := plan.Name.ValueString()
		updateReq.Name = &name
	}

	if !plan.DisplayName.IsNull() && plan.DisplayName.ValueString() != "" {
		displayName := plan.DisplayName.ValueString()
		updateReq.DisplayName = &displayName
	}

	if !plan.Description.IsNull() && plan.Description.ValueString() != "" {
		description := plan.Description.ValueString()
		updateReq.Description = &description
	}

	if !plan.Status.IsNull() {
		status := plan.Status.ValueString()
		updateReq.Status = &status
	}

	if !plan.MaxWorkers.IsNull() {
		maxWorkers := int(plan.MaxWorkers.ValueInt64())
		updateReq.MaxWorkers = &maxWorkers
	}

	if !plan.HeartbeatInterval.IsNull() {
		heartbeatInterval := int(plan.HeartbeatInterval.ValueInt64())
		updateReq.HeartbeatInterval = &heartbeatInterval
	}

	// Convert tags
	if !plan.Tags.IsNull() {
		var tags []string
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.Tags = tags
	}

	// Convert settings
	if !plan.Settings.IsNull() {
		var settings map[string]string
		diags = plan.Settings.ElementsAs(ctx, &settings, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		// Convert to map[string]interface{}
		settingsInterface := make(map[string]interface{})
		for k, v := range settings {
			settingsInterface[k] = v
		}
		updateReq.Settings = settingsInterface
	}

	queue, err := r.client.UpdateWorkerQueue(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating worker queue", err.Error())
		return
	}

	// Update state with response
	plan.Name = types.StringValue(queue.Name)

	if queue.DisplayName != nil {
		plan.DisplayName = types.StringValue(*queue.DisplayName)
	} else {
		plan.DisplayName = types.StringNull()
	}

	if queue.Description != nil {
		plan.Description = types.StringValue(*queue.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if queue.Status != "" {
		plan.Status = types.StringValue(string(queue.Status))
	} else {
		plan.Status = types.StringValue("active")
	}

	if queue.MaxWorkers != nil {
		plan.MaxWorkers = types.Int64Value(int64(*queue.MaxWorkers))
	} else {
		plan.MaxWorkers = types.Int64Null()
	}

	plan.HeartbeatInterval = types.Int64Value(int64(queue.HeartbeatInterval))

	// Convert tags
	if queue.Tags != nil && len(queue.Tags) > 0 {
		tagValues := make([]attr.Value, len(queue.Tags))
		for i, tag := range queue.Tags {
			tagValues[i] = types.StringValue(tag)
		}
		plan.Tags = types.ListValueMust(types.StringType, tagValues)
	} else {
		plan.Tags = types.ListValueMust(types.StringType, []attr.Value{})
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
		plan.Settings = types.MapValueMust(types.StringType, settingsMap)
	} else {
		plan.Settings = types.MapValueMust(types.StringType, map[string]attr.Value{})
	}

	if queue.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(queue.UpdatedAt.String())
	} else {
		plan.UpdatedAt = types.StringNull()
	}

	plan.ActiveWorkers = types.Int64Value(int64(queue.ActiveWorkers))
	plan.TaskQueueName = types.StringValue(queue.TaskQueueName)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workerQueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workerQueueResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkerQueue(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting worker queue", err.Error())
		return
	}
}

func (r *workerQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
