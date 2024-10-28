package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Region string
	Label  string
}

func DataBasic(t testing.TB, id, label string) string {
	return acceptance.ExecuteTemplate(t,
		"region_data_basic", TemplateData{Region: id, Label: label})
}
