package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct{}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"child_account_data_basic", nil)
}
