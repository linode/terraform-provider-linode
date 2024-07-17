package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label   string
	Cluster string
	Region  string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_basic", TemplateData{Label: label})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_updates", TemplateData{Label: label})
}

func ClusterLimited(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_limited", TemplateData{Label: label, Cluster: cluster})
}

func Limited(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_limited", TemplateData{Label: label, Region: region})
}
