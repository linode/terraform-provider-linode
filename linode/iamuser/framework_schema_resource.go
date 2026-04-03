package iamuser

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var EntityAccessType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":    types.Int64Type,
		"type":  types.StringType,
		"roles": types.ListType{ElemType: types.StringType},
	},
}

var frameworkResourceSchema = schema.Schema{
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
			Optional:    true,
			Computed:    true,
			ElementType: types.StringType,
		},
		"entity_access": schema.ListNestedAttribute{
			Description: "The user entity level access.",
			Optional:    true,
			Computed:    true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "The ID of the entity.",
						Required:    true,
					},
					"type": schema.StringAttribute{
						Description: "The entity category.",
						Required:    true,
					},
					"roles": schema.ListAttribute{
						Description: "The roles for this entity and user.",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
		},
	},
}
