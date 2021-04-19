package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"context"
)

func dataSourceLinodeKernel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinodeKernelRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "The unique ID of this Kernel.",
				Required:    true,
			},
			"architecture": {
				Type:        schema.TypeString,
				Description: "The architecture of this Kernel.",
				Computed:    true,
			},
			"deprecated": {
				Type:        schema.TypeBool,
				Description: "Whether or not this Kernel is deprecated.",
				Computed:    true,
			},
			"kvm": {
				Type:        schema.TypeBool,
				Description: "If this Kernel is suitable for KVM Linodes.",
				Computed:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "The friendly name of this Kernel.",
				Computed:    true,
			},
			"pvops": {
				Type:        schema.TypeBool,
				Description: "If this Kernel is suitable for paravirtualized operations.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Linux Kernel version.",
				Computed:    true,
			},
			"xen": {
				Type:        schema.TypeBool,
				Description: "If this Kernel is suitable for Xen Linodes.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeKernelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

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
