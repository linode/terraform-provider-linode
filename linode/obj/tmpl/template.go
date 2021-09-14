package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	objectbucket "github.com/linode/terraform-provider-linode/linode/objbucket/tmpl"
	objectkey "github.com/linode/terraform-provider-linode/linode/objkey/tmpl"
)

type TemplateData struct {
	Bucket objectbucket.TemplateData
	Key    objectkey.TemplateData

	Content string
	Source  string
}

func Basic(t *testing.T, name, keyName, content string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_basic", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
		})
}

func Base64(t *testing.T, name, keyName, content string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_base64", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
		})
}

func Source(t *testing.T, name, keyName, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_source", TemplateData{
			Bucket: objectbucket.TemplateData{Label: name},
			Key:    objectkey.TemplateData{Label: keyName},
			Source: source,
		})
}

func Updates(t *testing.T, name, keyName, content string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_updates", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
		})
}
