package instanceip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
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
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_instance_ip")

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
	d.Set("vpc_nat_1_1", []map[string]any{flattenVPCNAT1To1(ip.VPCNAT1To1)})
	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_instance_ip")

	linodeID := d.Get("linode_id").(int)
	private := d.Get("public").(bool)
	applyImmediately := d.Get("apply_immediately").(bool)

	client := meta.(*helper.ProviderMeta).Client

	ip, err := client.AddInstanceIPAddress(ctx, linodeID, private)
	if err != nil {
		return diag.Errorf("failed to create instance (%d) ip: %s", linodeID, err)
	}

	ctx = tflog.SetField(ctx, "ip", ip.Address)
	tflog.Info(ctx, "Allocated Instance IP address")

	rdns := d.Get("rdns").(string)
	if rdns != "" {
		if _, err := client.UpdateIPAddress(ctx, ip.Address, linodego.IPAddressUpdateOptions{
			RDNS: &rdns,
		}); err != nil {
			return diag.Errorf("failed to set RDNS for instance (%d) ip (%s): %s", linodeID, ip.Address, err)
		}

		tflog.Info(ctx, "Updated RDNS for IP address", map[string]any{
			"rdns": rdns,
		})
	}

	d.SetId(ip.Address)

	// Only reboot the associated instance with apply_immediately == true
	if applyImmediately {
		// TODO(jxriddle): check instance status and react to state

		// setting bootConfig to 0 to use current bootConfig
		if diagErr := helper.RebootInstance(ctx, d, linodeID, meta, 0); diagErr != nil {
			return diagErr
		}
		tflog.Info(ctx, "Rebooting instance for apply_immediately")
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_instance_ip")

	linodeID := d.Get("linode_id").(int)
	address := d.Id()
	rdns := d.Get("rdns").(string)

	client := meta.(*helper.ProviderMeta).Client

	if d.HasChange("rdns") {
		updateOptions := linodego.IPAddressUpdateOptions{}
		if rdns != "" {
			updateOptions.RDNS = &rdns
		}

		tflog.Info(ctx, "Updating RDNS for IP", map[string]any{
			"rdns": rdns,
		})
		if _, err := client.UpdateIPAddress(ctx, address, linodego.IPAddressUpdateOptions{
			RDNS: &rdns,
		}); err != nil {
			return diag.Errorf("failed to update RDNS for instance (%d) ip (%s): %s", linodeID, address, err)
		}
	}
	return nil
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_instance_ip")

	address := d.Id()
	linodeID := d.Get("linode_id").(int)

	client := meta.(*helper.ProviderMeta).Client

	if err := client.DeleteInstanceIPAddress(ctx, linodeID, address); err != nil {
		return diag.Errorf("failed to delete instance (%d) ip (%s): %s", linodeID, address, err)
	}

	tflog.Info(ctx, "Deleted Instance IP")

	return nil
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"linode_id": d.Get("linode_id").(int),
		"id":        d.Id(),
	})
}

func flattenVPCNAT1To1(data *linodego.InstanceIPNAT1To1) map[string]any {
	if data == nil {
		return nil
	}

	return map[string]any{
		"address":   data.Address,
		"vpc_id":    data.VPCID,
		"subnet_id": data.SubnetID,
	}
}
