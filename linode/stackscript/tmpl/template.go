package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Script string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscript_basic", TemplateData{Label: label})
}

func CodeChange(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscript_code_change", TemplateData{Label: label})
}

func DataBasic(t *testing.T, script string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscript_data_basic", TemplateData{
			Script: script,
		})
}
