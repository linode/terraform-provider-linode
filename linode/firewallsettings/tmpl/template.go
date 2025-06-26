package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	DataSourceName string
}

func Data(t testing.TB, datasourceName string) string {
	return acceptance.ExecuteTemplate(t, "data_linode_firewall_settings", TemplateData{
		DataSourceName: datasourceName,
	})
}
