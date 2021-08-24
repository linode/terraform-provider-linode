package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/account"
	"github.com/linode/terraform-provider-linode/linode/balancer"
	"github.com/linode/terraform-provider-linode/linode/balancerconfig"
	"github.com/linode/terraform-provider-linode/linode/balancernode"
	"github.com/linode/terraform-provider-linode/linode/bucket"
	"github.com/linode/terraform-provider-linode/linode/domain"
	"github.com/linode/terraform-provider-linode/linode/domainrecord"
	"github.com/linode/terraform-provider-linode/linode/firewall"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"
	"github.com/linode/terraform-provider-linode/linode/images"
	"github.com/linode/terraform-provider-linode/linode/lke"
	"github.com/linode/terraform-provider-linode/linode/object"
	"github.com/linode/terraform-provider-linode/linode/objectcluster"
	"github.com/linode/terraform-provider-linode/linode/objectkey"
	"github.com/linode/terraform-provider-linode/linode/rdns"
	"github.com/linode/terraform-provider-linode/linode/token"
)

// Provider creates and manages the resources in a Linode configuration.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_TOKEN", nil),
				Description: "The token that allows you access to your Linode account",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LINODE_URL", nil),
				Description:  "The HTTP(S) API address of the Linode API to use.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"ua_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_UA_PREFIX", nil),
				Description: "An HTTP User-Agent Prefix to prepend in API requests.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_API_VERSION", nil),
				Description: "An HTTP User-Agent Prefix to prepend in API requests.",
			},

			"skip_instance_ready_poll": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip waiting for a linode_instance resource to be running.",
			},

			"skip_instance_delete_poll": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Skip waiting for a linode_instance resource to finish deleting.",
			},

			"min_retry_delay_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Minimum delay in milliseconds before retrying a request.",
			},
			"max_retry_delay_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum delay in milliseconds before retrying a request.",
			},

			"event_poll_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("LINODE_EVENT_POLL_MS", 300),
				Description: "The rate in milliseconds to poll for events.",
			},

			"lke_event_poll_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "The rate in milliseconds to poll for LKE events.",
			},

			"lke_node_ready_poll_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     500,
				Description: "The rate in milliseconds to poll for an LKE node to be ready.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"linode_account":                account.DataSource(),
			"linode_domain":                 domain.DataSource(),
			"linode_domain_record":          domainrecord.DataSource(),
			"linode_firewall":               firewall.DataSource(),
			"linode_image":                  image.DataSource(),
			"linode_images":                 images.DataSource(),
			"linode_instances":              dataSourceLinodeInstances(),
			"linode_instance_backups":       dataSourceLinodeInstanceBackups(),
			"linode_instance_type":          dataSourceLinodeInstanceType(),
			"linode_kernel":                 dataSourceLinodeKernel(),
			"linode_lke_cluster":            lke.DataSource(),
			"linode_networking_ip":          dataSourceLinodeNetworkingIP(),
			"linode_nodebalancer":           balancer.DataSource(),
			"linode_nodebalancer_node":      balancernode.DataSource(),
			"linode_nodebalancer_config":    balancerconfig.DataSource(),
			"linode_object_storage_cluster": objectcluster.DataSource(),
			"linode_profile":                dataSourceLinodeProfile(),
			"linode_region":                 dataSourceLinodeRegion(),
			"linode_sshkey":                 dataSourceLinodeSSHKey(),
			"linode_stackscript":            dataSourceLinodeStackscript(),
			"linode_user":                   dataSourceLinodeUser(),
			"linode_vlans":                  dataSourceLinodeVLANs(),
			"linode_volume":                 dataSourceLinodeVolume(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_domain":                domain.Resource(),
			"linode_domain_record":         domainrecord.Resource(),
			"linode_firewall":              firewall.Resource(),
			"linode_image":                 image.Resource(),
			"linode_instance":              resourceLinodeInstance(),
			"linode_instance_ip":           resourceLinodeInstanceIP(),
			"linode_lke_cluster":           lke.Resource(),
			"linode_nodebalancer":          balancer.Resource(),
			"linode_nodebalancer_node":     balancernode.Resource(),
			"linode_nodebalancer_config":   balancerconfig.Resource(),
			"linode_object_storage_key":    objectkey.Resource(),
			"linode_object_storage_bucket": bucket.Resource(),
			"linode_object_storage_object": object.Resource(),
			"linode_rdns":                  rdns.Resource(),
			"linode_sshkey":                resourceLinodeSSHKey(),
			"linode_stackscript":           resourceLinodeStackscript(),
			"linode_token":                 token.Resource(),
			"linode_user":                  resourceLinodeUser(),
			"linode_volume":                resourceLinodeVolume(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(ctx, d, terraformVersion)
	}
	return provider
}

func providerConfigure(
	ctx context.Context, d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	config := &helper.Config{
		AccessToken: d.Get("token").(string),
		APIURL:      d.Get("url").(string),
		APIVersion:  d.Get("api_version").(string),
		UAPrefix:    d.Get("ua_prefix").(string),

		SkipInstanceReadyPoll:  d.Get("skip_instance_ready_poll").(bool),
		SkipInstanceDeletePoll: d.Get("skip_instance_delete_poll").(bool),

		MinRetryDelayMilliseconds: d.Get("min_retry_delay_ms").(int),
		MaxRetryDelayMilliseconds: d.Get("max_retry_delay_ms").(int),

		EventPollMilliseconds:    d.Get("event_poll_ms").(int),
		LKEEventPollMilliseconds: d.Get("lke_event_poll_ms").(int),

		LKENodeReadyPollMilliseconds: d.Get("lke_node_ready_poll_ms").(int),
	}
	config.TerraformVersion = terraformVersion
	client := config.Client()

	// Ping the API for an empty response to verify the configuration works
	if _, err := client.ListTypes(ctx, linodego.NewListOptions(100, "")); err != nil {
		return nil, diag.Errorf("Error connecting to the Linode API: %s", err)
	}
	return &helper.ProviderMeta{
		Client: client,
		Config: config,
	}, nil
}
