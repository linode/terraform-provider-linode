package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func DataList(t *testing.T) string {
	return acceptance.ExecuteTemplate(t, "networking_ip_data_list", nil)
}

func DataFilterReserved(t *testing.T) string {
	return acceptance.ExecuteTemplate(t, "networking_ip_data_filtered", nil)
}
