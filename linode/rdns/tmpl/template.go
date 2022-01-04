package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label            string
	WaitForAvailable bool
}

func Basic(t *testing.T, label string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_basic", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
		})
}

func Changed(t *testing.T, label string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_changed", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
		})
}

func Deleted(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_deleted", TemplateData{
			Label: label,
		})
}
