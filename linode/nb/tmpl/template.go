package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func Updates(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_updates", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataBasic(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataFirewalls(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_firewalls", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func Firewall(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_firewall", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func FirewallUpdate(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_firewall_updates", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func VPC(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_vpc", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataVPC(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_vpc", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}
