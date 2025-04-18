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

func Basic(t testing.TB, username, email string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"user_basic", TemplateData{
			Username:   username,
			Email:      email,
			Restricted: restricted,
		})
}

func Grants(t testing.TB, username, email string) string {
	return acceptance.ExecuteTemplate(t,
		"user_grants", TemplateData{
			Username: username,
			Email:    email,
		})
}

func GrantsUpdate(t testing.TB, username, email, instance string) string {
	return acceptance.ExecuteTemplate(t,
		"user_grants_update", TemplateData{
			Username:  username,
			Email:     email,
			InstLabel: instance,
		})
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"user_data_basic", nil)
}

func DataNoUser(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"user_data_nouser", nil)
}

func DataGrants(t testing.TB, username, email string) string {
	return acceptance.ExecuteTemplate(t,
		"data_grants", TemplateData{
			Username: username,
			Email:    email,
		})
}
