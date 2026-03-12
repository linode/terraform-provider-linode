package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label        string
	AlertChannel int
}

func DataBasic(t testing.TB, label string, alertChannel int) string {
	return acceptance.ExecuteTemplate(t,
		"alert_definitions_data_basic", TemplateData{
			Label:        label,
			AlertChannel: alertChannel,
		})
}

func DataFilter(t testing.TB, label string, alertChannel int) string {
	return acceptance.ExecuteTemplate(t,
		"alert_definitions_data_filter", TemplateData{
			Label:        label,
			AlertChannel: alertChannel,
		})
}
