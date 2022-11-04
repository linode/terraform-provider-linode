package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label   string
	Cluster string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_basic", TemplateData{Label: label})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_updates", TemplateData{Label: label})
}

func Limited(t *testing.T, label, cluster string) string {
	return acceptance.ExecuteTemplate(t,
		"object_key_limited", TemplateData{Label: label, Cluster: cluster})
}
