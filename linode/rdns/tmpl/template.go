package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	CreateTimeout    string
	UpdateTimeout    string
	WaitForAvailable bool
}

func Basic(t testing.TB, label, region string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_basic", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
			Region:           region,
		})
}

func Changed(t testing.TB, label, region string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_changed", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
			Region:           region,
		})
}

func Deleted(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_deleted", TemplateData{
			Label:  label,
			Region: region,
		})
}

func WithTimeout(t testing.TB, label, region, createTimeout, updateTimeout string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_with_timeout", TemplateData{
			Label:         label,
			Region:        region,
			CreateTimeout: createTimeout,
			UpdateTimeout: updateTimeout,
		})
}

func WithTimeoutUpdated(t testing.TB, label, region, createTimeout, updateTimeout string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_with_timeout_updated", TemplateData{
			Label:         label,
			Region:        region,
			CreateTimeout: createTimeout,
			UpdateTimeout: updateTimeout,
		})
}
