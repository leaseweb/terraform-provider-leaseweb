package ipmgmt

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &nullRouteResource{}
	_ resource.ResourceWithImportState = &nullRouteResource{}
)

type nullRouteResourceModel struct {
	AssignedContract     types.Object `tfsdk:"assigned_contract"`
	AutomaticUnnullingAt types.String `tfsdk:"automatic_unnulling_at"`
	Comment              types.String `tfsdk:"comment"`
	EquipmentID          types.String `tfsdk:"equipment_id"`
	ID                   types.String `tfsdk:"id"`
	IP                   types.String `tfsdk:"ip"`
	NulledAt             types.String `tfsdk:"nulled_at"`
	NulledBy             types.String `tfsdk:"nulled_by"`
	NullLevel            types.Int32  `tfsdk:"null_level"`
	TicketID             types.String `tfsdk:"ticket_id"`
	UnnulledAt           types.String `tfsdk:"unnulled_at"`
	UnnulledBy           types.String `tfsdk:"unnulled_by"`
}

func adaptNullRouteToResourceModel(
	nullRoutedIP ipmgmt.NullRoutedIP,
	diags *diag.Diagnostics,
	ctx context.Context,
) *nullRouteResourceModel {
	sdkAssignedContract, _ := nullRoutedIP.GetAssignedContractOk()
	assignedContract := utils.AdaptNullableSdkModelToResourceObject(
		sdkAssignedContract,
		map[string]attr.Type{"id": types.StringType},
		ctx,
		func(contract ipmgmt.AssignedContract) assignedContractResourceModel {
			return assignedContractResourceModel{
				ID: basetypes.NewStringValue(contract.GetId()),
			}
		},
		diags,
	)
	if diags.HasError() {
		return nil
	}

	automaticUnnullingAt, _ := nullRoutedIP.GetAutomatedUnnullingAtOk()
	comment, _ := nullRoutedIP.GetCommentOk()
	ticketID, _ := nullRoutedIP.GetTicketIdOk()
	unnulledAt, _ := nullRoutedIP.GetUnnulledAtOk()
	unnulledBy, _ := nullRoutedIP.GetUnnulledByOk()

	return &nullRouteResourceModel{
		AssignedContract:     assignedContract,
		AutomaticUnnullingAt: utils.AdaptNullableTimeToStringValue(automaticUnnullingAt),
		Comment:              basetypes.NewStringPointerValue(comment),
		EquipmentID:          basetypes.NewStringValue(nullRoutedIP.GetEquipmentId()),
		ID:                   basetypes.NewStringValue(nullRoutedIP.GetId()),
		IP:                   basetypes.NewStringValue(nullRoutedIP.GetIp()),
		NulledAt:             basetypes.NewStringValue(nullRoutedIP.GetNulledAt().String()),
		NulledBy:             basetypes.NewStringValue(nullRoutedIP.GetNulledBy()),
		NullLevel:            basetypes.NewInt32Value(nullRoutedIP.GetNullLevel()),
		TicketID:             basetypes.NewStringPointerValue(ticketID),
		UnnulledAt:           utils.AdaptNullableTimeToStringValue(unnulledAt),
		UnnulledBy:           basetypes.NewStringPointerValue(unnulledBy),
	}
}

type assignedContractResourceModel struct {
	ID types.String `tfsdk:"id"`
}

type nullRouteResource struct {
	utils.ResourceAPI
}

func (n nullRouteResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	resource.ImportStatePassthroughID(
		ctx,
		path.Root("id"),
		request,
		response,
	)
}

