package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label            string
	ApplyImmediately bool
}

func Basic(t *testing.T, instanceLabel string, applyImmediately bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_basic", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
		})
}

func NoBoot(t *testing.T, instanceLabel string, applyImmediately bool) string {
	return acceptance.ExecuteTemplate(t,
		"instance_ip_no_boot", TemplateData{
			Label:            instanceLabel,
			ApplyImmediately: applyImmediately,
		})
}
