package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"nbs_data_basic", TemplateData{
			Label:  label,
			Region: region,
		})
}
