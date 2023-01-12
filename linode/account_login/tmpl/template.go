package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	ID int
}

func DataBasic(t *testing.T, id int) string {
	return acceptance.ExecuteTemplate(t,
		"account_login_data_basic", TemplateData{
			ID: id,
		})
}
