package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Region string
}

func DataBasic(t *testing.T, region string) string {
	return acceptance.ExecuteTemplate(t,
		"account_availability_data_basic", TemplateData{
			Region: region,
		})
}
