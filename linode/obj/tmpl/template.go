package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	objectbucket "github.com/linode/terraform-provider-linode/v2/linode/objbucket/tmpl"
	objectkey "github.com/linode/terraform-provider-linode/v2/linode/objkey/tmpl"
)

type TemplateData struct {
	Bucket  objectbucket.TemplateData
	Key     objectkey.TemplateData
	Cluster string
	Region  string

	Content string
	Source  string
}

func BasicWithCluster(t *testing.T, name, cluster, keyName, content, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_basic", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name, Cluster: cluster},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Source:  source,
			Cluster: cluster,
		})
}

func Basic(t *testing.T, name, region, keyName, content, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_basic", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name, Region: region},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Source:  source,
			Region:  region,
		})
}

func Updates(t *testing.T, name, region, keyName, content, source string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_updates", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name, Region: region},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Source:  source,
			Region:  region,
		})
}

func CredsConfiged(t *testing.T, name, region, keyName, content string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_creds_configed", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name, Region: region},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Region:  region,
		})
}

func TempKeys(t *testing.T, name, region, keyName, content string) string {
	return acceptance.ExecuteTemplate(t,
		"object_object_temp_keys", TemplateData{
			Bucket:  objectbucket.TemplateData{Label: name, Region: region},
			Key:     objectkey.TemplateData{Label: keyName},
			Content: content,
			Region:  region,
		})
}
