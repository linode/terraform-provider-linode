package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilterLabel(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_filter_label", TemplateData{
			Label:  label,
			Region: region,
		})
}
