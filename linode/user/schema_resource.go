package user

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var resourceSchema = map[string]*schema.Schema{
	"email": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The email of the user.",
	},
	"username": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The username of the user.",
	},
	"restricted": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, the user must be explicitly granted access to platform actions and entities.",
	},
	"ssh_keys": {
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Computed:    true,
		Description: "SSH keys to add to the user profile.",
	},
	"tfa_enabled": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "If the User has Two Factor Authentication (TFA) enabled.",
	},
	"global_grants": {
		Type:        schema.TypeList,
		Description: "A structure containing the Account-level grants a User has.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem:        resourceLinodeUserGrantsGlobal(),
	},
	"domain_grant":       resourceLinodeUserGrantsEntitySet(),
	"firewall_grant":     resourceLinodeUserGrantsEntitySet(),
	"image_grant":        resourceLinodeUserGrantsEntitySet(),
	"linode_grant":       resourceLinodeUserGrantsEntitySet(),
	"longview_grant":     resourceLinodeUserGrantsEntitySet(),
	"nodebalancer_grant": resourceLinodeUserGrantsEntitySet(),
	"stackscript_grant":  resourceLinodeUserGrantsEntitySet(),
	"volume_grant":       resourceLinodeUserGrantsEntitySet(),
}

func resourceLinodeUserGrantsGlobal() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_access": {
				Type: schema.TypeString,
				Description: "The level of access this User has to Account-level actions, like billing information. " +
					"A restricted User will never be able to manage users.",
				Optional: true,
			},
			"add_domains": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add Domains.",
				Optional:    true,
				Default:     false,
			},
			"add_firewalls": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add Firewalls.",
				Optional:    true,
				Default:     false,
			},
			"add_images": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add Images.",
				Optional:    true,
				Default:     false,
			},
			"add_linodes": {
				Type:        schema.TypeBool,
				Description: "If true, this User may create Linodes.",
				Optional:    true,
				Default:     false,
			},
			"add_longview": {
				Type:        schema.TypeBool,
				Description: "If true, this User may create Longview clients and view the current plan.",
				Optional:    true,
				Default:     false,
			},
			"add_nodebalancers": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add NodeBalancers.",
				Optional:    true,
				Default:     false,
			},
			"add_stackscripts": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add StackScripts.",
				Optional:    true,
				Default:     false,
			},
			"add_volumes": {
				Type:        schema.TypeBool,
				Description: "If true, this User may add Volumes.",
				Optional:    true,
				Default:     false,
			},
			"cancel_account": {
				Type:        schema.TypeBool,
				Description: "If true, this User may cancel the entire Account.",
				Optional:    true,
				Default:     false,
			},
			"longview_subscription": {
				Type:        schema.TypeBool,
				Description: "If true, this User may manage the Account’s Longview subscription.",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceLinodeUserGrantsEntity() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the entity this grant applies to.",
			},
			"permissions": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The level of access this User has to this entity. If null, this User has no access.",
			},
		},
	}
}

func resourceLinodeUserGrantsEntitySet() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeSet,
		Description: "A set containing all of the user's active grants.",
		Optional:    true,
		Computed:    true,
		Elem:        resourceLinodeUserGrantsEntity(),
	}
}
