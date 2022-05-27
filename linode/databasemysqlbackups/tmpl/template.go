package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"

	"testing"
)

type TemplateData struct {
	Engine      string
	Label       string
	BackupLabel string
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_backups_data_basic", data)
}
