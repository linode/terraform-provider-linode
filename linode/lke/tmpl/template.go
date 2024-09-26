package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TaintData struct {
	Effect string
	Key    string
	Value  string
}

type TemplateData struct {
	Label            string
	K8sVersion       string
	HighAvailability bool
	Region           string
	ACLEnabled       bool
	IPv4             string
	IPv6             string
	Taints           []TaintData
	Labels           map[string]string
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

func ControlPlane(t *testing.T, name, version, region, ipv4, ipv6 string, ha, enabled bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_control_plane", TemplateData{
			Label:            name,
			HighAvailability: ha,
			K8sVersion:       version,
			Region:           region,
			IPv4:             ipv4,
			IPv6:             ipv6,
			ACLEnabled:       enabled,
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

func DataControlPlane(t *testing.T, name, version, region, ipv4, ipv6 string, ha, enabled bool) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_control_plane", TemplateData{
			Label:            name,
			HighAvailability: ha,
			K8sVersion:       version,
			Region:           region,
			IPv4:             ipv4,
			IPv6:             ipv6,
			ACLEnabled:       enabled,
		})
}

func DataTaintsLabels(t *testing.T, name, version, region string, taints []TaintData, labels map[string]string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_cluster_data_taints_labels", TemplateData{
			Label:      name,
			K8sVersion: version,
			Region:     region,
			Labels:     labels,
			Taints:     taints,
		})
}
