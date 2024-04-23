package nbconfig

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

var frameworkResourceSchemaV1 = schema.Schema{
	Version:    1,
	Attributes: getSchemaAttributes(1),
}

var frameworkResourceSchemaV0 = schema.Schema{
	Version:    0,
	Attributes: getSchemaAttributes(0),
}

var NodeStatusTypeV0 = types.StringType

var NodeStatusTypeV1 = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"up":   types.Int64Type,
		"down": types.Int64Type,
	},
}

func getSchemaAttributes(version int) map[string]schema.Attribute {
	result := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the Linode NodeBalancer Config.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
		"nodebalancer_id": schema.Int64Attribute{
			Description: "The ID of the NodeBalancer to access.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"protocol": schema.StringAttribute{
			Description: "The protocol this port is configured to serve. If this is set to https you must " +
				"include an ssl_cert and an ssl_key.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.ProtocolHTTP),
					string(linodego.ProtocolHTTPS),
					string(linodego.ProtocolTCP),
				),
			},
			Default:  stringdefault.StaticString(string(linodego.ProtocolHTTP)),
			Optional: true,
			Computed: true,
		},
		"proxy_protocol": schema.StringAttribute{
			Description: "The version of ProxyProtocol to use for the underlying NodeBalancer. " +
				"This requires protocol to be `tcp`. Valid values are `none`, `v1`, and `v2`.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.ProxyProtocolV1),
					string(linodego.ProxyProtocolV2),
					string(linodego.ProxyProtocolNone),
				),
			},
			Default:  stringdefault.StaticString(string(linodego.ProxyProtocolNone)),
			Optional: true,
			Computed: true,
		},
		"port": schema.Int64Attribute{
			Description: "The TCP port this Config is for. These values must be unique across configs on a " +
				"single NodeBalancer (you can't have two configs for port 80, for example). While some ports imply " +
				"some protocols, no enforcement is done and you may configure your NodeBalancer however is useful to " +
				"you. For example, while port 443 is generally used for HTTPS, you do not need SSL configured to have " +
				"a NodeBalancer listening on port 443.",
			Validators: []validator.Int64{
				int64validator.Between(1, 65535),
			},
			Default:  int64default.StaticInt64(80),
			Optional: true,
			Computed: true,
		},
		"check_interval": schema.Int64Attribute{
			Description: "How often, in seconds, to check that backends are up and serving requests.",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"check_timeout": schema.Int64Attribute{
			Description: "How long, in seconds, to wait for a check attempt before considering it failed. (1-30)",
			Validators: []validator.Int64{
				int64validator.Between(1, 30),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"check_attempts": schema.Int64Attribute{
			Description: "How many times to attempt a check before considering a backend to be down. (1-30)",
			Validators: []validator.Int64{
				int64validator.Between(1, 30),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"algorithm": schema.StringAttribute{
			Description: "What algorithm this NodeBalancer should use for routing traffic to backends: roundrobin, " +
				"leastconn, source",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.AlgorithmRoundRobin),
					string(linodego.AlgorithmLeastConn),
					string(linodego.AlgorithmSource),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"stickiness": schema.StringAttribute{
			Description: "Controls how session stickiness is handled on this port: 'none', 'table', 'http_cookie'",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.StickinessNone),
					string(linodego.StickinessTable),
					string(linodego.StickinessHTTPCookie),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"check": schema.StringAttribute{
			Description: "The type of check to perform against backends to ensure they are serving requests. This is used " +
				"to determine if backends are up or down. If none no check is performed. connection requires only a connection " +
				"to the backend to succeed. http and http_body rely on the backend serving HTTP, and that the response returned " +
				"matches what is expected.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.CheckNone),
					string(linodego.CheckConnection),
					string(linodego.CheckHTTP),
					string(linodego.CheckHTTPBody),
				),
			},
			Optional: true,
			Computed: true,
		},
		"check_path": schema.StringAttribute{
			Description: "The URL path to check on each backend. If the backend does not respond to this request it is " +
				"considered to be down.",
			Optional: true,
			Computed: true,
		},
		"check_body": schema.StringAttribute{
			Description: "This value must be present in the response body of the check in order for it to pass. " +
				"If this value is not present in the response body of a check request, the backend is considered to be down",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Optional: true,
			Computed: true,
		},
		"check_passive": schema.BoolAttribute{
			Description: "If true, any response from this backend with a 5xx status code will be enough for it to " +
				"be considered unhealthy and taken out of rotation.",
			Default:  booldefault.StaticBool(true),
			Optional: true,
			Computed: true,
		},
		"cipher_suite": schema.StringAttribute{
			Description: "What ciphers to use for SSL connections served by this NodeBalancer. `legacy` is " +
				"considered insecure and should only be used if necessary.",
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.CipherLegacy),
					string(linodego.CipherRecommended),
				),
			},
			Default:  stringdefault.StaticString(string(linodego.CipherRecommended)),
			Optional: true,
			Computed: true,
		},
		"ssl_commonname": schema.StringAttribute{
			Description: "The read-only common name automatically derived from the SSL certificate assigned to this " +
				"NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned " +
				"to your NodeBalancerConfig.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
		"ssl_fingerprint": schema.StringAttribute{
			Description: "The read-only fingerprint automatically derived from the SSL certificate assigned to this " +
				"NodeBalancerConfig. Please refer to this field to verify that the appropriate certificate is assigned to " +
				"your NodeBalancerConfig.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Computed: true,
		},
		"ssl_cert": schema.StringAttribute{
			Description: "The certificate this port is serving. This is not returned. If set, this field will come " +
				"back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify the certificate.",
			Optional:  true,
			Sensitive: true,
		},
		"ssl_key": schema.StringAttribute{
			Description: "The private key corresponding to this port's certificate. This is not returned. If set, this " +
				"field will come back as `<REDACTED>`. Please use the ssl_commonname and ssl_fingerprint to identify " +
				"the certificate.",
			Optional:  true,
			Sensitive: true,
		},
	}

	nodeStatusDescription := "A structure containing information about the health of the backends for this port. This " +
		"information is updated periodically as checks are performed against backends."

	if version == 0 {
		result["node_status"] = schema.MapAttribute{
			Description: nodeStatusDescription,
			Computed:    true,
			ElementType: NodeStatusTypeV0,
		}
	} else {
		result["node_status"] = schema.ListAttribute{
			Description: nodeStatusDescription,
			Computed:    true,
			ElementType: NodeStatusTypeV1,
		}
	}

	return result
}
