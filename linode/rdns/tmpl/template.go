package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_basic", TemplateData{Label: label})
}

func Changed(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_changed", TemplateData{Label: label})
}

func Deleted(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"rdns_deleted", TemplateData{Label: label})
}
