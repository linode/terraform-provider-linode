package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Image  string
	Region string
}

func DataBasic(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_basic", TemplateData{Image: image, Region: region})
}

func DataLatest(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest", TemplateData{Image: image, Region: region})
}

func DataLatestEmpty(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest_empty", TemplateData{Image: image, Region: region})
}

func DataOrder(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_order", TemplateData{Image: image, Region: region})
}

func DataSubstring(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_substring", TemplateData{Image: image, Region: region})
}

func DataClientFilter(t *testing.T, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_clientfilter", TemplateData{Image: image, Region: region})
}
