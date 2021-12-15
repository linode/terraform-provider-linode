package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label            string
	K8sVersion       string
	HighAvailability bool
	PoolTag          string
}

func Basic(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_basic", TemplateData{Label: name, K8sVersion: version})
}

func Updates(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_updates", TemplateData{Label: name, K8sVersion: version})
}

func ManyPools(t *testing.T, name, k8sVersion string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_many_pools", TemplateData{
			Label:      name,
			K8sVersion: k8sVersion,
		})
}

func ComplexPools(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_complex_pools", TemplateData{Label: name, K8sVersion: version})
}

func Autoscaler(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler", TemplateData{Label: name, K8sVersion: version})
}

func AutoscalerUpdates(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_updates", TemplateData{Label: name, K8sVersion: version})
}

func AutoscalerManyPools(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_many_pools", TemplateData{Label: name, K8sVersion: version})
}

func ControlPlane(t *testing.T, name, version string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_control_plane", TemplateData{Label: name, HighAvailability: ha, K8sVersion: version})
}

func PoolBasic(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_pool_basic", TemplateData{Label: name, K8sVersion: version})
}

func PoolTag(t *testing.T, name, version, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_pool_tag", TemplateData{Label: name, K8sVersion: version, PoolTag: tag})
}

func DataBasic(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_basic", TemplateData{Label: name, K8sVersion: version})
}

func DataAutoscaler(t *testing.T, name, version string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_autoscaler", TemplateData{Label: name, K8sVersion: version})
}

func DataControlPlane(t *testing.T, name, version string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_control_plane", TemplateData{Label: name, HighAvailability: ha, K8sVersion: version})
}
