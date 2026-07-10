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

func DataBasic(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destinations_data_basic", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
		})
}

func BucketOnly(t testing.TB, label, region, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"monitor_logs_destinations_bucket_only", TemplateData{
			Label:   label,
			Region:  region,
			Cluster: cluster,
		})
}
