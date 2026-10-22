// Package client implements access to facades.
package client

import (
	"github.com/leaseweb/leaseweb-go-sdk/dedicatedserver/v2"
	"github.com/leaseweb/leaseweb-go-sdk/dns"
	"github.com/leaseweb/leaseweb-go-sdk/ipmgmt"
	"github.com/leaseweb/leaseweb-go-sdk/publiccloud"
)

const userAgentBase = "leaseweb-terraform"

// The Client handles instantiation of the SDK.
type Client struct {
	PubliccloudAPI     publiccloud.PubliccloudAPI
	DedicatedserverAPI dedicatedserver.DedicatedserverAPI
	DNSAPI             dns.DnsAPI
	IPmgmtAPI          ipmgmt.IpmgmtAPI
}

type Optional struct {
	Host   *string
	Scheme *string
}

func NewClient(token string, optional Optional, version string) Client {
	publiccloudCFG := publiccloud.NewConfiguration()
	dedicatedserverCFG := dedicatedserver.NewConfiguration()
	dnsCFG := dns.NewConfiguration()
	ipmgmtCFG := ipmgmt.NewConfiguration()

	if optional.Host != nil {
		publiccloudCFG.Host = *optional.Host
		dedicatedserverCFG.Host = *optional.Host
		dnsCFG.Host = *optional.Host
		ipmgmtCFG.Host = *optional.Host
	}
	if optional.Scheme != nil {
		publiccloudCFG.Scheme = *optional.Scheme
		dedicatedserverCFG.Scheme = *optional.Scheme
		dnsCFG.Scheme = *optional.Scheme
		ipmgmtCFG.Scheme = *optional.Scheme
	}

	userAgent := userAgentBase + "-" + version

	publiccloudCFG.AddDefaultHeader("X-LSW-Auth", token)
	publiccloudCFG.UserAgent = userAgent

	dedicatedserverCFG.AddDefaultHeader("X-LSW-Auth", token)
	dedicatedserverCFG.UserAgent = userAgent

	dnsCFG.AddDefaultHeader("X-LSW-Auth", token)
	dnsCFG.UserAgent = userAgent

	ipmgmtCFG.AddDefaultHeader("X-LSW-Auth", token)
	ipmgmtCFG.UserAgent = userAgent

	publiccloudAPI := publiccloud.NewAPIClient(publiccloudCFG)
	dedicatedserverAPI := dedicatedserver.NewAPIClient(dedicatedserverCFG)
	dnsAPI := dns.NewAPIClient(dnsCFG)
	ipmgmtAPI := ipmgmt.NewAPIClient(ipmgmtCFG)

	return Client{
		PubliccloudAPI:     publiccloudAPI.PubliccloudAPI,
		DedicatedserverAPI: dedicatedserverAPI.DedicatedserverAPI,
		DNSAPI:             dnsAPI.DnsAPI,
		IPmgmtAPI:          ipmgmtAPI.IpmgmtAPI,
	}
}
