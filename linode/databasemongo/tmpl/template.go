package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"

	"testing"
)

type TemplateData struct {
	Engine          string
	Label           string
	AllowedIP       string
	ClusterSize     int
	Encrypted       bool
	SSLConnection   bool
	StorageEngine   string
	CompressionType string
}

func Basic(t *testing.T, label, engine string) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongo_basic", TemplateData{
			Engine: engine,
			Label:  label,
		})
}

func Complex(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongo_complex", data)
}

func ComplexUpdates(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongo_complex_updates", data)
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongo_data_basic", data)
}
