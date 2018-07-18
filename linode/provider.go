package linode

import (
	"fmt"
	"net/http"

	"github.com/chiefy/linodego"
	"github.com/displague/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/version"
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

		DataSourcesMap: map[string]*schema.Resource{
			"linode_ipv6_pool":  dataSourceLinodeComputeIPv6Pool(),
			"linode_ipv6_range": dataSourceLinodeComputeIPv6Range(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_instance":            resourceLinodeInstance(),
			"linode_nodebalancer":        resourceLinodeNodeBalancer(),
			"linode_nodebalancer_config": resourceLinodeNodeBalancerConfig(),
			"linode_nodebalancer_node":   resourceLinodeNodeBalancerNode(),
			"linode_volume":              resourceLinodeVolume(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token, ok := d.Get("token").(string)
	if !ok {
		return nil, fmt.Errorf("The Linode API Token was not valid")
	}
	var httpTransport http.Transport
	transport := logging.NewTransport("Linode", &httpTransport)
	client := linodego.NewClient(&token, transport)

	projectURL := "https://www.terraform.io"
	userAgent := fmt.Sprintf("Terraform/%s (+%s)",
		version.String(), projectURL)

	client.SetUserAgent(userAgent)

	// Ping the API for an empty response to verify the configuration works
	_, err := client.ListTypes(linodego.NewListOptions(100, ""))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the Linode API because %s", err)
	}

	return client, nil
}
