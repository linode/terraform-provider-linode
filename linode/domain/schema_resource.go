package domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

const domainSecondsDescription = "Valid values are 0, 30, 120, 300, 3600, 7200, 14400, 28800, " +
	"57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to " +
	"the nearest valid value."

var resourceSchema = map[string]*schema.Schema{
	"domain": {
		Type: schema.TypeString,
		Description: "The domain this Domain represents. These must be unique in our system; you cannot have " +
			"two Domains representing the same domain.",
		Required: true,
	},
	"type": {
		Type: schema.TypeString,
		Description: "If this Domain represents the authoritative source of information for the domain it " +
			"describes, or if it is a read-only copy of a master (also called a slave).",
		InputDefault: "master",
		ValidateFunc: validation.StringInSlice([]string{"master", "slave"}, false),
		Required:     true,
		ForceNew:     true,
	},
	"group": {
		Type:         schema.TypeString,
		Description:  "The group this Domain belongs to. This is for display purposes only.",
		ValidateFunc: validation.StringLenBetween(0, 50),
		Optional:     true,
	},
	"status": {
		Type:         schema.TypeString,
		Description:  "Used to control whether this Domain is currently being rendered.",
		Optional:     true,
		Computed:     true,
		InputDefault: "active",
	},
	"description": {
		Type:         schema.TypeString,
		Description:  "A description for this Domain. This is for display purposes only.",
		ValidateFunc: validation.StringLenBetween(0, 255),
		Optional:     true,
	},
	"master_ips": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "The IP addresses representing the master DNS for this Domain.",
		Optional:    true,
	},
	"axfr_ips": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Description: "The list of IPs that may perform a zone transfer for this Domain. This is potentially " +
			"dangerous, and should be set to an empty list unless you intend to use it.",
		Optional: true,
	},
	"ttl_sec": {
		Type: schema.TypeInt,
		Description: "'Time to Live' - the amount of time in seconds that this Domain's records may be " +
			"cached by resolvers or other domain servers. " + domainSecondsDescription,
		Optional:         true,
		DiffSuppressFunc: helper.DomainSecondsDiffSuppressor(),
	},
	"retry_sec": {
		Type: schema.TypeInt,
		Description: "The interval, in seconds, at which a failed refresh should be retried. " +
			domainSecondsDescription,
		Optional:         true,
		DiffSuppressFunc: helper.DomainSecondsDiffSuppressor(),
	},
	"expire_sec": {
		Type: schema.TypeInt,
		Description: "The amount of time in seconds that may pass before this Domain is no longer " +
			domainSecondsDescription,
		Optional:         true,
		DiffSuppressFunc: helper.DomainSecondsDiffSuppressor(),
	},
	"refresh_sec": {
		Type: schema.TypeInt,
		Description: "The amount of time in seconds before this Domain should be refreshed. " +
			domainSecondsDescription,
		Optional:         true,
		DiffSuppressFunc: helper.DomainSecondsDiffSuppressor(),
	},
	"soa_email": {
		Type:        schema.TypeString,
		Description: "Start of Authority email address. This is required for master Domains.",
		Optional:    true,
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
}
