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

func Basic(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_basic", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func Updates(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_updates", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func ManyPools(t *testing.T, name, k8sVersion, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_many_pools", TemplateData{
			Label:      name,
			K8sVersion: k8sVersion,
			Region:     region,
		})
}

func ComplexPools(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_complex_pools", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func Autoscaler(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func AutoscalerUpdates(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_updates", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func AutoscalerManyPools(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_many_pools", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func ControlPlane(t *testing.T, name, version, region string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_control_plane", TemplateData{
			Label:            name,
			HighAvailability: ha,
			K8sVersion:       version,
			Region:           region,
		})
}

func NoCount(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_no_count", TemplateData{
			Label:      name,
			K8sVersion: version,
			Region:     region,
		})
}

func AutoscalerNoCount(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_autoscaler_no_count", TemplateData{
			Label:      name,
			K8sVersion: version,
			Region:     region,
		})
}

func DataBasic(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_basic", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func DataAutoscaler(t *testing.T, name, version, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_autoscaler", TemplateData{Label: name, K8sVersion: version, Region: region})
}

func DataControlPlane(t *testing.T, name, version, region string, ha bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_control_plane", TemplateData{
			Label:            name,
			HighAvailability: ha,
			K8sVersion:       version,
			Region:           region,
		})
}
