package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	MaxCapacityBytes int64
}

func Basic(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}

func Updates(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_updates", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}

func DataBasic(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_data_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}