func (n nullRouteResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Null route ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"assigned_contract": schema.SingleNestedAttribute{
				Description: "An explanation about the purpose of this instance",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the contract connected to the IP",
					},
				},
			},
			"automatic_unnulling_at": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The date and time when the null route is to be deactivated. The date and time should be specified using the `2019-09-08 00:00:00 +0000 UTC` format. If this field is not present then the null route will not be automatically removed",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} UTC$`), "must be specified using the RFC3339 format (`yyyy-mm-ddThh:mm:ssZ`)"),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A comment to be stored with the null route (e.g. null route reason)",
			},
			"equipment_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the equipment which was assigned to the IP at the time of null route creation",
			},
			"nulled_at": schema.StringAttribute{
				Computed:    true,
				Description: "Null route date",
			},
			"nulled_by": schema.StringAttribute{
				Computed:    true,
				Description: "Email address of the user who created the null route or 'LeaseWeb' if null route was created by LeaseWeb",
			},
			"null_level": schema.Int32Attribute{
				Computed:    true,
				Description: "Null route permission level. If greater than 1 then the null route can only be removed by LeaseWeb",
			},
			"ticket_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A reference to be stored with the null route",
			},
			"unnulled_at": schema.StringAttribute{
				Computed:    true,
				Description: "The date and time when the null route has been removed. If null then the null route is still active",
			},
			"unnulled_by": schema.StringAttribute{
				Computed:    true,
				Description: "Email address of the user who removed the null route or 'LeaseWeb' if null route was removed by LeaseWeb",
			},
		},
	}
}

func (n nullRouteResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan nullRouteResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	if utils.AdaptStringPointerValueToNullableString(plan.IP) == nil {
		response.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"Attribute not set",
			"Attribute ip value must be set to create a null route",
		)
		return
	}

	opts := ipmgmt.NewNullRouteIPOpts()
	if !plan.AutomaticUnnullingAt.IsNull() && !plan.AutomaticUnnullingAt.IsUnknown() {
		automatedUnnullingAt, err := time.Parse(
			"2006-01-02 15:04:05 -0700 MST",
			plan.AutomaticUnnullingAt.ValueString(),
		)
		if err != nil {
			utils.GeneralError(&response.Diagnostics, ctx, err)
			return
		}
		opts.SetAutomatedUnnullingAt(automatedUnnullingAt)
	}

	opts.Comment = utils.AdaptStringPointerValueToNullableString(plan.Comment)
	opts.TicketId = utils.AdaptStringPointerValueToNullableString(plan.TicketID)

	nullRoutedIP, httpResponse, err := n.IPmgmtAPI.NullRouteIP(
		ctx,
		plan.IP.ValueString(),
	).NullRouteIPOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptNullRouteToResourceModel(
		*nullRoutedIP,
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (n nullRouteResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState nullRouteResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	if originalState.ID.IsNull() {
		response.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Attribute not set",
			"Attribute id value must be set to inspect a null route",
		)
		return
	}

	nullRoutedIP, httpResponse, err := n.IPmgmtAPI.InspectNullRouteHistory(
		ctx,
		originalState.ID.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptNullRouteToResourceModel(
		*nullRoutedIP,
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (n nullRouteResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan nullRouteResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	if plan.ID.IsNull() || plan.ID.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Attribute not set",
			"Attribute id value must be set to update a null route",
		)
		return
	}

	opts := ipmgmt.NewUpdateNullRouteOpts()
	if !plan.AutomaticUnnullingAt.IsNull() && !plan.AutomaticUnnullingAt.IsUnknown() {
		automatedUnnullingAt, err := time.Parse(
			"2006-01-02 15:04:05 -0700 MST",
			plan.AutomaticUnnullingAt.ValueString(),
		)
		if err != nil {
			utils.GeneralError(&response.Diagnostics, ctx, err)
			return
		}
		opts.SetAutomatedUnnullingAt(automatedUnnullingAt)
	}
	if !plan.Comment.IsNull() && !plan.Comment.IsUnknown() {
		opts.SetComment(plan.Comment.ValueString())
	}
	if !plan.TicketID.IsNull() && !plan.TicketID.IsUnknown() {
		opts.SetTicketId(plan.TicketID.ValueString())
	}
	nullRoutedIP, httpResponse, err := n.IPmgmtAPI.UpdateNullRoute(
		ctx,
		plan.ID.ValueString(),
	).UpdateNullRouteOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptNullRouteToResourceModel(
		*nullRoutedIP,
		&response.Diagnostics,
		ctx,
	)
	if response.Diagnostics.HasError() {
		return
	}

	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (n nullRouteResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state nullRouteResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	if state.IP.IsNull() || state.IP.IsUnknown() {
		response.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"Attribute not set",
			"Attribute ip value must be set to delete a null route",
		)
		return
	}

	httpResponse, err := n.IPmgmtAPI.RemoveIPNullRoute(
		ctx,
		state.IP.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewNullRouteResource() resource.Resource {
	return &nullRouteResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "ipmgmt_null_route",
		},
	}
}
