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

func Basic(t testing.TB, instanceLabel, region string, applyImmediately bool, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_basic", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
			Region:           region,
			RootPass:         rootPass,
		})
}

func NoBoot(t testing.TB, instanceLabel, region string, applyImmediately bool, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_no_boot", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
			Region:           region,
			RootPass:         rootPass,
		})
}
