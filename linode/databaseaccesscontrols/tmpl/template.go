package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Engine    string
	Label     string
	Region    string
	AllowedIP string
}

func MySQL(t *testing.T, label, engine, ip, region string) string {
	return acceptance.ExecuteTemplate(t,
		"database_access_controls_mysql", TemplateData{
			Engine:    engine,
			Label:     label,
			Region:    region,
			AllowedIP: ip,
		})
}

func MongoDB(t *testing.T, label, engine, ip, region string) string {
	return acceptance.ExecuteTemplate(t,
		"database_access_controls_mongodb", TemplateData{
			Engine:    engine,
			Label:     label,
			AllowedIP: ip,
			Region:    region,
		})
}

func PostgreSQL(t *testing.T, label, engine, ip, region string) string {
	return acceptance.ExecuteTemplate(t,
		"database_access_controls_postgresql", TemplateData{
			Engine:    engine,
			Label:     label,
			AllowedIP: ip,
			Region:    region,
		})
}
