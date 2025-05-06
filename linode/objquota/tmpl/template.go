package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	QuotaID string
}

func DataBasic(t testing.TB, quotaID string) string {
	return acceptance.ExecuteTemplate(t,
		"object_quota_data_basic", TemplateData{
			QuotaID: quotaID,
		})
}
