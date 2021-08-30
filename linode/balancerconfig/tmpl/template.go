package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	nodebalancer "github.com/linode/terraform-provider-linode/linode/balancer/tmpl"
)

type TemplateData struct {
	NodeBalancer nodebalancer.TemplateData
	SSLCert      string
	SSLKey       string
}

func Basic(t *testing.T, nodebalancerName string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_config_basic", TemplateData{
			NodeBalancer: nodebalancer.TemplateData{
				Label: nodebalancerName,
			}})
}

func Updates(t *testing.T, nodebalancerName string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_config_updates", TemplateData{
			NodeBalancer: nodebalancer.TemplateData{
				Label: nodebalancerName,
			}})
}

func SSL(t *testing.T, nodebalancerName, cert, privKey string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_config_ssl", TemplateData{
			SSLCert: cert,
			SSLKey:  privKey,
			NodeBalancer: nodebalancer.TemplateData{
				Label: nodebalancerName,
			}})
}

func ProxyProtocol(t *testing.T, nodebalancerName string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_config_proxy_protocol", TemplateData{
			NodeBalancer: nodebalancer.TemplateData{
				Label: nodebalancerName,
			}})
}

func DataBasic(t *testing.T, nodebalancerName string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_config_data_basic", TemplateData{
			NodeBalancer: nodebalancer.TemplateData{
				Label: nodebalancerName,
			}})
}
