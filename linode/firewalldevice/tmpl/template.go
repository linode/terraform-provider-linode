package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_basic", TemplateData{
			Label: label,
		})
}

func Detached(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"firewall_device_detached", TemplateData{
			Label: label,
		})
}
