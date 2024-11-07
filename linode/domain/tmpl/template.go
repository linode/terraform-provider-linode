package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Domain string
}

func Basic(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_basic", TemplateData{Domain: domain})
}

func Updates(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_updates", TemplateData{Domain: domain})
}

func RoundedSec(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_rounded_sec", TemplateData{Domain: domain})
}

func ZeroSec(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_zero_sec", TemplateData{Domain: domain})
}

func IPS(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_ips", TemplateData{Domain: domain})
}

func IPSUpdates(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_ips_updates", TemplateData{Domain: domain})
}

func DataBasic(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_data_basic", TemplateData{Domain: domain})
}

func DataByID(t testing.TB, domain string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_data_byid", TemplateData{Domain: domain})
}
