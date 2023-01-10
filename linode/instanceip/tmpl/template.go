package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label            string
	ApplyImmediately bool
	Region           string
}

func Basic(t *testing.T, instanceLabel, region string, applyImmediately bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_basic", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
			Region:           region,
		})
}

func DataBasic(t *testing.T, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_data_basic", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}

func NoBoot(t *testing.T, instanceLabel, region string, applyImmediately bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_no_boot", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
			Region:           region,
		})
}
