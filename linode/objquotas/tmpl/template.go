package tmpl

import (
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"testing"
)

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t, "object_storage_quotas_basic", nil)
}
