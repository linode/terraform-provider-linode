package kernel

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"

	"context"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id := d.Get("id").(string)

	kernel, err := client.GetKernel(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get kernel: %s", err)
	}

	d.Set("architecture", kernel.Architecture)
	d.Set("deprecated", kernel.Deprecated)
	d.Set("kvm", kernel.KVM)
	d.Set("label", kernel.Label)
	d.Set("pvops", kernel.PVOPS)
	d.Set("version", kernel.Version)
	d.Set("xen", kernel.XEN)

	d.SetId(id)

	return nil
}
