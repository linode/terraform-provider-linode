package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilterEmpty(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbs_data_filter_empty", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilter(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbs_data_filter", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataOrder(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbs_data_order", TemplateData{
			Label:  label,
			Region: region,
		})
}
