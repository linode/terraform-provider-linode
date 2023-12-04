package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
	IPv4   string
}

func DataBasic(t *testing.T, label, region, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnets_data_basic", TemplateData{
			Label:  label,
			Region: region,
			IPv4:   ipv4,
		})
}

func DataFilterLabel(t *testing.T, label, region, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnets_data_filter_label", TemplateData{
			Label:  label,
			Region: region,
			IPv4:   ipv4,
		})
}
