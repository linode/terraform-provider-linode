package nbconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var statusObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"up":   types.Int64Type,
		"down": types.Int64Type,
	},
}

var frameworkDatasourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer config.",
			Required:    true,
		},
		"nodebalancer_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer to access.",
			Required:    true,
		},
		"protocol": schema.StringAttribute{
			Description: "The protocol this port is configured to serve. If this is set to https you must include " +
				"an ssl_cert and an ssl_key.",
			Computed: true,
		},
		"proxy_protocol": schema.StringAttribute{
			Description: "The version of ProxyProtocol to use for the underlying NodeBalancer. This requires " +
				"protocol to be `tcp`. Valid values are `none`, `v1`, and `v2`.",
			Computed: true,
		},
		"port": schema.Int64Attribute{
			Description: "The TCP port this Config is for. These values must be unique across configs on a single " +
				"NodeBalancer (you can't have two configs for port 80, for example). While some ports imply some " +
				"protocols, no enforcement is done and you may configure your NodeBalancer however is useful to you. " +
				"For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have a " +
				"NodeBalancer listening on port 443.",
			Computed: true,
		},
		"check_interval": schema.Int64Attribute{
			Description: "How often, in seconds, to check that backends are up and serving requests.",
			Computed:    true,
		},
		"check_timeout": schema.Int64Attribute{
			Description: "How long, in seconds, to wait for a check attempt before considering it failed. (1-30)",
			Computed:    true,
		},
		"check_attempts": schema.Int64Attribute{
			Description: "How many times to attempt a check before considering a backend to be down. (1-30)",
			Computed:    true,
		},
		"algorithm": schema.StringAttribute{
			Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, " +
				"leastconn, source",
			Computed: true,
		},
		"stickiness": schema.StringAttribute{
			Description: "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
			Computed:    true,
		},
		"check": schema.StringAttribute{
			Description: "The type of check to perform against backends to ensure they are serving requests. " +
				"This is used to determine if backends are up or down. If none no check is performed. " +
				"connection requires only a connection to the backend to succeed. http and http_body rely on the " +
				"backend serving HTTP, and that the response returned matches what is expected.",
			Computed: true,
		},
		"check_path": schema.StringAttribute{
			Description: "The URL path to check on each backend. If the backend does not respond to this request " +
				"it is considered to be down.",
			Computed: true,
		},
		"check_body": schema.StringAttribute{
			Description: "This value must be present in the response body of the check in order for it to pass. " +
				"If this value is not present in the response body of a check request, the backend is considered to be down",
			Computed: true,
		},
		"check_passive": schema.BoolAttribute{
			Description: "If true, any response from this backend with a 5xx status code will be enough for it to " +
				"be considered unhealthy and taken out of rotation.",
			Computed: true,
		},
		"cipher_suite": schema.StringAttribute{
			Description: "What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is " +
				"considered insecure and should only be used if necessary.",
			Computed: true,
		},
		"ssl_commonname": schema.StringAttribute{
			Description: "The read-only common name automatically derived from the SSL certificate assigned to " +
				"this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate " +
				"is assigned to your NodeBalancerConfig.",
			Computed: true,
		},
		"ssl_fingerprint": schema.StringAttribute{
			Description: "The read-only fingerprint automatically derived from the SSL certificate assigned to " +
				"this NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate " +
				"is assigned to your NodeBalancerConfig.",
			Computed: true,
		},
		"node_status": schema.ListAttribute{
			Description: "A structure containing information about the health of the backends for this port. " +
				"This information is updated periodically as checks are performed against backends.",
			Computed:    true,
			ElementType: statusObjectType,
		},
	},
}
