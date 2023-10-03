package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_basic", nil)
}

func DataFilterLabel(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpcs_data_filter_label", TemplateData{
			Label:  label,
			Region: region,
		})
}
