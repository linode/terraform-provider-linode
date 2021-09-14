package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/nb/tmpl"
	config "github.com/linode/terraform-provider-linode/linode/nbconfig/tmpl"
)

type TemplateData struct {
	Label    string
	Instance InstanceTemplateData
	Config   config.TemplateData
}

type InstanceTemplateData struct {
	Label  string
	PubKey string
}

func Basic(t *testing.T, nodebalancer string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_basic",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:  nodebalancer,
				PubKey: acceptance.PublicKeyMaterial,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label: nodebalancer,
				},
			},
		})
}

func Updates(t *testing.T, nodebalancer string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_updates",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:  nodebalancer,
				PubKey: acceptance.PublicKeyMaterial,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label: nodebalancer,
				},
			},
		})
}

func DataBasic(t *testing.T, nodebalancer string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_data_basic",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:  nodebalancer,
				PubKey: acceptance.PublicKeyMaterial,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label: nodebalancer,
				},
			},
		})
}
