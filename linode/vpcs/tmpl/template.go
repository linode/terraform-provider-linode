package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label     string
	Region    string
	IPv4Range string
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataDualStack(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_dual_stack", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilterLabel(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_filter_label", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataIPv4(t testing.TB, label, region, ipv4Range string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_ipv4", TemplateData{
			Label:     label,
			Region:    region,
			IPv4Range: ipv4Range,
		})
}
