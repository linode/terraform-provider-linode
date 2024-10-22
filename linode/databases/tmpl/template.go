package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	databasemysqltmpl "github.com/linode/terraform-provider-linode/v2/linode/databasemysql/tmpl"
)

type TemplateData struct {
	DB     databasemysqltmpl.TemplateData
	Engine string
	Label  string
}

func ByLabel(t testing.TB, engineVersion, instLabel, dsLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_label", TemplateData{
			DB:    databasemysqltmpl.TemplateData{Engine: engineVersion, Label: instLabel, Region: region},
			Label: dsLabel,
		})
}

func ByEngine(t testing.TB, engineVersion, label, engine, region string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_engine", TemplateData{
			DB:     databasemysqltmpl.TemplateData{Engine: engineVersion, Label: label, Region: region},
			Engine: engine,
		})
}
