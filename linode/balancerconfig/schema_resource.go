package balancerconfig

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
)

var resourceSchemaStatus = map[string]*schema.Schema{
	"up": {
		Type:        schema.TypeInt,
		Description: "The number of backends considered to be 'UP' and healthy, and that are serving requests.",
		Computed:    true,
	},
	"down": {
		Type: schema.TypeInt,
		Description: "The number of backends considered to be 'DOWN' and unhealthy. These are not in " +
			"rotation, and not serving requests.",
		Computed: true,
	},
}

var resourceSchema = map[string]*schema.Schema{
	"nodebalancer_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the NodeBalancer to access.",
		Required:    true,
		ForceNew:    true,
	},
	"protocol": {
		Type: schema.TypeString,
		Description: "The protocol this port is configured to serve. If this is set to https you must " +
			"include an ssl_cert and an ssl_key.",
		StateFunc: func(val interface{}) string {
			return strings.ToLower(val.(string))
		},
		Optional: true,
		Default:  linodego.ProtocolHTTP,
	},
	"proxy_protocol": {
		Type: schema.TypeString,
		Description: "The version of ProxyProtocol to use for the underlying NodeBalancer. " +
			"This requires protocol to be `tcp`. Valid values are `none`, `v1`, and `v2`.",
		ValidateFunc: validation.StringInSlice([]string{"none", "v1", "v2"}, false),
		Optional:     true,
		Default:      linodego.ProxyProtocolNone,
	},
	"port": {
		Type: schema.TypeInt,
		Description: "The TCP port this Config is for. These values must be unique across configs on a " +
			"single NodeBalancer (you can't have two configs for port 80, for example). While some ports imply " +
			"some protocols, no enforcement is done and you may configure your NodeBalancer however is useful to " +
			"you. For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have " +
			"a NodeBalancer listening on port 443.",
		ValidateFunc: validation.IntBetween(1, 65535),
		Optional:     true,
		Default:      80,
	},
	"check_interval": {
		Type:        schema.TypeInt,
		Description: "How often, in seconds, to check that backends are up and serving requests.",
		Optional:    true,
		Computed:    true,
	},
	"check_timeout": {
		Type:         schema.TypeInt,
		Description:  "How long, in seconds, to wait for a check attempt before considering it failed. (1-30)",
		ValidateFunc: validation.IntBetween(1, 30),
		Optional:     true,
		Computed:     true,
	},
	"check_attempts": {
		Type:         schema.TypeInt,
		Description:  "How many times to attempt a check before considering a backend to be down. (1-30)",
		ValidateFunc: validation.IntBetween(1, 30),
		Optional:     true,
		Computed:     true,
	},
	"algorithm": {
		Type: schema.TypeString,
		Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, " +
			"leastconn, source",
		ValidateFunc: validation.StringInSlice([]string{"roundrobin", "leastconn", "source"}, false),
		Optional:     true,
		Computed:     true,
	},
	"stickiness": {
		Type:         schema.TypeString,
		Description:  "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
		ValidateFunc: validation.StringInSlice([]string{"none", "table", "http_cookie"}, false),
		Optional:     true,
		Computed:     true,
	},
	"check": {
		Type: schema.TypeString,
		Description: "The type of check to perform against backends to ensure they are serving requests. This is used " +
			"to determine if backends are up or down. If none no check is performed. connection requires only a connection " +
			"to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned " +
			"matches what is expected.",
		ValidateFunc: validation.StringInSlice([]string{"none", "connection", "http", "http_body"}, false),
		Optional:     true,
		Computed:     true,
	},
	"check_path": {
		Type: schema.TypeString,
		Description: "The URL path to check on each backend. If the backend does not respond to this request it is " +
			"considered to be down.",
		Optional: true,
		Computed: true,
	},
	"check_body": {
		Type: schema.TypeString,
		Description: "This value must be present in the response body of the check in order for it to pass. " +
			"If this value is not present in the response body of a check request, the backend is considered to be down",
		Optional: true,
		Computed: true,
	},
	"check_passive": {
		Type: schema.TypeBool,
		Description: "If true, any response from this backend with a 5xx status code will be enough for it to " +
			"be considered unhealthy and taken out of rotation.",
		Optional: true,
		Computed: true,
	},
	"cipher_suite": {
		Type: schema.TypeString,
		Description: "What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is " +
			"considered insecure and should only be used if necessary.",
		ValidateFunc: validation.StringInSlice([]string{"recommended", "legacy"}, false),
		Optional:     true,
		Computed:     true,
	},
	"ssl_commonname": {
		Type: schema.TypeString,
		Description: "The read-only common name automatically derived from the SSL certificate assigned to this " +
			"NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned " +
			"to your NodeBalancerConfig.",
		Computed: true,
	},
	"ssl_fingerprint": {
		Type: schema.TypeString,
		Description: "The read-only fingerprint automatically derived from the SSL certificate assigned to this " +
			"NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned to " +
			"your NodeBalancerConfig.",
		Computed: true,
	},
	"ssl_cert": {
		Type: schema.TypeString,
		Description: "The certificate this port is serving. This is not returned. If set, this field will come " +
			"back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.",
		Optional:  true,
		Sensitive: true,
	},
	"ssl_key": {
		Type: schema.TypeString,
		Description: "The private key corresponding to this port's certificate. This is not returned. If set, this " +
			"field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify " +
			"the certificate.",
		Optional:  true,
		Sensitive: true,
	},
	"node_status": {
		Type: schema.TypeList,
		Description: "A structure containing information about the health of the backends for this port. This " +
			"information is updated periodically as checks are performed against backends.",
		Computed: true,
		Elem:     resourceStatus(),
	},
}
