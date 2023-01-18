package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_networking_data_basic", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}
