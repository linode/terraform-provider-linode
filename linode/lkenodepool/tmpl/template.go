package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TaintData struct {
	Effect string
	Key    string
	Value  string
}

type TemplateData struct {
	ClusterID         string
	ClusterLabel      string
	Region            string
	K8sVersion        string
	PoolTag           string
	PoolNodeType      string
	NodeCount         int
	AutoscalerEnabled bool
	AutoscalerMin     int
	AutoscalerMax     int
	Taints            []TaintData
	Labels            map[string]string
	Label             string
	FirewallID        *int
	UpdateStrategy    string
	DiskEncryption    string
}

func Generate(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "nodepool_template", *data)
}

func EnterpriseBasic(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "lke_e_nodepool", *data)
}

func DataBasic(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "lke_nodepool_data_basic", *data)
}

func DataClusterNotFound(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "lke_nodepool_data_cluster_not_found", *data)
}

func DataNodePoolNotFound(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "lke_nodepool_data_nodepool_not_found", *data)
}
