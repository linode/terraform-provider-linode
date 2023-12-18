package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	SSHKey string
}

func Basic(t *testing.T, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_basic", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func Updates(t *testing.T, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_updates", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func DataBasic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_data_basic", TemplateData{Label: label})
}
