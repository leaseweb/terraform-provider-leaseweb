package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver/v2"
	"github.com/leaseweb/leaseweb-go-sdk/dns"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
	"github.com/leaseweb/terraform-provider-leaseweb/internal/provider/client"
)

func generateTypeName(providerTypeName string, name string) string {
	return fmt.Sprintf("%s_%s", providerTypeName, name)
}

func getCoreClient(
	providerData any,
	diagnostics *diag.Diagnostics,
) *client.Client {
	if providerData == nil {
		return nil
	}

	coreClient, ok := providerData.(client.Client)

	if !ok {
		diagnostics.AddError(
			"Unexpected Configure Type",
			fmt.Sprintf(
				"Expected client.Client, got: %T. Please report this issue to the provider developers.",
				providerData,
			),
		)
		return nil
	}

	return &coreClient
}

// ResourceAPI contains reusable Configure & Metadata functions for resources.
type ResourceAPI struct {
	Name               string
	PubliccloudAPI     publiccloud.PubliccloudAPI
	DedicatedserverAPI dedicatedserver.DedicatedserverAPI
	DNSAPI             dns.DnsAPI
	IPmgmtAPI          ipmgmt.IpmgmtAPI
}

func (p *ResourceAPI) Configure(
	_ context.Context,
	request resource.ConfigureRequest,
	response *resource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	p.PubliccloudAPI = coreClient.PubliccloudAPI
	p.DedicatedserverAPI = coreClient.DedicatedserverAPI
	p.DNSAPI = coreClient.DNSAPI
	p.IPmgmtAPI = coreClient.IPmgmtAPI
}

func (p *ResourceAPI) Metadata(
	_ context.Context,
	request resource.MetadataRequest,
	response *resource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, p.Name)
}

// DataSourceAPI contains reusable Configure & Metadata functions for data sources.
type DataSourceAPI struct {
	Name               string
	PubliccloudAPI     publiccloud.PubliccloudAPI
	DedicatedserverAPI dedicatedserver.DedicatedserverAPI
	DNSAPI             dns.DnsAPI
	IPmgmtAPI          ipmgmt.IpmgmtAPI
}

func (d *DataSourceAPI) Configure(
	_ context.Context,
	request datasource.ConfigureRequest,
	response *datasource.ConfigureResponse,
) {
	coreClient := getCoreClient(request.ProviderData, &response.Diagnostics)
	if coreClient == nil {
		return
	}

	d.DedicatedserverAPI = coreClient.DedicatedserverAPI
	d.PubliccloudAPI = coreClient.PubliccloudAPI
	d.DNSAPI = coreClient.DNSAPI
	d.IPmgmtAPI = coreClient.IPmgmtAPI
}

func (d *DataSourceAPI) Metadata(
	_ context.Context,
	request datasource.MetadataRequest,
	response *datasource.MetadataResponse,
) {
	response.TypeName = generateTypeName(request.ProviderTypeName, d.Name)
}
