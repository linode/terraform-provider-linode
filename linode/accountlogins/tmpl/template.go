package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Username   string
	IP         string
	Restricted bool
	Status     string
}

func DataBasic(t *testing.T) string {
	return acceptance.ExecuteTemplate(t,
		"account_logins_data_basic", nil)
}

func DataFilterRestricted(t *testing.T, username, ip, status string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_logins_data_filter_by_restricted", TemplateData{
			Username:   username,
			IP:         ip,
			Restricted: restricted,
			Status:     status,
		})
}

func DataFilterUsername(t *testing.T, username, ip, status string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_logins_data_filter_by_username", TemplateData{
			Username:   username,
			IP:         ip,
			Restricted: restricted,
			Status:     status,
		})
}

func DataFilterIP(t *testing.T, username, ip, status string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_logins_data_filter_by_ip", TemplateData{
			Username:   username,
			IP:         ip,
			Restricted: restricted,
			Status:     status,
		})
}

func DataFilterStatus(t *testing.T, username, ip, status string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"account_logins_data_filter_by_status", TemplateData{
			Username:   username,
			IP:         ip,
			Restricted: restricted,
			Status:     status,
		})
}
