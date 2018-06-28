package linode

import (
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
				Type:         schema.TypeString,
				Description:  "The protocol this port is configured to serve. If this is set to https you must include an ssl_cert and an ssl_key.",
				Required:     true,
				InputDefault: "us-east",
				Default:      linodego.ProtocolHTTP,
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
				Default:     0,
			},
			"check_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "How long, in seconds, to wait for a check attempt before considering it failed.",
				Optional:    true,
				Default:     0,
			},
			"check_attempts": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "How many times to attempt a check before considering a backend to be down.",
				Optional:    true,
				Default:     0,
			},
			"algorithm": &schema.Schema{
				Type:        schema.TypeString,
				Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, leastconn, source",
				Optional:    true,
				Default:     linodego.AlgorithmRoundRobin,
			},
			"stickiness": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
				Optional:    true,
				Default:     linodego.StickinessNone,
			},
			"check": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The type of check to perform against backends to ensure they are serving requests. This is used to determine if backends are up or down. If none no check is performed. connection requires only a connection to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned matches what is expected.",
				Optional:    true,
				Default:     linodego.CheckHTTP,
			},
			"check_path": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The URL path to check on each backend. If the backend does not respond to this request it is considered to be down.",
				Optional:    true,
			},
			"check_body": &schema.Schema{
				Type:        schema.TypeString,
				Description: "This value must be present in the response body of the check in order for it to pass. If this value is not present in the response body of a check request, the backend is considered to be down",
				Optional:    true,
			},
			"check_passive": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "If true, any response from this backend with a 5xx status code will be enough for it to be considered unhealthy and taken out of rotation.",
				Optional:    true,
				Default:     true,
			},
			"cipher_suite": &schema.Schema{
				Type:        schema.TypeString,
				Description: "What ciphers to use for SSL connections served by this NodeBalancer. legacy is considered insecure and should only be used if necessary.",
				Optional:    true,
				Default:     linodego.CipherRecommended,
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
		},
	}
}

func resourceLinodeNodeBalancerConfigExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode NodeBalancerConfig ID %s as int because %s", d.Id(), err)
	}
	nodebalancerID, err := strconv.ParseInt(d.Get("nodebalancer_id").(string), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Get("nodebalancer"), err)
	}

	_, err = client.GetNodeBalancerConfig(int(nodebalancerID), int(id))
	if err != nil {
		return false, fmt.Errorf("Failed to get Linode NodeBalancerConfig ID %s because %s", d.Id(), err)
	}
	return true, nil
}

func syncConfigResourceData(d *schema.ResourceData, config *linodego.NodeBalancerConfig) {
	d.Set("algorithm", config.Algorithm)
	d.Set("check", config.Check)
	d.Set("check_attempts", config.CheckAttempts)
	d.Set("check_body", config.CheckBody)
	d.Set("check_interval", config.CheckInterval)
	d.Set("check_passive", config.CheckPassive)
	d.Set("check_path", config.CheckPath)
	d.Set("port", config.Port)
	d.Set("protocol", config.Protocol)
	d.Set("ssl_fingerprint", config.SSLFingerprint)
	d.Set("ssl_commonname", config.SSLCommonName)
}

func resourceLinodeNodeBalancerConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancerConfig ID %s as int because %s", d.Id(), err)
	}
	nodebalancerID, err := strconv.ParseInt(d.Get("nodebalancer_id").(string), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Get("nodebalancer"), err)
	}

	nodebalancer, err := client.GetNodeBalancerConfig(int(nodebalancerID), int(id))

	if err != nil {
		return fmt.Errorf("Failed to find the specified Linode NodeBalancerConfig because %s", err)
	}

	syncConfigResourceData(d, nodebalancer)

	return nil
}

func resourceLinodeNodeBalancerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode NodeBalancerConfig")
	}

	createOpts := linodego.NodeBalancerConfigCreateOptions{
		NodeBalancerID: d.Get("nodebalancer_id").(int),
		Algorithm:      d.Get("algorithm").(linodego.ConfigAlgorithm),
		Check:          d.Get("check").(linodego.ConfigCheck),
		CheckAttempts:  d.Get("check_attempts").(int),
		CheckBody:      d.Get("check_body").(string),
		CheckInterval:  d.Get("check_interval").(int),
		CheckPassive:   d.Get("check_passive").(bool),
		CheckPath:      d.Get("check_path").(string),
		Port:           d.Get("port").(int),
		Protocol:       d.Get("protocol").(linodego.ConfigProtocol),
		SSLCert:        d.Get("ssl_cert").(string),
		SSLKey:         d.Get("ssl_key").(string),
	}

	config, err := client.CreateNodeBalancerConfig(&createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode NodeBalancerConfig in because %s", err)
	}
	d.SetId(fmt.Sprintf("%d", config.ID))

	syncConfigResourceData(d, config)

	return nil
}

func resourceLinodeNodeBalancerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancerConfig ID %s as int because %s", d.Id(), err)
	}
	nodebalancerID, err := strconv.ParseInt(d.Get("nodebalancer_id").(string), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Get("nodebalancer"), err)
	}

	config, err := client.GetNodeBalancerConfig(int(nodebalancerID), int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current NodeBalancerConfig because %s", err)
	}

	updateOpts := linodego.NodeBalancerConfigUpdateOptions{
		NodeBalancerID: int(nodebalancerID),
		Algorithm:      d.Get("algorithm").(linodego.ConfigAlgorithm),
		Check:          d.Get("check").(linodego.ConfigCheck),
		CheckAttempts:  d.Get("check_attempts").(int),
		CheckBody:      d.Get("check_body").(string),
		CheckInterval:  d.Get("check_interval").(int),
		CheckPassive:   d.Get("check_passive").(bool),
		CheckPath:      d.Get("check_path").(string),
		Port:           d.Get("port").(int),
		Protocol:       d.Get("protocol").(linodego.ConfigProtocol),
		SSLCert:        d.Get("ssl_cert").(string),
		SSLKey:         d.Get("ssl_key").(string),
	}

	if config, err = client.UpdateNodeBalancerConfig(int(nodebalancerID), int(id), updateOpts); err != nil {
		return err
	}
	syncConfigResourceData(d, config)

	return nil
}

func resourceLinodeNodeBalancerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancerConfig ID %s as int because %s", d.Id(), err)
	}
	nodebalancerID, err := strconv.ParseInt(d.Get("nodebalancer_id").(string), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode NodeBalancer ID %s as int because %s", d.Get("nodebalancer"), err)
	}
	err = client.DeleteNodeBalancerConfig(int(nodebalancerID), int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode NodeBalancerConfig %d because %s", id, err)
	}
	return nil
}
