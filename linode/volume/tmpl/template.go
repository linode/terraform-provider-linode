package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	PubKey string
}

func Basic(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_basic", TemplateData{Label: volume})
}

func Updates(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_updates", TemplateData{Label: volume})
}

func Resized(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_resized", TemplateData{Label: volume})
}

func Attached(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_attached", TemplateData{Label: volume})
}

func ReAttached(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_reattached", TemplateData{Label: volume})
}

func ClonedStep1(t *testing.T, volume, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_cloned_step1", TemplateData{
			Label:  volume,
			PubKey: pubKey,
		})
}

func ClonedStep2(t *testing.T, volume, pubKey string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_cloned_step2", TemplateData{
			Label:  volume,
			PubKey: pubKey,
		})
}

func DataBasic(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_data_basic", TemplateData{Label: volume})
}
