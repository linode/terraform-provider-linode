package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
)

type TemplateData struct {
	Label  string
	SSHKey string
	ID     string
}

func Basic(t testing.TB, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_basic", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func Updates(t testing.TB, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_updates", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func DataBasic(t testing.TB, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_data_basic", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func DataByID(t testing.TB, label, sshKey string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_data_byid", TemplateData{
			Label:  label,
			SSHKey: sshKey,
		})
}

func DataMissing(t testing.TB, label string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_data_missing", TemplateData{Label: label})
}

func DataMissingByID(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"sshkey_data_missing_by_id", TemplateData{ID: id})
}
