package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct{}

func DataAll(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"database_engines_data_all", nil)
}

func DataByEngine(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"database_engines_data_by_engine", nil)
}
