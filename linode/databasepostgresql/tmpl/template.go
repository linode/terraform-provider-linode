package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Engine                string
	Label                 string
	AllowedIP             string
	ClusterSize           int
	Encrypted             bool
	SSLConnection         bool
	ReplicationType       string
	ReplicationCommitType string
}

func Basic(t *testing.T, label, engine string) string {
	return acceptance.ExecuteTemplate(t,
		"database_postgresql_basic", TemplateData{
			Engine: engine,
			Label:  label,
		})
}

func Complex(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_postgresql_complex", data)
}

func ComplexUpdates(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_postgresql_complex_updates", data)
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_postgresql_data_basic", data)
}
