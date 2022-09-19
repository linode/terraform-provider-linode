package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	objectbucket "github.com/linode/terraform-provider-linode/linode/objbucket/tmpl"
	objectkey "github.com/linode/terraform-provider-linode/linode/objkey/tmpl"
)

type TemplateData struct {
	Bucket  objectbucket.TemplateData
	Key     objectkey.TemplateData
	Cluster string

	Content string
	Source  string
}

func Basic(t *testing.T, name, cluster, keyName, content, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_basic", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Source:  source,
			Cluster: cluster,
		})
}

func Updates(t *testing.T, name, cluster, keyName, content, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_updates", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Source:  source,
			Cluster: cluster,
		})
}
