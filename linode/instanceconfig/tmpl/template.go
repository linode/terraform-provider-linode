package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"testing"
)

type TemplateData struct {
	Label string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_basic", TemplateData{
			Label: label,
		})
}

func Complex(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex", TemplateData{
			Label: label,
		})
}

func ComplexUpdates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_config_complex_updates", TemplateData{
			Label: label,
		})
}
