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
		"template_basic", TemplateData{Label: label})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"template_updates", TemplateData{Label: label})
}
