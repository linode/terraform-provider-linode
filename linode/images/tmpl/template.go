package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Image  string
	Region string
	Label  string
}

func DataBasic(t testing.TB, image, region, label string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_basic", TemplateData{Image: image, Region: region, Label: label})
}

func DataLatest(t testing.TB, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest", TemplateData{Image: image, Region: region})
}

func DataLatestEmpty(t testing.TB, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_latest_empty", TemplateData{Image: image, Region: region})
}

func DataOrder(t testing.TB, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_order", TemplateData{Image: image, Region: region})
}

func DataSubstring(t testing.TB, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_substring", TemplateData{Image: image, Region: region})
}

func DataClientFilter(t testing.TB, image, region string) string {
	return acceptance.ExecuteTemplate(t,
		"images_data_clientfilter", TemplateData{Image: image, Region: region})
}
