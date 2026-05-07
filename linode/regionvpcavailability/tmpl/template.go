package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	RegionID string
}

func DataBasic(t testing.TB, regionID string) string {
	return acceptance.ExecuteTemplate(t,
		"region_vpc_availability_data_basic", TemplateData{
			RegionID: regionID,
		})
}

func DataNoRegion(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"region_vpc_availability_data_noregion", nil)
}
