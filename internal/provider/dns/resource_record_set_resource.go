package dns

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/leaseweb/leaseweb-go-sdk/dns"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/utils"
)

var (
	_ resource.ResourceWithConfigure   = &resourceRecordSetResource{}
	_ resource.ResourceWithImportState = &resourceRecordSetResource{}
)

type resourceRecordSetResourceModel struct {
	Content    types.List   `tfsdk:"content"`
	DomainName types.String `tfsdk:"domain_name"`
	Name       types.String `tfsdk:"name"`
	TTL        types.Int32  `tfsdk:"ttl"`
	RecordType types.String `tfsdk:"type"`
}

func adaptResourceRecordSetDetailsToResourceRecordSetResourceResource(
	domainName string,
	resourceRecordSetDetails dns.ResourceRecordSetDetails,
	ctx context.Context,
	diagnostics *diag.Diagnostics,
) *resourceRecordSetResourceModel {
	content, diags := basetypes.NewListValueFrom(
		ctx,
		basetypes.StringType{},
		resourceRecordSetDetails.GetContent(),
	)
	if diags.HasError() {
		diagnostics.Append(diags...)
		return nil
	}

	return &resourceRecordSetResourceModel{
		DomainName: basetypes.NewStringValue(domainName),
		Content:    content,
		Name:       basetypes.NewStringValue(resourceRecordSetDetails.GetName()),
		TTL:        basetypes.NewInt32Value(int32(resourceRecordSetDetails.GetTtl())),
		RecordType: basetypes.NewStringValue(string(resourceRecordSetDetails.GetType())),
	}
}

type resourceRecordSetResource struct {
	utils.ResourceAPI
}

func (r *resourceRecordSetResource) ImportState(
	ctx context.Context,
	request resource.ImportStateRequest,
	response *resource.ImportStateResponse,
) {
	idParts := strings.Split(request.ID, ",")

	if len(idParts) != 3 || idParts[0] == "" || idParts[1] == "" || idParts[2] == "" {
		utils.UnexpectedImportIdentifierError(
			&response.Diagnostics,
			"domain_name,name,type",
			request.ID,
		)
		return
	}

	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("domain_name"),
		idParts[0],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("name"),
		idParts[1],
	)...)
	response.Diagnostics.Append(response.State.SetAttribute(
		ctx,
		path.Root("type"),
		idParts[2],
	)...)
}

func (r *resourceRecordSetResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	response *resource.SchemaResponse,
) {
	ttl := utils.NewIntMarkdownList(dns.AllowedTtlEnumValues)
	warningError := "**WARNING!** Changing this value once running will cause this record to be destroyed and a new one to be created."

	response.Schema = schema.Schema{
		Description: "Manage a DNS record",
		Attributes: map[string]schema.Attribute{
			"content": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Array of resource record set Content entries",
				Validators: []validator.List{
					listvalidator.NoNullValues(),
					listvalidator.SizeAtLeast(1),
				},
			},
			"domain_name": schema.StringAttribute{
				Required:    true,
				Description: "Domain Name",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the resource record set. " + warningError,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(regexp.MustCompile(`^.*\.$`), "must end in ."),
				},
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"ttl": schema.Int32Attribute{
				Required: true,
				Description: fmt.Sprintf(
					"Time to live of the resource record set. Valid options are %s",
					ttl.Markdown(),
				),
				Validators: []validator.Int32{
					int32validator.OneOf(ttl.ToInt32()...),
				},
			},
			"type": schema.StringAttribute{
				Required: true,
				Description: fmt.Sprintf(
					"Type of the resource record set. Valid options are %s",
					utils.StringTypeArrayToMarkdown(dns.AllowedResourceRecordSetTypeEnumValues),
				),
				Validators: []validator.String{
					stringvalidator.OneOf(utils.AdaptStringTypeArrayToStringArray(dns.AllowedResourceRecordSetTypeEnumValues)...),
				},
			},
		},
	}
}

func (r *resourceRecordSetResource) Create(
	ctx context.Context,
	request resource.CreateRequest,
	response *resource.CreateResponse,
) {
	var plan resourceRecordSetResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var content []string
	response.Diagnostics.Append(
		plan.Content.ElementsAs(ctx, &content, false)...,
	)
	if response.Diagnostics.HasError() {
		return
	}
	resourceRecordSetDetails, httpResponse, err := r.DNSAPI.CreateResourceRecordSet(
		ctx,
		plan.DomainName.ValueString(),
	).ResourceRecordSet(
		*dns.NewResourceRecordSet(
			plan.Name.ValueString(),
			dns.ResourceRecordSetType(plan.RecordType.ValueString()),
			content,
			dns.Ttl(plan.TTL.ValueInt32()),
		),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptResourceRecordSetDetailsToResourceRecordSetResourceResource(
		plan.DomainName.ValueString(),
		*resourceRecordSetDetails,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *resourceRecordSetResource) Read(
	ctx context.Context,
	request resource.ReadRequest,
	response *resource.ReadResponse,
) {
	var originalState resourceRecordSetResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &originalState)...)
	if response.Diagnostics.HasError() {
		return
	}

	resourceRecordSetDetails, httpResponse, err := r.DNSAPI.GetResourceRecordSet(
		ctx,
		originalState.DomainName.ValueString(),
		originalState.Name.ValueString(),
		originalState.RecordType.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptResourceRecordSetDetailsToResourceRecordSetResourceResource(
		originalState.DomainName.ValueString(),
		*resourceRecordSetDetails,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *resourceRecordSetResource) Update(
	ctx context.Context,
	request resource.UpdateRequest,
	response *resource.UpdateResponse,
) {
	var plan resourceRecordSetResourceModel
	response.Diagnostics.Append(request.Plan.Get(ctx, &plan)...)
	if response.Diagnostics.HasError() {
		return
	}

	var contents []string
	response.Diagnostics.Append(
		plan.Content.ElementsAs(ctx, &contents, false)...,
	)
	if response.Diagnostics.HasError() {
		return
	}
	opts := dns.NewUpdateResourceRecordSetOpts(
		contents,
		dns.Ttl(plan.TTL.ValueInt32()),
	)
	resourceRecordSetDetails, httpResponse, err := r.DNSAPI.UpdateResourceRecordSet(
		ctx,
		plan.DomainName.ValueString(),
		plan.Name.ValueString(),
		plan.RecordType.ValueString(),
	).UpdateResourceRecordSetOpts(*opts).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
		return
	}

	state := adaptResourceRecordSetDetailsToResourceRecordSetResourceResource(
		plan.DomainName.ValueString(),
		*resourceRecordSetDetails,
		ctx,
		&response.Diagnostics,
	)
	if response.Diagnostics.HasError() {
		return
	}
	response.Diagnostics.Append(response.State.Set(ctx, state)...)
}

func (r *resourceRecordSetResource) Delete(
	ctx context.Context,
	request resource.DeleteRequest,
	response *resource.DeleteResponse,
) {
	var state resourceRecordSetResourceModel
	response.Diagnostics.Append(request.State.Get(ctx, &state)...)
	if response.Diagnostics.HasError() {
		return
	}

	httpResponse, err := r.DNSAPI.DeleteResourceRecordSet(
		ctx,
		state.DomainName.ValueString(),
		state.Name.ValueString(),
		state.RecordType.ValueString(),
	).Execute()
	if err != nil {
		utils.SdkError(ctx, &response.Diagnostics, err, httpResponse)
	}
}

func NewResourceRecordSetsResource() resource.Resource {
	return &resourceRecordSetResource{
		ResourceAPI: utils.ResourceAPI{
			Name: "dns_resource_record_set",
		},
	}
}
