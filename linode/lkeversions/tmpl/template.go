package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"lke_versions_data_basic", nil)
}
