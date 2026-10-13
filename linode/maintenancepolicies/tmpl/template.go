package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"maintenance_policies_data_basic", nil)
}
