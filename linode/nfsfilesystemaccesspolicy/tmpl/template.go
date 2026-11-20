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
