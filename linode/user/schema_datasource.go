package user

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"username": {
		Type: schema.TypeString,
		Description: "This User's username. This is used for logging in, and may also be displayed alongside " +
			"actions the User performs (for example, in Events or public StackScripts).",
		Required: true,
	},
	"ssh_keys": {
		Type: schema.TypeList,
		Elem: &schema.Schema{Type: schema.TypeString},
		Description: "A list of SSH Key labels added by this User. These are the keys that will be deployed " +
			"if this User is included in the authorized_users field of a create Linode, rebuild Linode, or " +
			"create Disk request.",
		Computed: true,
	},
	"email": {
		Type: schema.TypeString,
		Description: "The email address for this User, for account management communications, and may be used " +
			"for other communications as configured.",
		Computed: true,
	},
	"restricted": {
		Type:        schema.TypeBool,
		Description: "If true, this User must be granted access to perform actions or access entities on this Account.",
		Computed:    true,
	},
	"global_grants": {
		Type:        schema.TypeList,
		Description: "A structure containing the Account-level grants a User has.",
		Computed:    true,
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
