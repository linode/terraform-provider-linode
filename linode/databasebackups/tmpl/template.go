package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Engine      string
	Label       string
	BackupLabel string
	Region      string
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_backups_data_basic", data)
}
