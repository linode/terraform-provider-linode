package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	Region string
	Tags   []string
}

func Basic(t *testing.T, region string) string {
	return acceptance.ExecuteTemplate(t,
		"reserved_ip_basic", TemplateData{
			Region: region,
		})
}

func WithTags(t *testing.T, region string, tags []string) string {
	return acceptance.ExecuteTemplate(t,
		"reserved_ip_with_tags", TemplateData{
			Region: region,
			Tags:   tags,
		})
}
