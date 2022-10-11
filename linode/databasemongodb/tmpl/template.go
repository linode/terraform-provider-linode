package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
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
	Region          string
}

func Basic(t *testing.T, label, engine, region string) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongodb_basic", TemplateData{
			Engine: engine,
			Label:  label,
			Region: region,
		})
}

func Complex(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongodb_complex", data)
}

func ComplexUpdates(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongodb_complex_updates", data)
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mongodb_data_basic", data)
}
