package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeNodeBalancerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeNodeBalancerConfigCreate,
		Read:   resourceLinodeNodeBalancerConfigRead,
		Update: resourceLinodeNodeBalancerConfigUpdate,
		Delete: resourceLinodeNodeBalancerConfigDelete,
		Exists: resourceLinodeNodeBalancerConfigExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"nodebalancer_id": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The ID of the NodeBalancer to access.",
				Required:    true,
				ForceNew:    true,
			},
			"protocol": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The protocol this port is configured to serve. If this is set to https you must include an ssl_cert and an ssl_key.",
				Optional:    true,
				Default:     linodego.ProtocolHTTP,
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The TCP port this Config is for. These values must be unique across configs on a single NodeBalancer (you can't have two configs for port 80, for example). While some ports imply some protocols, no enforcement is done and you may configure your NodeBalancer however is useful to you. For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have a NodeBalancer listening on port 443.",
				Optional:    true,
				Default:     80,
			},
			"check_interval": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "How often, in seconds, to check that backends are up and serving requests.",
				Optional:    true,
				Computed:    true,
			},
			"check_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "How long, in seconds, to wait for a check attempt before considering it failed. (1-30)",
				Optional:    true,
				Computed:    true,
			},
			"check_attempts": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "How many times to attempt a check before considering a backend to be down. (1-30)",
				Optional:    true,
				Computed:    true,
			},
			"algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, leastconn, source",
				Optional:    true,
				Computed:    true,
			},
			"stickiness": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
				Optional:    true,
				Computed:    true,
			},
			"check": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of check to perform against backends to ensure they are serving requests. This is used to determine if backends are up or down. If none no check is performed. connection requires only a connection to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned matches what is expected.",
				Optional:    true,
				Computed:    true,
			},
			"check_path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The URL path to check on each backend. If the backend does not respond to this request it is considered to be down.",
				Optional:    true,
				Computed:    true,
			},
			"check_body": &schema.Schema{
				Type:        schema.TypeString,
				Description: "This value must be present in the response body of the check in order for it to pass. If this value is not present in the response body of a check request, the backend is considered to be down",
				Optional:    true,
				Computed:    true,
			},
			"check_passive": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If true, any response from this backend with a 5xx status code will be enough for it to be considered unhealthy and taken out of rotation.",
				Optional:    true,
				Computed:    true,
			},
			"cipher_suite": &schema.Schema{
				Type:        schema.TypeString,
				Description: "What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is considered insecure and should only be used if necessary.",
				Optional:    true,
				Computed:    true,
			},
			"ssl_commonname": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The common name for the SSL certification this port is serving if this port is not configured to use SSL.",
				Computed:    true,
			},
			"ssl_fingerprint": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The fingerprint for the SSL certification this port is serving if this port is not configured to use SSL.",
				Computed:    true,
			},
			"ssl_cert": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The certificate this port is serving. This is not returned. If set, this field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.",
				Optional:    true,
			},
			"ssl_key": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The private key corresponding to this port's certificate. This is not returned. If set, this field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.",
				Optional:    true,
				Sensitive:   true,
			},
			"node_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_up": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of backends considered to be 'UP' and healthy, and that are serving requests.",
							Computed:    true,
						},
						"status_down": &schema.Schema{
							Type:        schema.TypeInt,
							Description: "The number of backends considered to be 'DOWN' and unhealthy. These are not in rotation, and not serving requests.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceLinodeNodeBalancerConfigExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return false, fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}

	_, err = client.GetNodeBalancerConfig(context.Background(), int(nodebalancerID), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode NodeBalancerConfig ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func syncConfigResourceData(d *schema.ResourceData, config *linodego.NodeBalancerConfig) {
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
	d.Set("ssl_key", config.SSLKey)
	d.Set("ssl_fingerprint", config.SSLFingerprint)
	d.Set("ssl_commonname", config.SSLCommonName)
	nodeStatus := map[string]interface{}{
		"up":   fmt.Sprintf("%d", config.NodesStatus.Up),
		"down": fmt.Sprintf("%d", config.NodesStatus.Down),
	}
	if err := d.Set("node_status", nodeStatus); err != nil {
		panic(err)
	}
}

func resourceLinodeNodeBalancerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}

	nodebalancer, err := client.GetNodeBalancerConfig(context.Background(), int(nodebalancerID), int(id))

	if err != nil {
		return fmt.Errorf("Error finding the specified Linode NodeBalancerConfig: %s", err)
	}

	syncConfigResourceData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode NodeBalancerConfig")
	}

	nodebalancerID := d.Get("nodebalancer_id").(int)

	createOpts := linodego.NodeBalancerConfigCreateOptions{
		Algorithm:     linodego.ConfigAlgorithm(d.Get("algorithm").(string)),
		Check:         linodego.ConfigCheck(d.Get("check").(string)),
		Stickiness:    linodego.ConfigStickiness(d.Get("stickiness").(string)),
		CheckAttempts: d.Get("check_attempts").(int),
		CheckBody:     d.Get("check_body").(string),
		CheckInterval: d.Get("check_interval").(int),
		CheckPath:     d.Get("check_path").(string),
		CheckTimeout:  d.Get("check_timeout").(int),
		Port:          d.Get("port").(int),
		Protocol:      linodego.ConfigProtocol(d.Get("protocol").(string)),
		SSLCert:       d.Get("ssl_cert").(string),
		SSLKey:        d.Get("ssl_key").(string),
	}

	if checkPassive, ok := d.GetOk("check_passive"); ok {
		checkPassive := checkPassive.(bool)
		createOpts.CheckPassive = &checkPassive
	}

	config, err := client.CreateNodeBalancerConfig(context.Background(), nodebalancerID, createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode NodeBalancerConfig: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", config.ID))
	d.Set("nodebalancer_id", nodebalancerID)

	syncConfigResourceData(d, config)

	return nil
}

func resourceLinodeNodeBalancerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %s as int", d.Get("nodebalancer_id"))
	}

	config, err := client.GetNodeBalancerConfig(context.Background(), nodebalancerID, int(id))
	if err != nil {
		return fmt.Errorf("Error fetching data about the current NodeBalancerConfig: %s", err)
	}

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		Algorithm:     linodego.ConfigAlgorithm(d.Get("algorithm").(string)),
		Check:         linodego.ConfigCheck(d.Get("check").(string)),
		Stickiness:    linodego.ConfigStickiness(d.Get("stickiness").(string)),
		CheckAttempts: d.Get("check_attempts").(int),
		CheckBody:     d.Get("check_body").(string),
		CheckInterval: d.Get("check_interval").(int),
		CheckPath:     d.Get("check_path").(string),
		Port:          d.Get("port").(int),
		Protocol:      linodego.ConfigProtocol(d.Get("protocol").(string)),
		SSLCert:       d.Get("ssl_cert").(string),
		SSLKey:        d.Get("ssl_key").(string),
	}

	if checkPassive, ok := d.GetOk("check_passive"); ok {
		checkPassive := checkPassive.(bool)
		updateOpts.CheckPassive = &checkPassive
	}

	if config, err = client.UpdateNodeBalancerConfig(context.Background(), int(nodebalancerID), int(id), updateOpts); err != nil {
		return err
	}
	syncConfigResourceData(d, config)

	return nil
}

func resourceLinodeNodeBalancerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return fmt.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	err = client.DeleteNodeBalancerConfig(context.Background(), nodebalancerID, int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode NodeBalancerConfig %d: %s", id, err)
	}
	return nil
}
