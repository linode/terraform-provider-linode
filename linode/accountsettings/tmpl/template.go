package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	LongviewSubscription string
	BackupsEnabled       bool
	NetworkHelper        bool
}

func Basic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_basic", nil)
}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_data_basic", nil)
}

func Updates(t *testing.T, longviewSubscription string, backupsEnabled, networkHelper bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_settings_updates", TemplateData{LongviewSubscription: longviewSubscription,
			BackupsEnabled: backupsEnabled, NetworkHelper: networkHelper})
}
