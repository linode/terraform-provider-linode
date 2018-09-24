package linode

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/logging"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/version"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

// Provider creates and manages the resources in a Linode configuration.
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
			"linode_instance_type": dataSourceLinodeInstanceType(),
			"linode_region":        dataSourceLinodeRegion(),
			"linode_image":         dataSourceLinodeImage(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_instance":            resourceLinodeInstance(),
			"linode_domain":              resourceLinodeDomain(),
			"linode_domain_record":       resourceLinodeDomainRecord(),
			"linode_nodebalancer":        resourceLinodeNodeBalancer(),
			"linode_nodebalancer_config": resourceLinodeNodeBalancerConfig(),
			"linode_nodebalancer_node":   resourceLinodeNodeBalancerNode(),
			"linode_volume":              resourceLinodeVolume(),
			"linode_stackscript":         resourceLinodeStackscript(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	token, ok := d.Get("token").(string)
	if !ok {
		return nil, fmt.Errorf("The Linode API Token was not valid")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	oauthTransport := &oauth2.Transport{
		Source: tokenSource,
	}
	loggingTransport := logging.NewTransport("Linode", oauthTransport)
	oauth2Client := &http.Client{
		Transport: loggingTransport,
	}

	client := linodego.NewClient(oauth2Client)

	projectURL := "https://www.terraform.io"
	userAgent := fmt.Sprintf("Terraform/%s (+%s) linodego/%s",
		version.String(), projectURL, linodego.Version)

	client.SetUserAgent(userAgent)

	// Ping the API for an empty response to verify the configuration works
	_, err := client.ListTypes(context.Background(), linodego.NewListOptions(100, ""))
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the Linode API: %s", err)
	}

	return client, nil
}
