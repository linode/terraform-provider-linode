package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Updates(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_updates", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}
