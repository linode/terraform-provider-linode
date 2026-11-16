package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label      string
	Region     string
	SubnetIPv4 string
}

func DataBasic(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_interfaces_data_basic", TemplateData{
		Label:      label,
		Region:     region,
		SubnetIPv4: subnetIPv4,
	})
}

func DataAddress(t testing.TB, label, region, subnetIPv4 string) string {
	return acceptance.ExecuteTemplate(t, "natgateway_interfaces_address", TemplateData{
		Label:      label,
		Region:     region,
		SubnetIPv4: subnetIPv4,
	})
}
