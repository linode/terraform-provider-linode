package domainrecord

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	"domain_id": {
		Type:        schema.TypeInt,
		Description: "The ID of the Domain to access.",
		Required:    true,
		ForceNew:    true,
	},
	"name": {
		Type: schema.TypeString,
		Description: "The name of this Record. This field's actual usage depends on the type of record this " +
			"represents. For A and AAAA records, this is the subdomain being associated with an IP address. " +
			"Generated for SRV records.",
		Optional:     true,
		Computed:     true, // This is true for SRV records
		ValidateFunc: validation.StringLenBetween(0, 100),
	},
	"record_type": {
		Type: schema.TypeString,
		Description: "The type of Record this is in the DNS system. For example, A records associate a " +
			"domain name with an IPv4 address, and AAAA records associate a domain name with an IPv6 address.",
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringInSlice(
			[]string{"A", "AAAA", "NS", "MX", "CNAME", "TXT", "SRV", "PTR", "CAA"}, false),
	},
	"ttl_sec": {
		Type: schema.TypeInt,
		Description: "'Time to Live' - the amount of time in seconds that this Domain's records may be " +
			"cached by resolvers or other domain servers. Valid values are 30, 120, 300, 3600, 7200, 14400, 28800, 57600, " +
			"86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest " +
			"valid value.",
		Optional:         true,
		DiffSuppressFunc: helper.DomainSecondsDiffSuppressor(),
	},
	"target": {
		Type: schema.TypeString,
		Description: "The target for this Record. This field's actual usage depends on the type of record " +
			"this represents. For A and AAAA records, this is the address the named Domain should resolve to.",
		Required:         true,
		DiffSuppressFunc: domainRecordTargetSuppressor,
	},
	"priority": {
		Type:         schema.TypeInt,
		Description:  "The priority of the target host. Lower values are preferred.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 255),
	},
	"protocol": {
		Type:        schema.TypeString,
		Description: "The protocol this Record's service communicates with. Only valid for SRV records.",
		Optional:    true,
	},
	"service": {
		Type:        schema.TypeString,
		Description: "The service this Record identified. Only valid for SRV records.",
		Optional:    true,
	},
	"tag": {
		Type:        schema.TypeString,
		Description: "The tag portion of a CAA record. It is invalid to set this on other record types.",
		Optional:    true,
	},
	"port": {
		Type:        schema.TypeInt,
		Description: "The port this Record points to.",
		Optional:    true,
	},
	"weight": {
		Type:        schema.TypeInt,
		Description: "The relative weight of this Record. Higher values are preferred.",
		Optional:    true,
	},
}
