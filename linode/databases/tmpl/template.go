package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	databasemysqlv2tmpl "github.com/linode/terraform-provider-linode/v3/linode/databasemysqlv2/tmpl"
)

type TemplateData struct {
	DB     databasemysqlv2tmpl.TemplateData
	Engine string
	Label  string
}

func ByLabel(t testing.TB, engineVersion, instLabel, dsLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_label", TemplateData{
			DB:    databasemysqlv2tmpl.TemplateData{EngineID: engineVersion, Label: instLabel, Region: region},
			Label: dsLabel,
		})
}

func ByEngine(t testing.TB, engineVersion, label, engine, region string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_engine", TemplateData{
			DB:     databasemysqlv2tmpl.TemplateData{EngineID: engineVersion, Label: label, Region: region},
			Engine: engine,
		})
}
