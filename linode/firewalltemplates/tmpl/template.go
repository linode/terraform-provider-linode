package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Slug string
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"data_linode_firewall_templates_basic", TemplateData{})
}

func DataFilter(t testing.TB, slug string) string {
	return acceptance.ExecuteTemplate(t,
		"data_linode_firewall_templates_filter", TemplateData{
			Slug: slug,
		},
	)
}
