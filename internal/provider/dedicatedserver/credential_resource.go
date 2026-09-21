package dedicatedserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver/v2"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.Resource              = &credentialResource{}
	_ resource.ResourceWithConfigure = &credentialResource{}
)

type credentialResource struct {
	utils.ResourceAPI
}

type credentialResourceModel struct {
	DedicatedServerID types.String `tfsdk:"dedicated_server_id"`
	Username          types.String `tfsdk:"username"`
	Type              types.String `tfsdk:"type"`
	Password          types.String `tfsdk:"password"`
}

func NewCredentialResource() resource.Resource {
	return &credentialResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "dedicated_server_credential",
		},
	}
}

func (c *credentialResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"dedicated_server_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the dedicated server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: `The username for the credentials`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: `The type of the credential. Valid options are: "OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"`,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"OPERATING_SYSTEM", "CONTROL_PANEL", "REMOTE_MANAGEMENT", "RESCUE_MODE", "SWITCH", "PDU", "FIREWALL", "LOAD_BALANCER"}...),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: `The password for the credentials`,
			},
		},
	}
}

func (c *credentialResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewCreateCredentialOpts(
		plan.Password.ValueString(),
		dedicatedserver.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
	)
	request := c.DedicatedserverAPI.CreateCredential(
		ctx,
		plan.DedicatedServerID.ValueString(),
	).CreateCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerID: plan.DedicatedServerID,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.DedicatedserverAPI.GetCredential(
		ctx,
		state.DedicatedServerID.ValueString(),
		dedicatedserver.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerID: state.DedicatedServerID,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan credentialResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := dedicatedserver.NewUpdateCredentialOpts(
		plan.Password.ValueString(),
	)
	request := c.DedicatedserverAPI.UpdateCredential(
		ctx,
		plan.DedicatedServerID.ValueString(),
		dedicatedserver.CredentialType(plan.Type.ValueString()),
		plan.Username.ValueString(),
	).UpdateCredentialOpts(*opts)
	result, response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(
			ctx,
			credentialResourceModel{
				DedicatedServerID: plan.DedicatedServerID,
				Type:              types.StringValue(string(result.GetType())),
				Password:          types.StringValue(result.GetPassword()),
				Username:          types.StringValue(result.GetUsername()),
			},
		)...,
	)
}

func (c *credentialResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state credentialResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := c.DedicatedserverAPI.DeleteCredential(
		ctx,
		state.DedicatedServerID.ValueString(),
		dedicatedserver.CredentialType(state.Type.ValueString()),
		state.Username.ValueString(),
	)
	response, err := request.Execute()
	if err != nil {
		utils.SdkError(ctx, &resp.Diagnostics, err, response)
	}
}
