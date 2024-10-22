package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Detached(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_detached", TemplateData{
			Label:  label,
			Region: region,
		})
}

func WithNodeBalancer(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_with_nodebalancer", TemplateData{
			Label:  label,
			Region: region,
		})
}
