package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Username   string
	Email      string
	Restricted bool
	InstLabel  string
}

func DataBasic(t testing.TB, username, email string) string {
	return acceptance.ExecuteTemplate(t,
		"users_data_basic", TemplateData{
			Username: username,
			Email:    email,
		})
}

func DataClientFilter(t testing.TB, username, email string) string {
	return acceptance.ExecuteTemplate(t,
		"users_data_clientfilter", TemplateData{
			Username: username,
			Email:    email,
		})
}

func DataSubstring(t testing.TB, username, email string) string {
	return acceptance.ExecuteTemplate(t,
		"users_data_substring", TemplateData{
			Username: username,
			Email:    email,
		})
}
