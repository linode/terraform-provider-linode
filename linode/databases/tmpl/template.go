package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	databasemysqltmpl "github.com/linode/terraform-provider-linode/linode/databasemysql/tmpl"
)

type TemplateData struct {
	DB     databasemysqltmpl.TemplateData
	Engine string
	Label  string
}

// TODO: resolve this dynamically at runtime
const engineSlug = "mysql/8.0.26"

func ByLabel(t *testing.T, instLabel, dsLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_label", TemplateData{
			DB:    databasemysqltmpl.TemplateData{Engine: engineSlug, Label: instLabel},
			Label: dsLabel,
		})
}

func ByEngine(t *testing.T, label, engine string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_engine", TemplateData{
			DB:     databasemysqltmpl.TemplateData{Engine: engineSlug, Label: label},
			Engine: engine,
		})
}
