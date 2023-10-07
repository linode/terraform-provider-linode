package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataBasic(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volumes_data_basic", TemplateData{Label: volume, Region: region})
}
