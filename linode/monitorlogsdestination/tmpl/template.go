package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label   string
	Region  string
	Cluster string
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
