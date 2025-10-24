package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label        string
	Region       string
	SubnetIPv4   string
	VLANLabel    string
	IPAMAddress  string
	IPv4Address  string
	IPv6Range    string
	IPv4Address2 string
	IPv6Range2   string
}

func VLANBasic(t testing.TB, label, region, vlanLabel, ipamAddress string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vlan_basic", TemplateData{
			Label:       label,
			Region:      region,
			VLANLabel:   vlanLabel,
			IPAMAddress: ipamAddress,
		})
}

func PublicBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func PublicWithIPv4(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_ipv4", TemplateData{
			Label:  label,
			Region: region,
		})
}

func PublicWithIPv6(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_ipv6", TemplateData{
			Label:  label,
			Region: region,
		})
}

func PublicWithIPv4AndIPv6(t testing.TB, label, region, ipv4Address, ipv6Range string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_ipv4_ipv6", TemplateData{
			Label:       label,
			Region:      region,
			IPv4Address: ipv4Address,
			IPv6Range:   ipv6Range,
		})
}

func PublicUpdatedIPv4AndIPv6(t testing.TB, label, region, ipv4Address, ipv4Address2, ipv6Range, ipv6Range2 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_updated_ipv4_ipv6", TemplateData{
			Label:        label,
			Region:       region,
			IPv4Address:  ipv4Address,
			IPv4Address2: ipv4Address2,
			IPv6Range:    ipv6Range,
			IPv6Range2:   ipv6Range2,
		})
}

func VPCBasic(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_basic", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		})
}

func VPCWithIPv4(t testing.TB, label, region, subnetIPv4, ipAddress string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_with_ipv4", TemplateData{
			Label:       label,
			Region:      region,
			SubnetIPv4:  subnetIPv4,
			IPv4Address: ipAddress,
		})
}

func VPCDefaultIP(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_default_ip", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		},
	)
}

func VPCUpdatedIPv4(t testing.TB, label, region, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_updated_ipv4", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: ipv4,
		})
}

func PublicDefaultIP(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_default_ip", TemplateData{
			Label:  label,
			Region: region,
		})
}

func PublicEmptyIPObjects(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_public_empty_ip_objects", TemplateData{
			Label:  label,
			Region: region,
		})
}

func VPCEmptyIPObjects(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_empty_ip_objects", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		})
}

func PublicDefaultRouteIPv6(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"public_default_route_ipv6", TemplateData{
			Label:  label,
			Region: region,
		})
}

func VPCDefaultRouteIPv4(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_default_route_ipv4", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		})
}

func VPCWithIPv60(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_with_ipv6_0", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		})
}

func VPCWithIPv61(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"interface_vpc_with_ipv6_1", TemplateData{
			Label:      label,
			Region:     region,
			SubnetIPv4: subnetIPv4,
		})
}
