package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"

	"testing"
)

type TemplateData struct {
	Engine    string
	Label     string
	AllowedIP string
}

func Basic(t *testing.T, label, engine, ip string) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_allow_list_basic", TemplateData{
			Engine:    engine,
			Label:     label,
			AllowedIP: ip,
		})
}
