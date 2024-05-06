package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label        string
	Region       string
	AffinityType string
	IsStrict     bool
}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"placement_group_data_basic", nil)
}

func Basic(t *testing.T, label, region, affinityType string, isStrict bool) string {
	return acceptance.ExecuteTemplate(t,
		"placement_group_basic", TemplateData{
			Label:        label,
			Region:       region,
			AffinityType: affinityType,
			IsStrict:     isStrict,
		})
}
