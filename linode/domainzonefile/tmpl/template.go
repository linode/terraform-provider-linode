package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	domain "github.com/linode/terraform-provider-linode/linode/domain/tmpl"
)

type TemplateData struct {
	Domain domain.TemplateData

	Record string
}

func Basic(t *testing.T, domainRecord string) string {
	return acceptance.ExecuteTemplate(t,
		"domain_zonefile_basic", TemplateData{
			Domain: domain.TemplateData{Domain: domainRecord},
			Record: domainRecord,
		})
}
