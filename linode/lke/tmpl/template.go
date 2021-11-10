package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label            string
	K8sVersion       string
	HighAvailability bool
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

func Autoscaler(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler", TemplateData{Label: name})
}

func AutoscalerUpdates(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_updates", TemplateData{Label: name})
}

func AutoscalerManyPools(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_many_pools", TemplateData{Label: name})
}

func ControlPlane(t *testing.T, name string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_control_plane", TemplateData{Label: name, HighAvailability: ha})
}

func DataBasic(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_basic", TemplateData{Label: name})
}

func DataAutoscaler(t *testing.T, name string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_autoscaler", TemplateData{Label: name})
}

func DataControlPlane(t *testing.T, name string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_control_plane", TemplateData{Label: name, HighAvailability: ha})
}
