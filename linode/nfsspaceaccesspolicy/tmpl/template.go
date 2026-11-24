package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label            string
	Region           string
	MaxCapacityBytes int64
	Enabled          bool
}

func Basic(t testing.TB, label, region string, maxCapacityBytes int64, enabled bool) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_space_access_policy_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
			Enabled:          enabled,
		})
}

func DataBasic(t testing.TB, label, region string, maxCapacityBytes int64, enabled bool) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_space_access_policy_data_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
			Enabled:          enabled,
		})
}
