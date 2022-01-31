package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Script string
}

func DataBasic(t *testing.T, label, script string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscripts_data_basic", TemplateData{
			Label:  label,
			Script: script,
		})
}

func DataSubString(t *testing.T, label, script string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscripts_data_substring", TemplateData{
			Label:  label,
			Script: script,
		})
}

func DataLatest(t *testing.T, label, script string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscripts_data_latest", TemplateData{
			Label:  label,
			Script: script,
		})
}

func DataClientFilter(t *testing.T, label, script string) string {
	return acceptance.ExecuteTemplate(t,
		"stackscripts_data_clientfilter", TemplateData{
			Label:  label,
			Script: script,
		})
}
