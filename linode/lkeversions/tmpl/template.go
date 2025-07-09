package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateDataTier struct {
	Tier string
}

func DataNoTier(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"lke_versions_data_no_tier", nil)
}

func DataTier(t testing.TB, tier string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_versions_data_tier", TemplateDataTier{
			Tier: tier,
		})
}
