package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Domain string
}

func DataBasic(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domains_data_basic", TemplateData{Domain: domain})
}

func DataFilter(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domains_data_filter", TemplateData{Domain: domain})
}

func DataAPIFilter(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domains_data_api_filter", TemplateData{Domain: domain})
}
