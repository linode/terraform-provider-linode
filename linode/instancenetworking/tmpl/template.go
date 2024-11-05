package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label         string
	Region        string
	IPv4          string
	InterfaceIPv4 string
}

func DataBasic(t testing.TB, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_basic", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}

func DataVPC(t testing.TB, label, region, subnetIPv4, interfaceIPv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_vpc", TemplateData{
			Label:         label,
			Region:        region,
			IPv4:          subnetIPv4,
			InterfaceIPv4: interfaceIPv4,
		})
}
