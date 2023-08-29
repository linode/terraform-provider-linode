package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Id string
}

func DataBasic(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_basic", TemplateData{
			Id: id,
		})
}

func DataFilter(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_filter", TemplateData{
			Id: id,
		})
}

func DataFilterEmpty(t *testing.T, id string) string {
	return acceptance.ExecuteTemplate(t,
		"kernels_data_filter_empty", TemplateData{
			Id: id,
		})
}
