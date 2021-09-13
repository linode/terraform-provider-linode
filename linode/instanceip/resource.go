package instanceip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	address := d.Id()
	linodeID := d.Get("linode_id").(int)
	ip, err := client.GetInstanceIPAddress(ctx, linodeID, address)
	if err != nil {
		return diag.Errorf("failed to get instance (%d) ip: %s", linodeID, err)
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

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	linodeID := d.Get("linode_id").(int)
	private := d.Get("public").(bool)
	applyImmediately := d.Get("apply_immediately").(bool)

	ip, err := client.AddInstanceIPAddress(ctx, linodeID, private)
	if err != nil {
		return diag.Errorf("failed to create instance (%d) ip: %s", linodeID, err)
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

	// Only reboot the associated instance with apply_immediately == true
	if applyImmediately {
		// TODO(jxriddle): check instance status and react to state

		// setting bootConfig to 0 to use current bootConfig
		if diagErr := helper.RebootInstance(ctx, d, linodeID, meta, 0); diagErr != nil {
			return diagErr
		}
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

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

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	address := d.Id()
	linodeID := d.Get("linode_id").(int)
	if err := client.DeleteInstanceIPAddress(ctx, linodeID, address); err != nil {
		return diag.Errorf("failed to delete instance (%d) ip (%s): %s", linodeID, address, err)
	}
	return nil
}
