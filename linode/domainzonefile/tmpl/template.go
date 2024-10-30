package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	domain "github.com/linode/terraform-provider-linode/v2/linode/domain/tmpl"
)

type TemplateData struct {
	Domain domain.TemplateData

	Record string
}

func Basic(t testing.TB, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_zonefile_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord},
			Record: domainRecord,
		})
}
