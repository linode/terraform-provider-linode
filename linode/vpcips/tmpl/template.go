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

func DataBasic(t testing.TB, label, region, ipv4_1, ipv4_2 string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_ips_data_basic", TemplateDataBasic{
			Label:  label,
			Region: region,
			IPv4_1: ipv4_1,
			IPv4_2: ipv4_2,
		})
}
