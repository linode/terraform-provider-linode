package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct{}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"profile_data_basic", nil)
}
