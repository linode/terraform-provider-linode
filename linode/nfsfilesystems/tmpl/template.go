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

func DataBasic(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystems_data_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}

func DataFilterLabel(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystems_data_filter_label", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}
