package stackscript

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkDataSourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The StackScript's unique ID.",
			Required:    true,
		},

		"label": schema.StringAttribute{
			Description: "The StackScript's label is for display purposes only.",
			Computed:    true,
		},
		"script": schema.StringAttribute{
			Description: "The script to execute when provisioning a new Linode with this StackScript.",
			Computed:    true,
		},
		"description": schema.StringAttribute{
			Description: "A description for the StackScript.",
			Computed:    true,
		},
		"rev_note": schema.StringAttribute{
			Description: "This field allows you to add notes for the set of revisions made to this StackScript.",
			Computed:    true,
		},
		"is_public": schema.BoolAttribute{
			Description: "This determines whether other users can use your StackScript. Once a StackScript is " +
				"made public, it cannot be made private.",
			Computed: true,
		},
		"images": schema.SetAttribute{
			Description: "An array of Image IDs representing the Images that this StackScript is compatible for " +
				"deploying with.",
			ElementType: types.StringType,
			Computed:    true,
		},
		"deployments_active": schema.Int64Attribute{
			Description: "Count of currently active, deployed Linodes created from this StackScript.",
			Computed:    true,
		},
		"user_gravatar_id": schema.StringAttribute{
			Description: "The Gravatar ID for the User who created the StackScript.",
			Computed:    true,
		},
		"deployments_total": schema.Int64Attribute{
			Description: "The total number of times this StackScript has been deployed.",
			Computed:    true,
		},
		"username": schema.StringAttribute{
			Description: "The User who created the StackScript.",
			Computed:    true,
		},
		"created": schema.StringAttribute{
			Description: "The date this StackScript was created.",
			Computed:    true,
		},
		"updated": schema.StringAttribute{
			Description: "The date this StackScript was updated.",
			Computed:    true,
		},
		"user_defined_fields": schema.ListAttribute{
			Description: "This is a list of fields defined with a special syntax inside this " +
				"StackScript that allow for supplying customized parameters during deployment.",
			Computed:    true,
			ElementType: udfObjectType,
		},
	},
}
