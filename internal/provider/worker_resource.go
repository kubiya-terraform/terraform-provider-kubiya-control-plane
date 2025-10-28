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

var _ resource.Resource = (*workerResource)(nil)
var _ resource.ResourceWithImportState = (*workerResource)(nil)

func NewWorkerResource() resource.Resource {
	return &workerResource{}
}

type workerResource struct {
	client *clients.Client
}

type workerResourceModel struct {
	ID              types.String `tfsdk:"id"`
	EnvironmentName types.String `tfsdk:"environment_name"`
	Hostname        types.String `tfsdk:"hostname"`
	Status          types.String `tfsdk:"status"`
	Metadata        types.String `tfsdk:"metadata"`
	RegisteredAt    types.String `tfsdk:"registered_at"`
	LastHeartbeat   types.String `tfsdk:"last_heartbeat"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func (r *workerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_worker"
}

func (r *workerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Worker in the Control Plane.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Worker ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_name": schema.StringAttribute{
				Description: "Environment name",
				Required:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "Worker hostname",
				Optional:    true,
			},
			"status": schema.StringAttribute{
				Description: "Worker status",
				Computed:    true,
			},
			"metadata": schema.StringAttribute{
				Description: "Worker metadata as JSON string",
				Optional:    true,
			},
			"registered_at": schema.StringAttribute{
				Description: "Registration timestamp",
				Computed:    true,
			},
			"last_heartbeat": schema.StringAttribute{
				Description: "Last heartbeat timestamp",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp",
				Computed:    true,
			},
		},
	}
}

func (r *workerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan workerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &entities.WorkerCreateRequest{
		EnvironmentName: plan.EnvironmentName.ValueString(),
	}

	if !plan.Hostname.IsNull() {
		hostname := plan.Hostname.ValueString()
		createReq.Hostname = &hostname
	}

	if !plan.Metadata.IsNull() {
		metadata, err := parseJSON(plan.Metadata.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid Metadata", fmt.Sprintf("Failed to parse metadata JSON: %s", err))
			return
		}
		createReq.WorkerMetadata = metadata
	}

	worker, err := r.client.CreateWorker(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating worker", err.Error())
		return
	}

	plan.ID = types.StringValue(worker.ID)
	plan.EnvironmentName = types.StringValue(worker.EnvironmentName)
	plan.Status = types.StringValue(string(worker.Status))

	if worker.RegisteredAt != nil {
		plan.RegisteredAt = types.StringValue(worker.RegisteredAt.String())
	}

	if worker.LastHeartbeat != nil {
		plan.LastHeartbeat = types.StringValue(worker.LastHeartbeat.String())
	}

	if worker.UpdatedAt != nil {
		plan.UpdatedAt = types.StringValue(worker.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state workerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	worker, err := r.client.GetWorker(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading worker", err.Error())
		return
	}

	state.ID = types.StringValue(worker.ID)
	state.EnvironmentName = types.StringValue(worker.EnvironmentName)
	state.Status = types.StringValue(string(worker.Status))

	if worker.LastHeartbeat != nil {
		state.LastHeartbeat = types.StringValue(worker.LastHeartbeat.String())
	}

	if worker.UpdatedAt != nil {
		state.UpdatedAt = types.StringValue(worker.UpdatedAt.String())
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *workerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Workers are typically immutable after creation
	var plan workerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *workerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state workerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorker(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting worker", err.Error())
		return
	}
}

func (r *workerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
