package linode

import (
	"fmt"

	golinode "github.com/chiefy/go-linode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_API_KEY", nil),
				Description: "The api key that allows you access to your linode account",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_linode": resourceLinodeLinode(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := golinode.NewClient(d.Get("key").(string), nil)

	_, err := client.ListKernels()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the linode api because %s", err)
	}

	return client, nil
}
