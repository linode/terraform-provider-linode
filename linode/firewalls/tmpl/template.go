package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataAll(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"data_linode_firewalls_all", TemplateData{
			Label:  label,
			Region: region,
		})
}

func DataFilter(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"data_linode_firewalls_filter", TemplateData{
			Label:  label,
			Region: region,
		})
}
