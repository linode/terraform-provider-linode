package linode

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func resourceLinodeRDNS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeRDNSCreate,
		ReadContext:   resourceLinodeRDNSRead,
		DeleteContext: resourceLinodeRDNSDelete,
		UpdateContext: resourceLinodeRDNSUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"address": {
				Type:         schema.TypeString,
				Description:  "The public Linode IPv4 or IPv6 address to operate on.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"rdns": {
				Type: schema.TypeString,
				Description: "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set " +
					"to a default value provided by Linode if not explicitly set.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 254),
			},
		},
	}
}

func resourceLinodeRDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	ip, err := client.GetIPAddress(ctx, ipStr)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode RDNS %q from state because it no longer exists", ipStr)
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode RDNS: %s", err)
	}

	d.Set("address", d.Id())
	d.Set("rdns", ip.RDNS)

	return nil
}

func resourceLinodeRDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	address := d.Get("address").(string)
	var rdns *string
	if rdnsRaw, ok := d.GetOk("rdns"); ok && len(rdnsRaw.(string)) > 0 {
		rdnsStr := rdnsRaw.(string)
		rdns = &rdnsStr
	}
	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: rdns,
	}
	ip, err := client.UpdateIPAddress(ctx, address, updateOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode RDNS: %s", err)
	}
	d.SetId(address)
	d.Set("rdns", ip.RDNS)

	return resourceLinodeRDNSRead(ctx, d, meta)
}

func resourceLinodeRDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	var rdns *string

	if rdnsRaw, ok := d.GetOk("rdns"); ok && len(rdnsRaw.(string)) > 0 {
		rdnsStr := rdnsRaw.(string)
		rdns = &rdnsStr
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: rdns,
	}

	if _, err := client.UpdateIPAddress(ctx, d.Id(), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode RDNS: %s", err)
	}

	return resourceLinodeRDNSRead(ctx, d, meta)
}

func resourceLinodeRDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: nil,
	}

	if _, err := client.UpdateIPAddress(ctx, d.Id(), updateOpts); err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error deleting Linode RDNS: %s", err)
	}

	d.SetId("")

	return nil
}
