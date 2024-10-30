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
}

func Generate(t testing.TB, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "nodepool_template", *data)
}
