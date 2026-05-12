package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label string
}

func DataSource(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"tag_datasource", TemplateData{
			Label: label,
		})
}
