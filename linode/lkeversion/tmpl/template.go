package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateDataNoTier struct {
	ID string
}

type TemplateDataTier struct {
	ID   string
	Tier string
}

func DataNoTier(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_version_data_no_tier", TemplateDataNoTier{
			ID: id,
		})
}

func DataTier(t testing.TB, id, tier string) string {
	return acceptance.ExecuteTemplate(t,
		"lke_version_data_tier", TemplateDataTier{
			ID:   id,
			Tier: tier,
		})
}
