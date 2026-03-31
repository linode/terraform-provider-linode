package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
)

type TemplateData struct {
	VolumeLabel string
	Region      string
	Username    string
	Email       string
	Restricted  bool
}

func Update(t testing.TB, label, region, username, email string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"iam_user_update", TemplateData{
			VolumeLabel: label,
			Region:      region,
			Username:    username,
			Email:       email,
			Restricted:  restricted,
		})
}

func DataBasic(t testing.TB) string {
	return acceptance.ExecuteTemplate(t,
		"user_data_basic", nil)
}
