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
	Inherit          bool
	SquashPolicy     string
}

func Basic(
	t testing.TB,
	label, region string,
	maxCapacityBytes int64,
	enabled, inherit bool,
	squashPolicy string,
) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_access_policy_basic", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
			Enabled:          enabled,
			Inherit:          inherit,
			SquashPolicy:     squashPolicy,
		})
}

func DataUnconfigured(t testing.TB, label, region string, maxCapacityBytes int64) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_access_policy_data_unconfigured", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
		})
}

func DataConfigured(
	t testing.TB,
	label, region string,
	maxCapacityBytes int64,
	enabled, inherit bool,
	squashPolicy string,
) string {
	return acceptance.ExecuteTemplate(t,
		"nfs_filesystem_access_policy_data_configured", TemplateData{
			Label:            label,
			Region:           region,
			MaxCapacityBytes: maxCapacityBytes,
			Enabled:          enabled,
			Inherit:          inherit,
			SquashPolicy:     squashPolicy,
		})
}
