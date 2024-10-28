package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	ID string
}

func DataBasic(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"obj_cluster_data_basic", TemplateData{ID: id})
}
