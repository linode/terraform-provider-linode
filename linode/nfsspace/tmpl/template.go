package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label       string
	Description string
}

func Basic(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_space_basic", TemplateData{
			Label: label,
		})
}

func Updates(t testing.TB, label, description string) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_space_updates", TemplateData{
			Label:       label,
			Description: description,
		})
}

func DataBasic(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_space_data_basic", TemplateData{
			Label: label,
		})
}
