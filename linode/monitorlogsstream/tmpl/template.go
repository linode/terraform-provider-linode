package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label    string
	StreamID string
	Region   string
}

func Lifecycle(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_lifecycle", TemplateData{
		Label:  label,
		Region: region,
	})
}

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_basic", TemplateData{
		Label:  label,
		Region: region,
	})
}

func Updates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_updates", TemplateData{
		Label:  label,
		Region: region,
	})
}

func DataBasic(t testing.TB, streamID string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_data_basic", TemplateData{
		StreamID: streamID,
	})
}

func DataNotFound(t testing.TB) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_data_not_found", TemplateData{})
}

func InvalidType(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_invalid_type", TemplateData{
		Label: label,
	})
}

func InvalidDestination(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t, "monitor_logs_stream_invalid_destination", TemplateData{
		Label: label,
	})
}
