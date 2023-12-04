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

func Basic(t *testing.T, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
		})
}

func Updates(t *testing.T, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_updates", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
		})
}

func ANoName(t *testing.T, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_a_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
		})
}

func AAAANoName(t *testing.T, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_aaaa_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
		})
}

func CAANoName(t *testing.T, d string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_caa_noname", TemplateData{
			Domain: domain.TemplateData{Domain: d},
			Record: d,
		})
}

func TTL(t *testing.T, domainRecord string, ttlSec int) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_ttl", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord + ".example"},
			Record: domainRecord,
			TTL:    ttlSec,
		})
}

func SRV(t *testing.T, domainName string, target string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_srv", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
			Target: target,
		})
}

func DataBasic(t *testing.T, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataID(t *testing.T, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_id", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataSRV(t *testing.T, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_srv", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}

func DataCAA(t *testing.T, domainName string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_record_data_caa", TemplateData{
			Domain: domain.TemplateData{Domain: domainName},
		})
}
