package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label string
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

func DataBasic(t *testing.T, volume string) string {
	return acceptance.ExecuteTemplate(t,
		"volume_data_basic", TemplateData{Label: volume})
}
