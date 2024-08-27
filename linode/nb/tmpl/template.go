package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func Updates(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_updates", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataBasic(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataFirewalls(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_firewalls", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func Firewall(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_firewall", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func FirewallUpdate(t *testing.T, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_firewall_updates", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}
