package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Updates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_updates", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func InvalidType(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_invalid_type", TemplateData{
			Label: label,
		})
}

func DataNotFound(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_data_not_found", nil)
}
