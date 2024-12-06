package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func NetworkingIPsAssign(t *testing.T, label string, region string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ips_assign",
		TemplateData{
			Label:  label,
			Region: region,
		})
}
