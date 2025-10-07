package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbvpcs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilter(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbvpcs_data_filter", TemplateData{
			Label:  label,
			Region: region,
		})
}
