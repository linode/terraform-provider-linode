package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	ClusterID string
	Tag       string
}

func Basic(t *testing.T, clusterID string, tag string) string {
	return acceptance.ExecuteTemplate(t,
		"nodepool_basic", TemplateData{
			ClusterID: clusterID,
			Tag:       tag,
		})
}
