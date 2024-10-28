package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func Basic(t *testing.T, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ips_basic", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}
