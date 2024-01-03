package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
	Port   int
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nb_configs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilter(t *testing.T, label, region string, port int) string {
	return acceptance.ExecuteTemplate(t,
		"nb_configs_data_filter", TemplateData{
			Label:  label,
			Region: region,
			Port:   port,
		})
}
