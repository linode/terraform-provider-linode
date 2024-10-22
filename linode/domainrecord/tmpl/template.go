package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	domain "github.com/linode/terraform-provider-linode/v2/linode/domain/tmpl"
)

type TemplateData struct {
	Domain domain.TemplateData

	Record string
	Target string

	TTL int
}

func Basic(t testing.TB, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
		})
}

func Updates(t testing.TB, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_updates", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
		})
}

func ANoName(t testing.TB, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_a_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
		})
}

func AAAANoName(t testing.TB, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_aaaa_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
		})
}

func CAANoName(t testing.TB, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_caa_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
			Record: d,
		})
}

func TTL(t testing.TB, domainRecord string, ttlSec int) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_ttl", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
			TTL:    ttlSec,
		})
}

func SRV(t testing.TB, domainName string, target string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_srv", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
			Target: target,
		})
}

func WithDomain(t testing.TB, domainName, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_with_domain", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
			Record: domainRecord,
		})
}

func DataBasic(t testing.TB, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataID(t testing.TB, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_id", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataSRV(t testing.TB, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_srv", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataCAA(t testing.TB, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_caa", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}
