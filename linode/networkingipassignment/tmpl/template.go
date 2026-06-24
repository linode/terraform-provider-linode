package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Region   string
	RootPass string
}

func NetworkingIPsAssign(t *testing.T, label, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"networking_ips_assign",
		TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}
