package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Region  string
	Address string
}

func DataBasic(t *testing.T, region string) string {
	return acceptance.ExecuteTemplate(t,
		"reserved_ip_data_basic", TemplateData{
			Region: region,
		})
}

// ReserveIP generates the Terraform configuration for reserving an IP address
func ReserveIP(t *testing.T, name, region string) string {
	return acceptance.ExecuteTemplate(t,
		"reserved_ip_basic", TemplateData{
			Region: region,
		})
}
