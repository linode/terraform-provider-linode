package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label            string
	K8sVersion       string
	HighAvailability bool
	Region           string
}

func DataBasic(t testing.TB, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_clusters_data_basic", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func DataFilter(t testing.TB, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_clusters_data_filter", TemplateData{Label: name, K8sVersion: version, Region: region})
}
