package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct{}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"account_data_basic", nil)
}
