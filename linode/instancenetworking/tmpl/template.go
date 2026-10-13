package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label         string
	Region        string
	RootPass      string
	IPv4          string
	InterfaceIPv4 string
}

func DataBasic(t testing.TB, instanceLabel, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_basic", TemplateData{
			Label:    instanceLabel,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataVPC(t testing.TB, label, region, subnetIPv4, interfaceIPv4, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_vpc", TemplateData{
			Label:         label,
			Region:        region,
			RootPass:      rootPass,
			IPv4:          subnetIPv4,
			InterfaceIPv4: interfaceIPv4,
		})
}

func DataBasic_withReservedField(t *testing.T, instanceLabel, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_basic_with_reserved", TemplateData{
			Label:    instanceLabel,
			Region:   region,
			RootPass: rootPass,
		})
}

func DataVPCDualStack(t testing.TB, label, region, subnetIPv4, interfaceIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_vpc_dual_stack", TemplateData{
			Label:         label,
			Region:        region,
			IPv4:          subnetIPv4,
			InterfaceIPv4: interfaceIPv4,
		})
}
