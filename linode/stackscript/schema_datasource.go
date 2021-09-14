package stackscript

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of this Stackscript.",
		Required:    true,
	},
	"label": {
		Type:        schema.TypeString,
		Description: "The StackScript's label is for display purposes only.",
		Computed:    true,
	},
	"script": {
		Type:        schema.TypeString,
		Description: "The script to execute when provisioning a new Linode with this StackScript.",
		Computed:    true,
	},
	"description": {
		Type:        schema.TypeString,
		Description: "A description for the StackScript.",
		Computed:    true,
	},
	"rev_note": {
		Type:        schema.TypeString,
		Description: "This field allows you to add notes for the set of revisions made to this StackScript.",
		Computed:    true,
	},
	"is_public": {
		Type: schema.TypeBool,
		Description: "This determines whether other users can use your StackScript. Once a StackScript is " +
			"made public, it cannot be made private.",
		Computed: true,
	},
	"images": {
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
		Description: "An array of Image IDs representing the Images that this StackScript is compatible for " +
			"deploying with.",
		Computed: true,
	},

	"deployments_active": {
		Type:        schema.TypeInt,
		Description: "Count of currently active, deployed Linodes created from this StackScript.",
		Computed:    true,
	},
	"user_gravatar_id": {
		Type:        schema.TypeString,
		Description: "The Gravatar ID for the User who created the StackScript.",
		Computed:    true,
	},
	"deployments_total": {
		Type:        schema.TypeInt,
		Description: "The total number of times this StackScript has been deployed.",
		Computed:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Description: "The User who created the StackScript.",
		Computed:    true,
	},
	"created": {
		Type:        schema.TypeString,
		Description: "The date this StackScript was created.",
		Computed:    true,
	},
	"updated": {
		Type:        schema.TypeString,
		Description: "The date this StackScript was updated.",
		Computed:    true,
	},
	"user_defined_fields": {
		Description: "This is a list of fields defined with a special syntax inside this StackScript that " +
			"allow for supplying customized parameters during deployment.",
		Type:       schema.TypeList,
		Computed:   true,
		Optional:   true,
		ConfigMode: schema.SchemaConfigModeAttr,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"label": {
					Type: schema.TypeString,
					Description: "A human-readable label for the field that will serve as the " +
						"input prompt for entering the value during deployment.",
					Computed: true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the field.",
					Computed:    true,
				},
				"example": {
					Type:        schema.TypeString,
					Description: "An example value for the field.",
					Computed:    true,
				},
				"one_of": {
					Type:        schema.TypeString,
					Description: "A list of acceptable single values for the field.",
					Computed:    true,
				},
				"many_of": {
					Type:        schema.TypeString,
					Description: "A list of acceptable values for the field in any quantity, combination or order.",
					Computed:    true,
				},
				"default": {
					Type:        schema.TypeString,
					Description: "The default value. If not specified, this value will be used.",
					Computed:    true,
				},
			},
		},
	},
}
