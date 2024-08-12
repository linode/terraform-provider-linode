package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateDataBasic struct {
	Label  string
	Region string
	IPv4_1 string
	IPv4_2 string
}

type TemplateDataFilter struct {
	Label  string
	Region string
	IPv4   string
}

func DataBasic(t *testing.T, label, region, ipv4_1, ipv4_2 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_ips_data_basic", TemplateDataBasic{
			Label:  label,
			Region: region,
			IPv4_1: ipv4_1,
			IPv4_2: ipv4_2,
		})
}

func DataFilterAddress(t *testing.T, label, region, ipv4 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_ips_data_filter_address", TemplateDataFilter{
			Label:  label,
			Region: region,
			IPv4:   ipv4,
		})
}
