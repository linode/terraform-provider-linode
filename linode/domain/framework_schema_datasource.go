package domain

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Description: "The Domain's unique ID.",
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.ConflictsWith(path.Expressions{
					path.MatchRoot("domain"),
				}...),
			},
		},
		"domain": schema.StringAttribute{
			Description: "The domain this Domain represents. These must be unique in our system; you cannot have " +
				"two Domains representing the same domain.",
			Optional: true,
		},
		"type": schema.StringAttribute{
			Description: "If this Domain represents the authoritative source of information for the domain it " +
				"describes, or if it is a read-only copy of a master (also called a slave).",
			Computed: true,
		},
		"group": schema.StringAttribute{
			Description: "The group this Domain belongs to. This is for display purposes only.",
			Computed:    true,
		},
		"status": schema.StringAttribute{
			Description: "Used to control whether this Domain is currently being rendered.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "A description for this Domain. This is for display purposes only.",
			Computed: true,
		},
		"master_ips": schema.SetAttribute{
			Description: "The IP addresses representing the master DNS for this Domain.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"axfr_ips": schema.SetAttribute{
			Description: "The list of IPs that may perform a zone transfer for this Domain. This is potentially " +
				"dangerous, and should be set to an empty list unless you intend to use it.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"ttl_sec": schema.Int64Attribute{
			Description: "'Time to Live' - the amount of time in seconds that this Domain's records may be " +
				"cached by resolvers or other domain servers. " + domainSecondsDescription,
			Computed: true,
		},
		"retry_sec": schema.Int64Attribute{
			Description: "The interval, in seconds, at which a failed refresh should be retried. " +
				domainSecondsDescription,
			Computed: true,
		},
		"expire_sec": schema.Int64Attribute{
			Description: "The amount of time in seconds that may pass before this Domain is no longer " +
				domainSecondsDescription,
			Computed: true,
		},
		"refresh_sec": schema.Int64Attribute{
			Description: "The amount of time in seconds before this Domain should be refreshed. " +
				domainSecondsDescription,
			Computed: true,
		},
		"soa_email": schema.StringAttribute{
			Description: "Start of Authority email address. This is required for master Domains.",
			Computed:    true,
		},
		"tags": schema.SetAttribute{
			Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			Computed:    true,
			ElementType: types.StringType,
		},
	},
}
