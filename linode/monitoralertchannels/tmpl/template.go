package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	AlertType string
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_alert_channels_data_basic", nil)
}

func DataFilter(t testing.TB, alertType string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_alert_channels_data_filter", TemplateData{
			AlertType: alertType,
		})
}
