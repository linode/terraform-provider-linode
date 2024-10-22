package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Id string
}

func DataBasic(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_basic", TemplateData{
			Id: id,
		})
}

func DataFilter(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_filter", TemplateData{
			Id: id,
		})
}

func DataFilterEmpty(t testing.TB, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_filter_empty", TemplateData{
			Id: id,
		})
}
