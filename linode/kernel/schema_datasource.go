package kernel

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
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
}
