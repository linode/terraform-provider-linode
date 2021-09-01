package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Image string
}

func DataBasic(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_basic", TemplateData{Image: image})
}
