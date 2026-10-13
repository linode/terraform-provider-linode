package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	QuotaID string
}

func DataBasic(t testing.TB, quotaID string) string {
	return acceptance.ExecuteTemplate(t, "object_storage_global_quotas_basic", TemplateData{
		QuotaID: quotaID,
	})
}
