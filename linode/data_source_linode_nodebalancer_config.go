package linode

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLinodeNodeBalancerConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLinodeNodeBalancerConfigRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer config.",
				Required:    true,
			},
			"nodebalancer_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer to access.",
				Required:    true,
			},
			"protocol": {
				Type: schema.TypeString,
				Description: "The protocol this port is configured to serve. If this is set to https you must include " +
					"an ssl_cert and an ssl_key.",
				Computed: true,
			},
			"proxy_protocol": {
				Type: schema.TypeString,
				Description: "The version of ProxyProtocol to use for the underlying NodeBalancer. This requires " +
					"protocol to be `tcp`. Valid values are `none`, `v1`, and `v2`.",
				Computed: true,
			},
			"port": {
				Type: schema.TypeInt,
				Description: "The TCP port this Config is for. These values must be unique across configs on a single " +
					"NodeBalancer (you can't have two configs for port 80, for example). While some ports imply some " +
					"protocols, no enforcement is done and you may configure your NodeBalancer however is useful to you. " +
					"For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have a " +
					"NodeBalancer listening on port 443.",
				Computed: true,
			},
			"check_interval": {
				Type:        schema.TypeInt,
				Description: "How often, in seconds, to check that backends are up and serving requests.",
				Computed:    true,
			},
			"check_timeout": {
				Type:        schema.TypeInt,
				Description: "How long, in seconds, to wait for a check attempt before considering it failed. (1-30)",
				Computed:    true,
			},
			"check_attempts": {
				Type:        schema.TypeInt,
				Description: "How many times to attempt a check before considering a backend to be down. (1-30)",
				Computed:    true,
			},
			"algorithm": {
				Type: schema.TypeString,
				Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, " +
					"leastconn, source",
				Computed: true,
			},
			"stickiness": {
				Type:        schema.TypeString,
				Description: "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
				Computed:    true,
			},
			"check": {
				Type: schema.TypeString,
				Description: "The type of check to perform against backends to ensure they are serving requests. " +
					"This is used to determine if backends are up or down. If none no check is performed. " +
					"connection requires only a connection to the backend to succeed. http and http_body rely on the " +
					"backend serving HTTP, and that the response returned matches what is expected.",
				Computed: true,
			},
			"check_path": {
				Type: schema.TypeString,
				Description: "The URL path to check on each backend. If the backend does not respond to this request " +
					"it is considered to be down.",
				Computed: true,
			},
			"check_body": {
				Type: schema.TypeString,
				Description: "This value must be present in the response body of the check in order for it to pass. " +
					"If this value is not present in the response body of a check request, the backend is considered to be down",
				Computed: true,
			},
			"check_passive": {
				Type: schema.TypeBool,
				Description: "If true, any response from this backend with a 5xx status code will be enough for it to " +
					"be considered unhealthy and taken out of rotation.",
				Computed: true,
			},
			"cipher_suite": {
				Type: schema.TypeString,
				Description: "What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is " +
					"considered insecure and should only be used if necessary.",
				Computed: true,
			},
			"ssl_commonname": {
				Type: schema.TypeString,
				Description: "The read-only common name automatically derived from the SSL certificate assigned to " +
					"this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate " +
					"is assigned to your NodeBalancerConfig.",
				Computed: true,
			},
			"ssl_fingerprint": {
				Type: schema.TypeString,
				Description: "The read-only fingerprint automatically derived from the SSL certificate assigned to " +
					"this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate " +
					"is assigned to your NodeBalancerConfig.",
				Computed: true,
			},
			"node_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     resourceLinodeNodeBalancerConfigNodeStatus(),
			},
		},
	}
}

func datasourceLinodeNodeBalancerConfigRead(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	id := d.Get("id").(int)
	nodebalancerID := d.Get("nodebalancer_id").(int)

	config, err := client.GetNodeBalancerConfig(context.Background(), nodebalancerID, id)
	if err != nil {
		return diag.Errorf("failed to get nodebalancer config %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(config.ID))
	d.Set("nodebalancer_id", config.NodeBalancerID)
	d.Set("algorithm", config.Algorithm)
	d.Set("stickiness", config.Stickiness)
	d.Set("check", config.Check)
	d.Set("check_attempts", config.CheckAttempts)
	d.Set("check_body", config.CheckBody)
	d.Set("check_interval", config.CheckInterval)
	d.Set("check_timeout", config.CheckTimeout)
	d.Set("check_passive", config.CheckPassive)
	d.Set("check_path", config.CheckPath)
	d.Set("cipher_suite", config.CipherSuite)
	d.Set("port", config.Port)
	d.Set("protocol", config.Protocol)
	d.Set("proxy_protocol", config.ProxyProtocol)
	d.Set("ssl_fingerprint", config.SSLFingerprint)
	d.Set("ssl_commonname", config.SSLCommonName)
	d.Set("node_status", []map[string]interface{}{{
		"up":   config.NodesStatus.Up,
		"down": config.NodesStatus.Down,
	}})

	return nil
}
