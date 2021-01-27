package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeInstanceIP() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeInstanceIPCreate,
		ReadContext:   resourceLinodeInstanceIPRead,
		UpdateContext: resourceLinodeInstanceIPUpdate,
		DeleteContext: resourceLinodeInstanceIPDelete,

		Schema: map[string]*schema.Schema{
			"linode_id": {
				Type:        schema.TypeInt,
				Description: "The ID of the Linode to allocate an IPv4 address for.",
				Required:    true,
				ForceNew:    true,
			},
			"public": {
				Type:        schema.TypeBool,
				Description: "Whether the IPv4 address is public or private.",
				Default:     true,
				Optional:    true,
				ForceNew:    true,
			},

			"address": {
				Type:        schema.TypeString,
				Description: "The resulting IPv4 address.",
				Computed:    true,
			},
			"gateway": {
				Type:        schema.TypeString,
				Description: "The default gateway for this address",
				Computed:    true,
			},
			"prefix": {
				Type:        schema.TypeInt,
				Description: "The number of bits set in the subnet mask.",
				Computed:    true,
			},
			"rdns": {
				Type:        schema.TypeString,
				Description: "The reverse DNS assigned to this address.",
				Optional:    true,
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The region this IP resides in.",
				Computed:    true,
			},
			"subnet_mask": {
				Type:        schema.TypeString,
				Description: "The mask that separates host bits from network bits for this address.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of IP address.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeInstanceIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	address := d.Id()
	linodeID := d.Get("linode_id").(int)
	ip, err := client.GetInstanceIPAddress(ctx, linodeID, address)
	if err != nil {
		diag.Errorf("failed to get instance (%d) ip: %s", linodeID, err)
	}

	d.Set("address", ip.Address)
	d.Set("gateway", ip.Gateway)
	d.Set("prefix", ip.Prefix)
	d.Set("rdns", ip.RDNS)
	d.Set("region", ip.Region)
	d.Set("subnet_mask", ip.SubnetMask)
	d.Set("type", ip.Type)
	return nil
}

func resourceLinodeInstanceIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)
	private := d.Get("public").(bool)
	ip, err := client.AddInstanceIPAddress(ctx, linodeID, private)
	if err != nil {
		diag.Errorf("failed to create instance (%d) ip: %s", linodeID, err)
	}

	rdns := d.Get("rdns").(string)
	if rdns != "" {
		if _, err := client.UpdateIPAddress(ctx, ip.Address, linodego.IPAddressUpdateOptions{
			RDNS: &rdns,
		}); err != nil {
			return diag.Errorf("failed to set RDNS for instance (%d) ip (%s): %s", linodeID, ip.Address, err)
		}
	}

	d.SetId(ip.Address)
	return resourceLinodeInstanceIPRead(ctx, d, meta)
}

func resourceLinodeInstanceIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	address := d.Id()
	linodeID := d.Get("linode_id").(int)
	rdns := d.Get("rdns").(string)
	if d.HasChange("rdns") {
		updateOptions := linodego.IPAddressUpdateOptions{}
		if rdns != "" {
			updateOptions.RDNS = &rdns
		}

		if _, err := client.UpdateIPAddress(ctx, address, linodego.IPAddressUpdateOptions{
			RDNS: &rdns,
		}); err != nil {
			return diag.Errorf("failed to update RDNS for instance (%d) ip (%s): %s", linodeID, address, err)
		}
	}
	return nil
}

func resourceLinodeInstanceIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	address := d.Id()
	linodeID := d.Get("linode_id").(int)
	if err := client.DeleteInstanceIPAddress(ctx, linodeID, address); err != nil {
		return diag.Errorf("failed to delete instance (%d) ip (%s): %s", linodeID, address, err)
	}
	return nil
}
