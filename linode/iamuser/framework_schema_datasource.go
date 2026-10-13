package iamuser

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"username": schema.StringAttribute{
			Description: "The username to work with.",
			Required:    true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"account_access": schema.ListAttribute{
			Description: "The user account level access.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"entity_access": schema.ListNestedAttribute{
			Description: "The user entity level access.",
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the entity.",
						Computed:    true,
					},
					"type": schema.StringAttribute{
						Description: "The entity category.",
						Computed:    true,
					},
					"roles": schema.ListAttribute{
						Description: "The roles for this entity and user.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
		},
	},
}
