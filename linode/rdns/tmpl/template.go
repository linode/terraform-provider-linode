package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	RootPass         string
	CreateTimeout    string
	UpdateTimeout    string
	WaitForAvailable bool
}

func Basic(t testing.TB, label, region, rootPass string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_basic", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
			Region:           region,
			RootPass:         rootPass,
		})
}

func Changed(t testing.TB, label, region, rootPass string, waitForAvailable bool) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_changed", TemplateData{
			Label:            label,
			WaitForAvailable: waitForAvailable,
			Region:           region,
			RootPass:         rootPass,
		})
}

func Deleted(t testing.TB, label, region, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_deleted", TemplateData{
			Label:    label,
			Region:   region,
			RootPass: rootPass,
		})
}

func WithTimeout(t testing.TB, label, region, rootPass, createTimeout, updateTimeout string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_with_timeout", TemplateData{
			Label:         label,
			Region:        region,
			RootPass:      rootPass,
			CreateTimeout: createTimeout,
			UpdateTimeout: updateTimeout,
		})
}

func WithTimeoutUpdated(t testing.TB, label, region, rootPass, createTimeout, updateTimeout string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_with_timeout_updated", TemplateData{
			Label:         label,
			Region:        region,
			RootPass:      rootPass,
			CreateTimeout: createTimeout,
			UpdateTimeout: updateTimeout,
		})
}
