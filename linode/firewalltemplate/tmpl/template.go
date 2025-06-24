package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Slug string
}

func DataBasic(t testing.TB, slug string) string {
	return acceptance.ExecuteTemplate(t,
		"data_linode_firewall_template_basic", TemplateData{
			Slug: slug,
		})
}
