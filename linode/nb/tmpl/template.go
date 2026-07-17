package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label          string
	Region         string
	Type           string
	K8sVersion     string
	NodeBalancerID int
}

func Basic(t testing.TB, nodebalancer, region, nbType string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
			Type:   nbType,
		})
}

func Updates(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_updates", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func DataBasic(t testing.TB, nodebalancer, region, nbType string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_basic", TemplateData{
			Label:  nodebalancer,
			Region: region,
			Type:   nbType,
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

func FrontendVPC(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_frontend_vpc", TemplateData{
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

func DataFrontendVPC(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_frontend_vpc", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func ReservedIP(t testing.TB, nodebalancer, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_reserved_ip", TemplateData{
			Label:  nodebalancer,
			Region: region,
		})
}

func ReservedIPOnly(t testing.TB, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_reserved_ip_only", TemplateData{
			Region: region,
		})
}

func LKEClusterData(t testing.TB, nodeBalancerID int) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_data_lke_cluster", TemplateData{
			NodeBalancerID: nodeBalancerID,
		})
}
