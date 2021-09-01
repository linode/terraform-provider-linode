package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	objkey "github.com/linode/terraform-provider-linode/linode/objectkey/tmpl"
)

type TemplateData struct {
	Key objkey.TemplateData

	Label       string
	ACL         string
	CORSEnabled bool
	Versioning  bool

	Cert    string
	PrivKey string
}

func Basic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_basic", TemplateData{Label: label})
}

func Updates(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_updates", TemplateData{Label: label})
}

func Access(t *testing.T, label, acl string, cors bool) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_access", TemplateData{
			Label:       label,
			ACL:         acl,
			CORSEnabled: cors,
		})
}

func Cert(t *testing.T, label, cert, privKey string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_cert", TemplateData{
			Label:   label,
			Cert:    cert,
			PrivKey: privKey,
		})
}

func Versioning(t *testing.T, label, keyName string, versioning bool) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_versioning", TemplateData{
			Key:        objkey.TemplateData{Label: keyName},
			Label:      label,
			Versioning: versioning,
		})
}

func LifeCycle(t *testing.T, label, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle", TemplateData{
			Key:   objkey.TemplateData{Label: keyName},
			Label: label,
		})
}

func LifeCycleUpdates(t *testing.T, label, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_lifecycle_updates", TemplateData{
			Key:   objkey.TemplateData{Label: keyName},
			Label: label,
		})
}

func DataBasic(t *testing.T, label string) string {
	return acceptance.ExecuteTemplate(t,
		"object_bucket_data_basic", TemplateData{
			Label: label,
		})
}
