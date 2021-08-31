package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	ID string
}

func DataBasic(t *testing.T, kernelID string) string {
	return acceptance.ExecuteTemplate(t,
		"kernel_data_basic", TemplateData{ID: kernelID})
}
