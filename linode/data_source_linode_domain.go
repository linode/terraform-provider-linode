package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeDomainRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).",
				Computed:    true,
			},
			"group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The group this Domain belongs to. This is for display purposes only.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used to control whether this Domain is currently being rendered.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description for this Domain. This is for display purposes only.",
			},
			"master_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IP addresses representing the master DNS for this Domain.",
				Computed:    true,
			},
			"axfr_ips": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.",
				Computed:    true,
			},
			"ttl_sec": {
				Type:        schema.TypeInt,
				Description: "'Time to Live' - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Computed:    true,
			},
			"retry_sec": {
				Type:        schema.TypeInt,
				Description: "The interval, in seconds, at which a failed refresh should be retried. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Computed:    true,
			},
			"expire_sec": {
				Type:        schema.TypeInt,
				Description: "The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Computed:    true,
			},
			"refresh_sec": {
				Type:        schema.TypeInt,
				Description: "The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Computed:    true,
			},
			"soa_email": {
				Type:        schema.TypeString,
				Description: "Start of Authority email address. This is required for master Domains.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	reqDomain := d.Get("id").(string)

	if reqDomain == "" {
		return fmt.Errorf("Domain id is required")
	}

	reqDomainInt, err := strconv.Atoi(reqDomain)

	if err != nil {
		return fmt.Errorf("Domain id number be a valid integer")
	}

	domain, err := client.GetDomain(context.Background(), reqDomainInt)
	if err != nil {
		return fmt.Errorf("Error listing domain: %s", err)
	}

	if domain != nil {
		d.SetId(strconv.Itoa(domain.ID))
		d.Set("domain", domain.Domain)
		d.Set("type", domain.Type)
		d.Set("group", domain.Group)
		d.Set("status", domain.Status)
		d.Set("description", domain.Description)
		d.Set("master_ips", domain.MasterIPs)
		d.Set("axfr_ips", domain.AXfrIPs)
		d.Set("ttl_sec", domain.TTLSec)
		d.Set("retry_sec", domain.RetrySec)
		d.Set("expire_sec", domain.ExpireSec)
		d.Set("refresh_sec", domain.RefreshSec)
		d.Set("soa_email", domain.SOAEmail)
		d.Set("tags", domain.Tags)
		return nil
	}

	d.SetId("")

	return fmt.Errorf("Domain %s was not found", reqDomain)
}
