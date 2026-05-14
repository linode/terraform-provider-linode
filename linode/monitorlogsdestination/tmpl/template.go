package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label       string
	Region      string
	EndpointURL string
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

func CustomHTTPS(t testing.TB, label, endpointURL string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_custom_https", TemplateData{
			Label:       label,
			EndpointURL: endpointURL,
		})
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}
