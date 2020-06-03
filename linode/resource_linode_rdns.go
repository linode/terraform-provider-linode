package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/linodego"
)

func resourceLinodeRDNS() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeRDNSCreateContext,
		ReadContext:   resourceLinodeRDNSReadContext,
		DeleteContext: resourceLinodeRDNSDeleteContext,
		UpdateContext: resourceLinodeRDNSUpdateContext,
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
				Type:         schema.TypeString,
				Description:  "The reverse DNS assigned to this address. For public IPv4 addresses, this will be set to a default value provided by Linode if not explicitly set.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 254),
			},
		},
	}
}

func resourceLinodeRDNSReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	ip, err := client.GetIPAddress(context.Background(), ipStr)

	if err != nil {
		return diag.Errorf("Error finding the specified Linode RDNS: %s", err)
	}

	d.Set("address", d.Id())
	d.Set("rdns", ip.RDNS)

	return nil
}

func resourceLinodeRDNSCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(linodego.Client)
	if !ok {
		return diag.Errorf("Invalid Client when creating Linode RDNS")
	}

	var address = d.Get("address").(string)
	var rdns *string
	if rdnsRaw, ok := d.GetOk("rdns"); ok && len(rdnsRaw.(string)) > 0 {
		rdnsStr := rdnsRaw.(string)
		rdns = &rdnsStr
	}
	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: rdns,
	}
	ip, err := client.UpdateIPAddress(context.Background(), address, updateOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode RDNS: %s", err)
	}
	d.SetId(address)
	d.Set("rdns", ip.RDNS)

	return resourceLinodeRDNSReadContext(ctx, d, meta)
}

func resourceLinodeRDNSUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
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

	if _, err := client.UpdateIPAddress(context.Background(), d.Id(), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode RDNS: %s", err)
	}

	return resourceLinodeRDNSReadContext(ctx, d, meta)
}

func resourceLinodeRDNSDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	ipStr := d.Id()

	if len(ipStr) == 0 {
		return diag.Errorf("Error parsing Linode RDNS ID %s as IP string", ipStr)
	}

	updateOpts := linodego.IPAddressUpdateOptions{
		RDNS: nil,
	}

	if _, err := client.UpdateIPAddress(context.Background(), d.Id(), updateOpts); err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error deleting Linode RDNS: %s", err)
	}

	d.SetId("")

	return nil
}
