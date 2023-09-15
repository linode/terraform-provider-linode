package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Detached(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_detached", TemplateData{
			Label:  label,
			Region: region,
		})
}

func WithNodeBalancer(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_with_nodebalancer", TemplateData{
			Label:  label,
			Region: region,
		})
}
