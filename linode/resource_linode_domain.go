package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeDomainCreate,
		Read:   resourceLinodeDomainRead,
		Update: resourceLinodeDomainUpdate,
		Delete: resourceLinodeDomainDelete,
		Exists: resourceLinodeDomainExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The domain this Domain represents. These must be unique in our system; you cannot have two Domains representing the same domain.",
				Required:    true,
			},
			"domain_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "If this Domain represents the authoritative source of information for the domain it describes, or if it is a read-only copy of a master (also called a slave).",
				Default:     "master",
				Optional:    true,
				ForceNew:    true,
			},
			"group": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The group this Domain belongs to. This is for display purposes only.",
				Optional:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Used to control whether this Domain is currently being rendered.",
				Optional:    true,
				Default:     "active",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Description: "A description for this Domain. This is for display purposes only.",
				Optional:    true,
			},
			"master_ips": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IP addresses representing the master DNS for this Domain.",
				Optional:    true,
			},
			"axfr_ips": &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The list of IPs that may perform a zone transfer for this Domain. This is potentially dangerous, and should be set to an empty list unless you intend to use it.",
				Optional:    true,
			},
			"ttl_sec": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "'Time to Live' - the amount of time in seconds that this Domain's records may be cached by resolvers or other domain servers. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Optional:    true,
			},
			"retry_sec": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The interval, in seconds, at which a failed refresh should be retried. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Optional:    true,
			},
			"expire_sec": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The amount of time in seconds that may pass before this Domain is no longer authoritative. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Optional:    true,
			},
			"refresh_sec": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The amount of time in seconds before this Domain should be refreshed. Valid values are 300, 3600, 7200, 14400, 28800, 57600, 86400, 172800, 345600, 604800, 1209600, and 2419200 - any other value will be rounded to the nearest valid value.",
				Optional:    true,
			},
			"soa_email": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Start of Authority email address. This is required for master Domains.",
				Optional:    true,
			},
		},
	}
}

func syncResourceData(d *schema.ResourceData, domain *linodego.Domain) {
	d.Set("domain", domain.Domain)
	d.Set("domain_type", domain.Type)
	d.Set("group", domain.Group)
	d.Set("status", domain.Status)
	d.Set("description", domain.Description)
	d.Set("master_ips", domain.MasterIPs)
	d.Set("afxr_ips", domain.AXfrIPs)
	d.Set("ttl_sec", domain.TTLSec)
	d.Set("retry_sec", domain.RetrySec)
	d.Set("expire_sec", domain.ExpireSec)
	d.Set("refresh_sec", domain.RefreshSec)
	d.Set("soa_email", domain.SOAEmail)
}

func resourceLinodeDomainExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Error parsing Linode Domain ID %s as int: %s", d.Id(), err)
	}

	_, err = client.GetDomain(context.Background(), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return false, nil
		}

		return false, fmt.Errorf("Error getting Linode Domain ID %s: %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Domain ID %s as int: %s", d.Id(), err)
	}

	domain, err := client.GetDomain(context.Background(), int(id))

	if err != nil {
		return fmt.Errorf("Error finding the specified Linode Domain: %s", err)
	}

	syncResourceData(d, domain)

	return nil
}

func resourceLinodeDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Domain")
	}

	createOpts := linodego.DomainCreateOptions{
		Domain:      d.Get("domain").(string),
		Type:        linodego.DomainType(d.Get("domain_type").(string)),
		Group:       d.Get("group").(string),
		Description: d.Get("description").(string),
		SOAEmail:    d.Get("soa_email").(string),
		RetrySec:    d.Get("retry_sec").(int),
		ExpireSec:   d.Get("expire_sec").(int),
		RefreshSec:  d.Get("refresh_sec").(int),
		TTLSec:      d.Get("ttl_sec").(int),
	}

	if v, ok := d.GetOk("master_ips"); ok {
		var masterIPS []string
		for _, ip := range v.([]interface{}) {
			masterIPS = append(masterIPS, ip.(string))
		}

		createOpts.MasterIPs = masterIPS
	}

	if v, ok := d.GetOk("axfr_ips"); ok {
		var AXfrIPs []string
		for _, ip := range v.([]interface{}) {
			AXfrIPs = append(AXfrIPs, ip.(string))
		}

		createOpts.AXfrIPs = AXfrIPs
	}

	domain, err := client.CreateDomain(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("Error creating a Linode Domain: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", domain.ID))
	syncResourceData(d, domain)

	return resourceLinodeDomainRead(d, meta)
}

func resourceLinodeDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Domain id %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.DomainUpdateOptions{
		Domain:      d.Get("domain").(string),
		Status:      linodego.DomainStatus(d.Get("status").(string)),
		Type:        linodego.DomainType(d.Get("domain_type").(string)),
		Group:       d.Get("group").(string),
		Description: d.Get("description").(string),
		SOAEmail:    d.Get("soa_email").(string),
		RetrySec:    d.Get("retry_sec").(int),
		ExpireSec:   d.Get("expire_sec").(int),
		RefreshSec:  d.Get("refresh_sec").(int),
		TTLSec:      d.Get("ttl_sec").(int),
	}

	if v, ok := d.GetOk("master_ips"); ok {
		var masterIPS []string
		for _, ip := range v.([]interface{}) {
			masterIPS = append(masterIPS, ip.(string))
		}

		updateOpts.MasterIPs = masterIPS
	}

	if v, ok := d.GetOk("axfr_ips"); ok {
		var AXfrIPs []string
		for _, ip := range v.([]interface{}) {
			AXfrIPs = append(AXfrIPs, ip.(string))
		}

		updateOpts.AXfrIPs = AXfrIPs
	}

	domain, err := client.UpdateDomain(context.Background(), int(id), updateOpts)
	if err != nil {
		return fmt.Errorf("Error updating Linode Domain %d: %s", id, err)
	}
	syncResourceData(d, domain)

	return nil
}

func resourceLinodeDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing Linode Domain id %s as int", d.Id())
	}
	err = client.DeleteDomain(context.Background(), int(id))
	if err != nil {
		return fmt.Errorf("Error deleting Linode Domain %d: %s", id, err)
	}
	d.SetId("")

	return nil
}
