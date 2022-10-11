package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
)

type TemplateData struct {
	Label  string
	Region string
}

func SingleNode(t *testing.T, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_shared_ips_single_node", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}

func DualNode(t *testing.T, instanceLabel, region string) string {
	return acceptance.ExecuteTemplate(t,
		"instance_shared_ips_dual_node", TemplateData{
			Label:  instanceLabel,
			Region: region,
		})
}
