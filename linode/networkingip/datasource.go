package networkingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	reqImage := d.Get("address").(string)

	if reqImage == "" {
		return diag.Errorf("NetworkingIP address is required")
	}

	address, err := client.GetIPAddress(ctx, reqImage)
	if err != nil {
		return diag.Errorf("Error listing addresses: %s", err)
	}

	if address != nil {
		d.SetId(address.Address)
		d.Set("address", address.Address)
		d.Set("gateway", address.Gateway)
		d.Set("subnet_mask", address.SubnetMask)
		d.Set("prefix", address.Prefix)
		d.Set("type", address.Type)
		d.Set("public", address.Public)
		d.Set("rdns", address.RDNS)
		d.Set("linode_id", address.LinodeID)
		d.Set("region", address.Region)
		return nil
	}

	d.SetId("")

	return diag.Errorf("NetworkingIP address %s was not found", reqImage)
}
