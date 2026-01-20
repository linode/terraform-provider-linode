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

func Basic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lock_basic", TemplateData{
			Label:    label,
			Region:   region,
			LockType: "cannot_delete",
		})
}

func WithSubresources(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lock_basic", TemplateData{
			Label:    label,
			Region:   region,
			LockType: "cannot_delete_with_subresources",
		})
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"lock_data_basic", TemplateData{
			Label:    label,
			Region:   region,
			LockType: "cannot_delete",
		})
}
