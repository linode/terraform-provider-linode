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

func DataLatest(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest", TemplateData{Image: image})
}

func DataLatestEmpty(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest_empty", TemplateData{Image: image})
}

func DataOrder(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_order", TemplateData{Image: image})
}

func DataSubstring(t *testing.T, image string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_substring", TemplateData{Image: image})
}
