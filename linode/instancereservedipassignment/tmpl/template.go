package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label            string
	ApplyImmediately bool
	Region           string
	Address          string
}

func AddReservedIP(t *testing.T, instanceLabel, region string, address string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_add_reservedIP", TemplateData{
			Label:   instanceLabel,
			Region:  region,
			Address: address,
		})
}
