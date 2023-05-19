package linode

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/accountlogin"
	"github.com/linode/terraform-provider-linode/linode/accountlogins"
	"github.com/linode/terraform-provider-linode/linode/accountsettings"
	"github.com/linode/terraform-provider-linode/linode/databaseaccesscontrols"
	"github.com/linode/terraform-provider-linode/linode/databasebackups"
	"github.com/linode/terraform-provider-linode/linode/databaseengines"
	"github.com/linode/terraform-provider-linode/linode/databasemysql"
	"github.com/linode/terraform-provider-linode/linode/databasemysqlbackups"
	"github.com/linode/terraform-provider-linode/linode/databasepostgresql"
	"github.com/linode/terraform-provider-linode/linode/databases"
	"github.com/linode/terraform-provider-linode/linode/domain"
	"github.com/linode/terraform-provider-linode/linode/domainrecord"
	"github.com/linode/terraform-provider-linode/linode/domainzonefile"
	"github.com/linode/terraform-provider-linode/linode/firewall"
	"github.com/linode/terraform-provider-linode/linode/firewalldevice"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/image"
	"github.com/linode/terraform-provider-linode/linode/images"
	"github.com/linode/terraform-provider-linode/linode/instance"
	"github.com/linode/terraform-provider-linode/linode/instanceconfig"
	"github.com/linode/terraform-provider-linode/linode/instancedisk"
	"github.com/linode/terraform-provider-linode/linode/instanceip"
	"github.com/linode/terraform-provider-linode/linode/instancenetworking"
	"github.com/linode/terraform-provider-linode/linode/instancesharedips"
	"github.com/linode/terraform-provider-linode/linode/instancetype"
	"github.com/linode/terraform-provider-linode/linode/instancetypes"
	"github.com/linode/terraform-provider-linode/linode/ipv6range"
	"github.com/linode/terraform-provider-linode/linode/lke"
	"github.com/linode/terraform-provider-linode/linode/nb"
	"github.com/linode/terraform-provider-linode/linode/nbconfig"
	"github.com/linode/terraform-provider-linode/linode/nbnode"
	"github.com/linode/terraform-provider-linode/linode/obj"
	"github.com/linode/terraform-provider-linode/linode/objbucket"
	"github.com/linode/terraform-provider-linode/linode/objcluster"
	"github.com/linode/terraform-provider-linode/linode/objkey"
	"github.com/linode/terraform-provider-linode/linode/rdns"
	"github.com/linode/terraform-provider-linode/linode/region"
	"github.com/linode/terraform-provider-linode/linode/regions"
	"github.com/linode/terraform-provider-linode/linode/sshkey"
	"github.com/linode/terraform-provider-linode/linode/stackscripts"
	"github.com/linode/terraform-provider-linode/linode/user"
	"github.com/linode/terraform-provider-linode/linode/vlan"
	"github.com/linode/terraform-provider-linode/linode/volume"
)

// Provider creates and manages the resources in a Linode configuration.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The token that allows you access to your Linode account",
			},
			"config_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The HTTP(S) API address of the Linode API to use.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"ua_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An HTTP User-Agent Prefix to prepend in API requests.",
			},
			"api_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The version of Linode API.",
			},

			"skip_instance_ready_poll": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip waiting for a linode_instance resource to be running.",
			},

			"skip_instance_delete_poll": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Skip waiting for a linode_instance resource to finish deleting.",
			},

			"disable_internal_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disable the internal caching system that backs certain Linode API requests.",
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
				Description: "The rate in milliseconds to poll for events.",
			},
			"lke_event_poll_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The rate in milliseconds to poll for LKE events.",
			},

			"lke_node_ready_poll_ms": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The rate in milliseconds to poll for an LKE node to be ready.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"linode_account_login":          accountlogin.DataSource(),
			"linode_account_logins":         accountlogins.DataSource(),
			"linode_account_settings":       accountsettings.DataSource(),
			"linode_database_backups":       databasebackups.DataSource(),
			"linode_database_engines":       databaseengines.DataSource(),
			"linode_database_mysql":         databasemysql.DataSource(),
			"linode_database_postgresql":    databasepostgresql.DataSource(),
			"linode_database_mysql_backups": databasemysqlbackups.DataSource(),
			"linode_databases":              databases.DataSource(),
			"linode_domain":                 domain.DataSource(),
			"linode_domain_zonefile":        domainzonefile.DataSource(),
			"linode_firewall":               firewall.DataSource(),
			"linode_image":                  image.DataSource(),
			"linode_images":                 images.DataSource(),
			"linode_instances":              instance.DataSource(),
			"linode_instance_type":          instancetype.DataSource(),
			"linode_instance_types":         instancetypes.DataSource(),
			"linode_instance_networking":    instancenetworking.DataSource(),
			"linode_ipv6_range":             ipv6range.DataSource(),
			"linode_lke_cluster":            lke.DataSource(),
			"linode_nodebalancer":           nb.DataSource(),
			"linode_nodebalancer_node":      nbnode.DataSource(),
			"linode_nodebalancer_config":    nbconfig.DataSource(),
			"linode_object_storage_bucket":  objbucket.DataSource(),
			"linode_object_storage_cluster": objcluster.DataSource(),
			"linode_region":                 region.DataSource(),
			"linode_regions":                regions.DataSource(),
			"linode_sshkey":                 sshkey.DataSource(),
			"linode_stackscripts":           stackscripts.DataSource(),
			"linode_user":                   user.DataSource(),
			"linode_vlans":                  vlan.DataSource(),
			"linode_volume":                 volume.DataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_account_settings":         accountsettings.Resource(),
			"linode_database_access_controls": databaseaccesscontrols.Resource(),
			"linode_database_mysql":           databasemysql.Resource(),
			"linode_database_postgresql":      databasepostgresql.Resource(),
			"linode_domain":                   domain.Resource(),
			"linode_domain_record":            domainrecord.Resource(),
			"linode_firewall":                 firewall.Resource(),
			"linode_firewall_device":          firewalldevice.Resource(),
			"linode_image":                    image.Resource(),
			"linode_instance":                 instance.Resource(),
			"linode_instance_config":          instanceconfig.Resource(),
			"linode_instance_disk":            instancedisk.Resource(),
			"linode_instance_ip":              instanceip.Resource(),
			"linode_instance_shared_ips":      instancesharedips.Resource(),
			"linode_ipv6_range":               ipv6range.Resource(),
			"linode_lke_cluster":              lke.Resource(),
			"linode_nodebalancer":             nb.Resource(),
			"linode_nodebalancer_node":        nbnode.Resource(),
			"linode_nodebalancer_config":      nbconfig.Resource(),
			"linode_object_storage_key":       objkey.Resource(),
			"linode_object_storage_bucket":    objbucket.Resource(),
			"linode_object_storage_object":    obj.Resource(),
			"linode_rdns":                     rdns.Resource(),
			"linode_sshkey":                   sshkey.Resource(),
			"linode_user":                     user.Resource(),
			"linode_volume":                   volume.Resource(),
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

