package user

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var frameworkResourceSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"email": schema.StringAttribute{
			Description:   "The email of the user.",
			Required:      true,
			PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
		},
		"username": schema.StringAttribute{
			Description: "The username of the user.",
			Required:    true,
		},
		"restricted": schema.BoolAttribute{
			Description: "If true, the user must be explicitly granted access to platform actions and entities.",
			Optional:    true,
			Default:     booldefault.StaticBool(false),
		},
		"ssh_keys": schema.ListAttribute{
			Description: "A list of SSH Key labels added by this User. These are the keys that will be deployed " +
				"if this User is included in the authorized_users field of a create Linode, rebuild Linode, or " +
				"create Disk request.",
			Computed:    true,
			ElementType: types.StringType,
		},
		"tfa_enabled": schema.BoolAttribute{
			Description: "A boolean value indicating if the User has Two Factor Authentication (TFA) enabled.",
			Computed:    true,
		},
		"password_created": schema.StringAttribute{
			Description: "The date and time when this User’s current password was created." +
				"User passwords are first created during the Account sign-up process, " +
				"and updated using the Reset Password webpage." +
				"null if this User has not created a password yet.",
			Computed: true,
		},
		"verified_phone_number": schema.StringAttribute{
			Description: "The phone number verified for this User Profile with the Phone Number Verify command." +
				"null if this User Profile has no verified phone number.",
			Computed: true,
		},
		"id": schema.StringAttribute{
			Description: "Unique identifier for this Resource.",
			Computed:    true,
		},
	},
	Blocks: map[string]schema.Block{
		"global_grants": schema.ListNestedBlock{
			Description: "A structure containing the Account-level grants a User has.",
			Validators:  []validator.List{listvalidator.SizeAtMost(1)},
			NestedObject: schema.NestedBlockObject{
				Attributes: map[string]schema.Attribute{
					"account_access": schema.StringAttribute{
						Description: "The level of access this User has to Account-level actions, like billing information. " +
							"A restricted User will never be able to manage users.",
						Optional:   true,
						Computed:   true,
						Validators: []validator.String{stringvalidator.OneOf("read_only", "read_write")},
					},
					"add_databases": schema.BoolAttribute{
						Description: "If true, this User may add Databases.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_domains": schema.BoolAttribute{
						Description: "If true, this User may add Domains.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_firewalls": schema.BoolAttribute{
						Description: "If true, this User may add Firewalls.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_images": schema.BoolAttribute{
						Description: "If true, this User may add Images.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_linodes": schema.BoolAttribute{
						Description: "If true, this User may create Linodes.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_longview": schema.BoolAttribute{
						Description: "If true, this User may create Longview clients and view the current plan.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_nodebalancers": schema.BoolAttribute{
						Description: "If true, this User may add NodeBalancers.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_stackscripts": schema.BoolAttribute{
						Description: "If true, this User may add StackScripts.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"add_volumes": schema.BoolAttribute{
						Description: "If true, this User may add Volumes.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"cancel_account": schema.BoolAttribute{
						Description: "If true, this User may cancel the entire Account.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"longview_subscription": schema.BoolAttribute{
						Description: "If true, this User may manage the Account’s Longview subscription.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
		},
		"database_grant":     resourceLinodeUserGrantsEntitySet,
		"domain_grant":       resourceLinodeUserGrantsEntitySet,
		"firewall_grant":     resourceLinodeUserGrantsEntitySet,
		"image_grant":        resourceLinodeUserGrantsEntitySet,
		"linode_grant":       resourceLinodeUserGrantsEntitySet,
		"longview_grant":     resourceLinodeUserGrantsEntitySet,
		"nodebalancer_grant": resourceLinodeUserGrantsEntitySet,
		"stackscript_grant":  resourceLinodeUserGrantsEntitySet,
		"volume_grant":       resourceLinodeUserGrantsEntitySet,
	},
}

var resourceLinodeUserGrantsEntitySet = schema.SetNestedBlock{
	Description: "A set containing all of the user's active grants.",
	NestedObject: schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the entity this grant applies to.",
				Required:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "The level of access this User has to this entity. If null, this User has no access.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.OneOf("read_only", "read_write")},
			},
			"label": schema.StringAttribute{
				Description: "The current label of the entity this grant applies to, for display purposes.",
				Computed:    true,
			},
		},
	},
}
