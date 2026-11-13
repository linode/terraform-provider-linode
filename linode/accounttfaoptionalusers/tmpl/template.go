package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Username string
	Email    string
	Optional bool
}

func Basic(t testing.TB, username, email string, optional bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_tfa_optional_users_basic", TemplateData{
			Username: username,
			Email:    email,
			Optional: optional,
		})
}

func DataBasic(t testing.TB, username, email string, optional bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_tfa_optional_users_data_basic", TemplateData{
			Username: username,
			Email:    email,
			Optional: optional,
		})
}
