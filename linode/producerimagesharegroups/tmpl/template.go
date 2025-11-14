package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Label1 string
	Label2 string
}

func DataBasic(t testing.TB, label1, label2 string) string {
	return acceptance.ExecuteTemplate(t,
		"producer_image_share_groups_data_basic", TemplateData{
			Label1: label1,
			Label2: label2,
		})
}
