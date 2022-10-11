package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	PubKey string
	Region string
}

func Basic(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_basic", TemplateData{Label: volume, Region: region})
}

func Updates(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_updates", TemplateData{Label: volume, Region: region})
}

func Resized(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_resized", TemplateData{Label: volume, Region: region})
}

func Attached(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_attached", TemplateData{Label: volume, Region: region})
}

func ReAttached(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_reattached", TemplateData{Label: volume, Region: region})
}

func ClonedStep1(t *testing.T, volume, region, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_cloned_step1", TemplateData{
			Label:  volume,
			PubKey: pubKey,
			Region: region,
		})
}

func ClonedStep2(t *testing.T, volume, region, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_cloned_step2", TemplateData{
			Label:  volume,
			PubKey: pubKey,
			Region: region,
		})
}

func DataBasic(t *testing.T, volume, region string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_data_basic", TemplateData{Label: volume, Region: region})
}
