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
	Role        string
}

func Update(t testing.TB, label, region, username, email, role string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"iam_user_update", TemplateData{
			VolumeLabel: label,
			Region:      region,
			Username:    username,
			Email:       email,
			Restricted:  restricted,
			Role:        role,
		})
}

func UpdateAccount(t testing.TB, label, region, username, email, role string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"iam_user_update_account", TemplateData{
			VolumeLabel: label,
			Region:      region,
			Username:    username,
			Email:       email,
			Restricted:  restricted,
			Role:        role,
		})
}

func DataBasic(t testing.TB, username, email string, restricted bool) string {
	return acceptance.ExecuteTemplate(t,
		"iam_user_data_basic", TemplateData{
			Username:   username,
			Email:      email,
			Restricted: restricted,
		})
}