func handleDefault(config *helper.Config, d *schema.ResourceData) diag.Diagnostics {
	if v, ok := d.GetOk("token"); ok {
		config.AccessToken = v.(string)
	} else {
		config.AccessToken = os.Getenv("LINODE_TOKEN")
	}

	if v, ok := d.GetOk("api_version"); ok {
		config.APIVersion = v.(string)
	} else {
		config.APIVersion = os.Getenv("LINODE_API_VERSION")
	}

	if v, ok := d.GetOk("config_path"); ok {
		config.ConfigPath = v.(string)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return diag.Errorf(
				"Failed to get user home directory: %s",
				err.Error(),
			)
		}
		config.ConfigPath = fmt.Sprintf("%s/.config/linode", homeDir)
	}

	if v, ok := d.GetOk("config_profile"); ok {
		config.ConfigProfile = v.(string)
	} else {
		config.ConfigProfile = "default"
	}

	if v, ok := d.GetOk("url"); ok {
		config.APIURL = v.(string)
	} else {
		config.APIURL = os.Getenv("LINODE_URL")
	}

	if v, ok := d.GetOk("ua_prefix"); ok {
		config.UAPrefix = v.(string)
	} else {
		config.UAPrefix = os.Getenv("LINODE_UA_PREFIX")
	}

	if v, ok := d.GetOk("event_poll_ms"); ok {
		config.EventPollMilliseconds = v.(int)
	} else {
		eventPollMs, err := strconv.ParseInt(os.Getenv("LINODE_EVENT_POLL_MS"), 10, 64)
		if err != nil {
			eventPollMs = 4000
		}
		config.EventPollMilliseconds = int(eventPollMs)
	}

	if v, ok := d.GetOk("lke_event_poll_ms"); ok {
		config.LKEEventPollMilliseconds = v.(int)
	} else {
		config.LKEEventPollMilliseconds = 3000
	}

	if v, ok := d.GetOk("lke_node_ready_poll_ms"); ok {
		config.LKENodeReadyPollMilliseconds = v.(int)
	} else {
		config.LKENodeReadyPollMilliseconds = 3000
	}

	return nil
}

func providerConfigure(
	ctx context.Context, d *schema.ResourceData, terraformVersion string,
) (interface{}, diag.Diagnostics) {
	config := &helper.Config{
		SkipInstanceReadyPoll:  d.Get("skip_instance_ready_poll").(bool),
		SkipInstanceDeletePoll: d.Get("skip_instance_delete_poll").(bool),

		DisableInternalCache: d.Get("disable_internal_cache").(bool),

		MinRetryDelayMilliseconds: d.Get("min_retry_delay_ms").(int),
		MaxRetryDelayMilliseconds: d.Get("max_retry_delay_ms").(int),
	}

	handleDefault(config, d)

	config.TerraformVersion = terraformVersion
	client, err := config.Client()
	if err != nil {
		return nil, diag.Errorf("failed to initialize client: %s", err)
	}

	// Ping the API for an empty response to verify the configuration works
	if _, err := client.ListTypes(ctx, linodego.NewListOptions(100, "")); err != nil {
		return nil, diag.Errorf("Error connecting to the Linode API: %s", err)
	}
	return &helper.ProviderMeta{
		Client: *client,
		Config: config,
	}, nil
}
