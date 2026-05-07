package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	VolumeLabel string
	Region      string
}

func DataBasic(t testing.TB, label, region string) string {
	return acceptance.ExecuteTemplate(t,
		"iam_entities_data_basic", TemplateData{
			VolumeLabel: label,
			Region:      region,
		})
}
