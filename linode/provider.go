package linode

import (
	"fmt"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_TOKEN", nil),
				Description: "The token that allows you access to your Linode account",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_linode": resourceLinodeLinode(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token, ok := d.Get("token").(string)
	if !ok {
		return nil, fmt.Errorf("The Linode API Token was not valid")
	}
	client := linodego.NewClient(&token, nil)

	// Ping the API for an empty response to verify the configuration works
	_, err := client.ListTypes(linodego.NewListOptions(100, ""))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the Linode API because %s", err)
	}

	return client, nil
}
