package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label    string
	Region   string
	LockType string
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"locks_data_basic", TemplateData{
			Label:    label,
			Region:   region,
			LockType: "cannot_delete",
		})
}

func DataFilter(t testing.TB, label, region, lockType string) string {
	return acceptance.ExecuteTemplate(t,
		"locks_data_filter", TemplateData{
			Label:    label,
			Region:   region,
			LockType: lockType,
		})
}
