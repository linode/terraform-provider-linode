package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label                string
	Region               string
	PlacementGroupType   string
	PlacementGroupPolicy string
}

func DataBasic(t *testing.T, label, region, placementGroupType string, placementGroupPolicy string) string {
	return acceptance.ExecuteTemplate(t,
		"placement_group_data_basic", TemplateData{
			Label:                label,
			Region:               region,
			PlacementGroupType:   placementGroupType,
			PlacementGroupPolicy: placementGroupPolicy,
		})
}

func Basic(t *testing.T, label, region, placementGroupType string, placementGroupPolicy string) string {
	return acceptance.ExecuteTemplate(t,
		"placement_group_basic", TemplateData{
			Label:                label,
			Region:               region,
			PlacementGroupType:   placementGroupType,
			PlacementGroupPolicy: placementGroupPolicy,
		})
}
