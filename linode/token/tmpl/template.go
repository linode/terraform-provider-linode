package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label string
}

type TemplateNewExpiryData struct {
	TemplateData
	Expiry string
}

type TemplateNewScopesData struct {
	TemplateData
	Scopes string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"token_basic", TemplateData{Label: label})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"token_updates", TemplateData{Label: label})
}

func RecreateNewScopes(t *testing.T, label, scopes string) string {
	return acceptance.ExecuteTemplate(t,
		"token_recreate_new_scopes", TemplateNewScopesData{
			TemplateData: TemplateData{Label: label},
			Scopes:       scopes,
		})
}

func RecreateNewExpiryDate(t *testing.T, label, expiry string) string {
	return acceptance.ExecuteTemplate(t,
		"token_recreate_new_expiry_date", TemplateNewExpiryData{
			TemplateData: TemplateData{Label: label},
			Expiry:       expiry,
		})
}
