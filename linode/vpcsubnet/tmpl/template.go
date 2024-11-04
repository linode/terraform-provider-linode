package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	IPv4   string
	Region string
}

func Basic(t testing.TB, label, ipv4, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_basic", TemplateData{
			Label:  label,
			IPv4:   ipv4,
			Region: region,
		})
}

func Updates(t testing.TB, label, ipv4, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_updates", TemplateData{
			Label:  label,
			IPv4:   ipv4,
			Region: region,
		})
}

func DataBasic(t testing.TB, label, ipv4, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_data_basic", TemplateData{
			Label:  label,
			IPv4:   ipv4,
			Region: region,
		})
}

func Attached(t testing.TB, label, ipv4, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_subnet_attached", TemplateData{
			Label:  label,
			IPv4:   ipv4,
			Region: region,
		})
}
