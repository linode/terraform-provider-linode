package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	DataSourceName string
	ResourceName   string
	FirewallID     int
}

func Basic(t testing.TB, resourceName string, testFirewallID int) string {
	return acceptance.ExecuteTemplate(t, "linode_firewall_settings_basic", TemplateData{
		ResourceName: resourceName,
		FirewallID:   testFirewallID,
	})
}

func Data(t testing.TB, datasourceName string) string {
	return acceptance.ExecuteTemplate(t, "data_linode_firewall_settings", TemplateData{
		DataSourceName: datasourceName,
	})
}
