package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/nb/tmpl"
	config "github.com/linode/terraform-provider-linode/v3/linode/nbconfig/tmpl"
)

type TemplateData struct {
	Label    string
	Instance InstanceTemplateData
	Config   config.TemplateData
}

type InstanceTemplateData struct {
	Label    string
	PubKey   string
	Region   string
	RootPass string
}

func Basic(t testing.TB, nodebalancer, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_basic",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:    nodebalancer,
				PubKey:   acceptance.PublicKeyMaterial,
				Region:   region,
				RootPass: rootPass,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label:  nodebalancer,
					Region: region,
				},
			},
		})
}

func Updates(t testing.TB, nodebalancer, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_updates",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:    nodebalancer,
				PubKey:   acceptance.PublicKeyMaterial,
				Region:   region,
				RootPass: rootPass,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label:  nodebalancer,
					Region: region,
				},
			},
		})
}

func VPC(t testing.TB, nodebalancer, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_vpc",
		TemplateData{
			Label: nodebalancer,
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label:  nodebalancer,
					Region: region,
				},
			},
		})
}

func DataBasic(t testing.TB, nodebalancer, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_data_basic",
		TemplateData{
			Label: nodebalancer,
			Instance: InstanceTemplateData{
				Label:    nodebalancer,
				PubKey:   acceptance.PublicKeyMaterial,
				Region:   region,
				RootPass: rootPass,
			},
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label:  nodebalancer,
					Region: region,
				},
			},
		})
}

func DataVPC(t testing.TB, nodebalancer, region string, rootPass string) string {
	return acceptance.ExecuteTemplate(t,
		"nodebalancer_node_data_vpc",
		TemplateData{
			Label: nodebalancer,
			Config: config.TemplateData{
				NodeBalancer: tmpl.TemplateData{
					Label:  nodebalancer,
					Region: region,
				},
			},
		})
}
