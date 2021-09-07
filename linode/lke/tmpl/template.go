package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label      string
	K8sVersion string
}

func Basic(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_basic", TemplateData{Label: name})
}

func Updates(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_updates", TemplateData{Label: name})
}

func ManyPools(t *testing.T, name, k8sVersion string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_many_pools", TemplateData{
			Label:      name,
			K8sVersion: k8sVersion,
		})
}

func ComplexPools(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_complex_pools", TemplateData{Label: name})
}

func DataBasic(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_basic", TemplateData{Label: name})
}
