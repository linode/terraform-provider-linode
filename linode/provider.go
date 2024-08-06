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
	"github.com/linode/terraform-provider-linode/v2/linode/databaseaccesscontrols"
	"github.com/linode/terraform-provider-linode/v2/linode/databasemysql"
	"github.com/linode/terraform-provider-linode/v2/linode/databasemysqlbackups"
	"github.com/linode/terraform-provider-linode/v2/linode/databasepostgresql"
	"github.com/linode/terraform-provider-linode/v2/linode/domain"
	"github.com/linode/terraform-provider-linode/v2/linode/domainrecord"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/instance"
	"github.com/linode/terraform-provider-linode/v2/linode/instanceconfig"
	"github.com/linode/terraform-provider-linode/v2/linode/lke"
	"github.com/linode/terraform-provider-linode/v2/linode/nbnode"
	"github.com/linode/terraform-provider-linode/v2/linode/obj"
	"github.com/linode/terraform-provider-linode/v2/linode/objbucket"
	"github.com/linode/terraform-provider-linode/v2/linode/user"
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The path to the Linode config file to use. (default `~/.config/linode`)",
			},
			"config_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Linode config profile to use. (default `default`)",
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

			"skip_implicit_reboots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, Linode Instances will not be rebooted on config and interface changes.",
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
			"obj_access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The access key to be used in linode_object_storage_bucket and linode_object_storage_object.",
			},
			"obj_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The secret key to be used in linode_object_storage_bucket and linode_object_storage_object.",
				Sensitive:   true,
			},
			"obj_use_temp_keys": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true, temporary object keys will be created implicitly at apply-time " +
					"for the linode_object_storage_object and linode_object_sorage_bucket resource.",
			},
			"obj_bucket_force_delete": {
				Type:     schema.TypeBool,
				Optional: true,
				Description: "If true, when deleting a linode_object_storage_bucket any objects " +
					"and versions will be force deleted.",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"linode_database_mysql_backups": databasemysqlbackups.DataSource(),
			"linode_instances":              instance.DataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"linode_database_access_controls": databaseaccesscontrols.Resource(),
			"linode_database_mysql":           databasemysql.Resource(),
			"linode_database_postgresql":      databasepostgresql.Resource(),
			"linode_domain":                   domain.Resource(),
			"linode_domain_record":            domainrecord.Resource(),
			"linode_instance":                 instance.Resource(),
			"linode_instance_config":          instanceconfig.Resource(),
			"linode_lke_cluster":              lke.Resource(),
			"linode_nodebalancer_node":        nbnode.Resource(),
			"linode_object_storage_bucket":    objbucket.Resource(),
			"linode_object_storage_object":    obj.Resource(),
			"linode_user":                     user.Resource(),
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
		eventPollMs, err := strconv.Atoi(os.Getenv("LINODE_EVENT_POLL_MS"))
		if err != nil {
			eventPollMs = 4000
		}
		config.EventPollMilliseconds = eventPollMs
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

	if v, ok := d.GetOk("obj_access_key"); ok {
		config.ObjAccessKey = v.(string)
	} else {
		config.ObjAccessKey = os.Getenv("LINODE_OBJ_ACCESS_KEY")
	}

	if v, ok := d.GetOk("obj_secret_key"); ok {
		config.ObjSecretKey = v.(string)
	} else {
		config.ObjSecretKey = os.Getenv("LINODE_OBJ_SECRET_KEY")
	}

	return nil
}

func providerConfigure(
	ctx context.Context, d *schema.ResourceData, terraformVersion string,
) (interface{}, diag.Diagnostics) {
	config := &helper.Config{
		SkipInstanceReadyPoll:  d.Get("skip_instance_ready_poll").(bool),
		SkipInstanceDeletePoll: d.Get("skip_instance_delete_poll").(bool),
		SkipImplicitReboots:    d.Get("skip_implicit_reboots").(bool),

		DisableInternalCache: d.Get("disable_internal_cache").(bool),

		MinRetryDelayMilliseconds: d.Get("min_retry_delay_ms").(int),
		MaxRetryDelayMilliseconds: d.Get("max_retry_delay_ms").(int),

		ObjUseTempKeys:       d.Get("obj_use_temp_keys").(bool),
		ObjBucketForceDelete: d.Get("obj_bucket_force_delete").(bool),
	}

	handleDefault(config, d)

	config.TerraformVersion = terraformVersion
	client, err := config.Client(ctx)
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
