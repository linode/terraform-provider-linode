package lock

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/linode/linodego"
)

var frameworkResourceSchema = schema.Schema{
	Description: "Manages a Linode Lock. Locks prevent accidental deletion of resources.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique ID of the Lock.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"entity_id": schema.Int64Attribute{
			Description: "The ID of the entity to lock.",
			Required:    true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplace(),
			},
		},
		"entity_type": schema.StringAttribute{
			Description: "The type of the entity to lock. Currently only 'linode' is supported. Note: Linodes that are part of an LKE cluster cannot be locked.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.EntityLinode),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"lock_type": schema.StringAttribute{
			Description: "The type of lock. Only one lock type can exist per resource at a time. Valid values are 'cannot_delete' (prevents deletion, rebuild, and transfer) and 'cannot_delete_with_subresources' (prevents deletion, rebuild, transfer, and deletion of subresources such as disks, configs, interfaces, and IP addresses).",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(linodego.LockTypeCannotDelete),
					string(linodego.LockTypeCannotDeleteWithSubresources),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"entity_label": schema.StringAttribute{
			Description: "The label of the locked entity.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"entity_url": schema.StringAttribute{
			Description: "The URL of the locked entity.",
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	},
}
