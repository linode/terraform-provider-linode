package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"account_data_basic", nil)
}
