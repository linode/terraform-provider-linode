package tmpl

import (
	"testing"

	objectbucket "github.com/linode/terraform-provider-linode/v2/linode/objbucket/tmpl"
	objectkey "github.com/linode/terraform-provider-linode/v2/linode/objkey/tmpl"

	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
)

type TemplateData struct {
	Bucket objectbucket.TemplateData
	Key    objectkey.TemplateData
}

func BasicDependency(t *testing.T, cluster, bucket, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_website_config_basic_dependency", TemplateData{
			Bucket: objectbucket.TemplateData{Label: bucket, Cluster: cluster},
			Key:    objectkey.TemplateData{Label: keyName},
		})
}

func UpdateDependency(t *testing.T, cluster, bucket, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_website_config_update_dependency", TemplateData{
			Bucket: objectbucket.TemplateData{Label: bucket, Cluster: cluster},
			Key:    objectkey.TemplateData{Label: keyName},
		})
}

func Basic(t *testing.T, cluster, bucket, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_website_config_basic", TemplateData{
			Bucket: objectbucket.TemplateData{Label: bucket, Cluster: cluster},
			Key:    objectkey.TemplateData{Label: keyName},
		})
}

func UpdatesBefore(t *testing.T, cluster, bucket, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_website_config_updates_before", TemplateData{
			Bucket: objectbucket.TemplateData{Label: bucket, Cluster: cluster},
			Key:    objectkey.TemplateData{Label: keyName},
		})
}

func UpdatesAfter(t *testing.T, cluster, bucket, keyName string) string {
	return acceptance.ExecuteTemplate(t,
		"object_website_config_updates_after", TemplateData{
			Bucket: objectbucket.TemplateData{Label: bucket, Cluster: cluster},
			Key:    objectkey.TemplateData{Label: keyName},
		})
}
