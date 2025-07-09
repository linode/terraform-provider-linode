package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"child_accounts_data_basic", nil)
}
