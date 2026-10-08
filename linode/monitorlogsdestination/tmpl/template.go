package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label       string
	Region      string
	Cluster     string
	EndpointURL string
	Username    string
	Password    string
}

func Basic(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_basic", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
		})
}

func Updates(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_updates", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
		})
}

func BucketOnly(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_bucket_only", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
		})
}

func DataBasic(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_data_basic", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
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

func TrafficPeakBasic(t testing.TB, label, endpointURL, username, password string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_basic", TemplateData{
			Label:       label,
			EndpointURL: endpointURL,
			Username:    username,
			Password:    password,
		})
}

func TrafficPeakUpdates(t testing.TB, label, endpointURL, username, password string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_updates", TemplateData{
			Label:       label,
			EndpointURL: endpointURL,
			Username:    username,
			Password:    password,
		})
}

func TrafficPeakInvalidAuthType(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_invalid_auth_type", TemplateData{Label: label})
}

func TrafficPeakDataBasic(t testing.TB, label, endpointURL, username, password string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_data_basic", TemplateData{
			Label:       label,
			EndpointURL: endpointURL,
			Username:    username,
			Password:    password,
		})
}

func TrafficPeakDataNoAuthType(t testing.TB, label, endpointURL, username, password string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_data_no_auth_type", TemplateData{
			Label:       label,
			EndpointURL: endpointURL,
			Username:    username,
			Password:    password,
		})
}

func TrafficPeakMissingAuth(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_missing_auth", TemplateData{Label: label})
}

func TrafficPeakMissingContentType(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destination_traffic_peak_missing_content_type", TemplateData{Label: label})
}
