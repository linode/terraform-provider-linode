package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label string
}

func DataBasic(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_spaces_data_basic", TemplateData{
			Label: label,
		})
}

func DataFilterLabel(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_spaces_data_filter_label", TemplateData{
			Label: label,
		})
}
