package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	ApplyImmediately bool
	Region           string
	RootPass         string
}

func AddReservedIP(t *testing.T, instanceLabel, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_add_reservedIP", TemplateData{
			Label:    instanceLabel,
			Region:   region,
			RootPass: rootPass,
		})
}
