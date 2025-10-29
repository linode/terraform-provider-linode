package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	ShareGroupLabel string
	TokenLabel      string
}

func DataBasic(t testing.TB, shareGroupLabel, tokenLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"consumer_image_share_group_tokens_data_basic", TemplateData{
			ShareGroupLabel: shareGroupLabel,
			TokenLabel:      tokenLabel,
		})
}
