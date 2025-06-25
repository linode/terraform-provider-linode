package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	EndpointType string
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(
		t, "objendpoints_data_basic", TemplateData{},
	)
}

func DataFilter(t testing.TB, endpointType string) string {
	return acceptance.ExecuteTemplate(
		t, "objendpoints_data_filter", TemplateData{
			EndpointType: endpointType,
		},
	)
}
