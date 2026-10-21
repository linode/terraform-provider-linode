package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	K8sVersion       string
	HighAvailability bool
	Region           string
	Tag1Name         string
	Tag2Name         string
}

func DataBasic(t testing.TB, name, version, region, tag1Name, tag2Name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_clusters_data_basic", TemplateData{Label: name, K8sVersion: version, Region: region, Tag1Name: tag1Name, Tag2Name: tag2Name})
}

func DataFilter(t testing.TB, name, version, region, tag1Name, tag2Name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_clusters_data_filter", TemplateData{Label: name, K8sVersion: version, Region: region, Tag1Name: tag1Name, Tag2Name: tag2Name})
}
