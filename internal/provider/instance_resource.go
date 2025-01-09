package provider

import (
	"context"
	"fmt"
	"terraform-provider-elice/internal/api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource              = &instanceResource{}
	_ resource.ResourceWithConfigure = &instanceResource{}
)

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client *api.Client
}

type instanceResourceModel struct {
	Id             types.String `tfsdk:"id"`
	Title          types.String `tfsdk:"title"`
	ImageId        types.String `tfsdk:"image_id"`
	InstanceTypeId types.String `tfsdk:"instance_type_id"`
	Disk           types.Int64  `tfsdk:"disk"`
}

func (r *instanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *instanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *instanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"image_id": schema.StringAttribute{
				Required: true,
			},
			"instance_type_id": schema.StringAttribute{
				Required: true,
			},
			"disk": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan instanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.CreateInstance(
		plan.Title.ValueString(),
		plan.ImageId.ValueString(),
		plan.InstanceTypeId.ValueString(),
		int(plan.Disk.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating instance",
			"Could not create instance, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(instance.Id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.client.GetInstance(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading EliceCloud Instance",
			"Could not read EliceCloud instance ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}
	state.Id = types.StringValue(instance.Id)
	state.Title = types.StringValue(instance.Title)
	state.ImageId = types.StringValue(instance.Image.Id)
	state.InstanceTypeId = types.StringValue(instance.InstanceType.Id)
	state.Disk = types.Int64Value(int64(instance.Disk))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan instanceResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.UpdateInstance(
		state.Id.ValueString(),
		plan.Title.ValueString(),
		int(plan.Disk.ValueInt64()),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating instance",
			"Could not update instance, unexpected error: "+err.Error(),
		)
		return
	}

	instance, err := r.client.GetInstance(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading EliceCloud Instance",
			"Could not read EliceCloud instance ID "+state.Id.ValueString()+": "+err.Error(),
		)
		return
	}

	plan.Id = types.StringValue(instance.Id)
	plan.Title = types.StringValue(instance.Title)
	plan.ImageId = types.StringValue(instance.Image.Id)
	plan.InstanceTypeId = types.StringValue(instance.InstanceType.Id)
	plan.Disk = types.Int64Value(int64(instance.Disk))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state instanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.client.DeleteInstance(state.Id.ValueString())

	resp.State.RemoveResource(ctx)
}
