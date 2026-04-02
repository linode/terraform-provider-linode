package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	BackupsEnabled          bool
	NetworkHelper           bool
	InterfacesForNewLinodes string
	MaintenancePolicy       string
}

func Basic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_basic", nil)
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_data_basic", nil)
}

func Updates(t testing.TB, interfacesForNewLinodes string, backupsEnabled, networkHelper bool, maintenancePolicy string) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_updates", TemplateData{
			BackupsEnabled:          backupsEnabled,
			NetworkHelper:           networkHelper,
			InterfacesForNewLinodes: interfacesForNewLinodes,
			MaintenancePolicy:       maintenancePolicy,
		})
}
