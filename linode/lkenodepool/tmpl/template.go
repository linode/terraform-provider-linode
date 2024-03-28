package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	ClusterID         string
	ClusterLabel      string
	Region            string
	K8sVersion        string
	PoolTag           string
	NodeCount         int
	AutoscalerEnabled bool
	AutoscalerMin     int
	AutoscalerMax     int
}

func Generate(t *testing.T, data *TemplateData) string {
	return acceptance.ExecuteTemplate(t, "nodepool_template", *data)
}
