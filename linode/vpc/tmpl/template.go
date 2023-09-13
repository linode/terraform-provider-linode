package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func WithSubnets(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_with_subnets", TemplateData{
			Label:  label,
			Region: region,
		})
}

func Updates(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_updates", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"vpc_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}
