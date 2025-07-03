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
	UpdateStrategy    string
}

func Generate(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "nodepool_template", *data)
}

func EnterpriseBasic(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "lke_e_nodepool", *data)
}
