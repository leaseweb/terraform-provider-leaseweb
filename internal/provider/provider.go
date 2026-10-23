package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dedicatedserver"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/dns"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/ipmgmt"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/publiccloud"
)

var (
	_ provider.Provider = &leasewebProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &leasewebProvider{
			version: version,
		}
	}
}

type leasewebProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type leasewebProviderModel struct {
	Host   types.String `tfsdk:"host"`
	Token  types.String `tfsdk:"token"`
	Scheme types.String `tfsdk:"scheme"`
}

func (p *leasewebProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "leaseweb"
	resp.Version = p.version
}

func (p *leasewebProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "Host for Leaseweb API, defaults to \"api.leaseweb.com\". May also be provided via LEASEWEB_HOST environment variable if present.",
			},
			"scheme": schema.StringAttribute{
				Optional:    true,
				Description: "Scheme for Leaseweb API, defaults to \"https\". May also be provided via LEASEWEB_SCHEME environment variable if present.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "The API token to use. By default it takes the value from the LEASEWEB_TOKEN environment variable if present.",
				Sensitive:   true,
			},
		},
	}
}

func (p *leasewebProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var config leasewebProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown Leaseweb API token",
			"The provider cannot create the Leaseweb API client as there is an unknown configuration value for the Leaseweb API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the LEASEWEB_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("LEASEWEB_HOST")
	scheme := os.Getenv("LEASEWEB_SCHEME")
	token := os.Getenv("LEASEWEB_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Scheme.IsNull() {
		scheme = config.Scheme.ValueString()
	}

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Leaseweb API token",
			"The provider cannot create the Leaseweb API client as there is a missing or empty value for the Leaseweb API token. "+
				"Set the token value in the configuration or use the LEASEWEB_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "leaseweb_host", host)
	ctx = tflog.SetField(ctx, "leaseweb_scheme", scheme)
	ctx = tflog.SetField(ctx, "leaseweb_token", token)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "leaseweb_token")

	optional := client.Optional{}
	if host != "" {
		optional.Host = &host
	}
	if scheme != "" {
		optional.Scheme = &scheme
	}

	coreClient := client.NewClient(token, optional, p.version)

	resp.DataSourceData = coreClient
	resp.ResourceData = coreClient

	tflog.Info(
		ctx,
		"Configured Leaseweb client",
		map[string]any{"success": true},
	)
}

func (p *leasewebProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		publiccloud.NewInstancesDataSource,
		publiccloud.NewCredentialDataSource,
		dedicatedserver.NewServerDataSource,
		dedicatedserver.NewServersDataSource,
		dedicatedserver.NewControlPanelsDataSource,
		dedicatedserver.NewOperatingSystemsDataSource,
		dedicatedserver.NewCredentialDataSource,
		publiccloud.NewImagesDataSource,
		publiccloud.NewLoadBalancersDataSource,
		publiccloud.NewLoadBalancerListenersDataSource,
		publiccloud.NewTargetGroupsDataSource,
		publiccloud.NewISOsDataSource,
		dns.NewResourceRecordSetsDataSource,
		ipmgmt.NewIPsDataSource,
		ipmgmt.NewNullRouteHistoryDataSource,
	}
}

func (p *leasewebProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		publiccloud.NewInstanceResource,
		publiccloud.NewCredentialResource,
		dedicatedserver.NewServerResource,
		dedicatedserver.NewCredentialResource,
		dedicatedserver.NewNotificationSettingDatatrafficResource,
		dedicatedserver.NewNotificationSettingBandwidthResource,
		dedicatedserver.NewInstallationResource,
		publiccloud.NewImageResource,
		publiccloud.NewLoadBalancerResource,
		publiccloud.NewLoadBalancerListenerResource,
		publiccloud.NewTargetGroupResource,
		publiccloud.NewIPResource,
		publiccloud.NewInstanceIsoResource,
		dns.NewResourceRecordSetsResource,
		ipmgmt.NewIPResource,
		ipmgmt.NewNullRouteResource,
	}
}
