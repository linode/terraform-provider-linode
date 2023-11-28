package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Engine          string
	Label           string
	Region          string
	AllowedIP       string
	ReplicationType string
	ClusterSize     int
	Encrypted       bool
	SSLConnection   bool
}

func Basic(t *testing.T, label, engine, region string) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_basic", TemplateData{
			Engine: engine,
			Label:  label,
			Region: region,
		})
}

func Complex(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_complex", data)
}

func ComplexUpdates(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_complex_updates", data)
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_data_basic", data)
}
