package sshkey

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/linode/terraform-provider-linode/linode/helper/customtypes"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"label": schema.StringAttribute{
			Description: "The label of the Linode SSH Key.",
			Required:    true,
		},
		"ssh_key": schema.StringAttribute{
			Description: "The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"created": schema.StringAttribute{
			Description: "The date this key was added.",
			Computed:    true,
			CustomType:  customtypes.RFC3339TimeStringType{},
		},
		// This field must be a string attribute for pass-through
		// imports to work as expected.
		"id": schema.StringAttribute{
			Description: "The unique identifier for this SSH key.",
			Computed:    true,
		},
	},
}
