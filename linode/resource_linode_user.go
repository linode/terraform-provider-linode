package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"

	"context"
	"fmt"
	"log"
)

var resourceLinodeUserGrantFields = []string{"global_grants", "domain_grant", "image_grant", "linode_grant",
	"longview_grant", "nodebalancer_grant", "stackscript_grant", "volume_grant"}

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

func resourceLinodeUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeUserCreate,
		ReadContext:   resourceLinodeUserRead,
		UpdateContext: resourceLinodeUserUpdate,
		DeleteContext: resourceLinodeUserDelete,
		Schema: map[string]*schema.Schema{
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
			"image_grant":        resourceLinodeUserGrantsEntitySet(),
			"linode_grant":       resourceLinodeUserGrantsEntitySet(),
			"longview_grant":     resourceLinodeUserGrantsEntitySet(),
			"nodebalancer_grant": resourceLinodeUserGrantsEntitySet(),
			"stackscript_grant":  resourceLinodeUserGrantsEntitySet(),
			"volume_grant":       resourceLinodeUserGrantsEntitySet(),
		},
	}
}

func resourceLinodeUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	createOpts := linodego.UserCreateOptions{
		Email:      d.Get("email").(string),
		Username:   d.Get("username").(string),
		Restricted: d.Get("restricted").(bool),
	}

	user, err := client.CreateUser(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create user: %s", err)
	}

	if userHasGrantsConfigured(d) {
		if err := updateUserGrants(ctx, d, meta); err != nil {
			return diag.Errorf("failed to set user grants (%s): %s", user.Username, err)
		}
	}

	d.SetId(user.Username)

	return resourceLinodeUserRead(ctx, d, meta)
}

func resourceLinodeUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	username := d.Get("username").(string)
	user, err := client.GetUser(ctx, username)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode User %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get user (%s): %s", username, err)
	}

	if user.Restricted {
		grants, err := client.GetUserGrants(ctx, username)
		if err != nil {
			return diag.Errorf("failed to get user grants (%s): %s", username, err)
		}

		d.Set("global_grants", []interface{}{flattenGrantsGlobal(&grants.Global)})

		d.Set("domain_grant", flattenGrantsEntities(grants.Domain))
		d.Set("image_grant", flattenGrantsEntities(grants.Image))
		d.Set("linode_grant", flattenGrantsEntities(grants.Linode))
		d.Set("longview_grant", flattenGrantsEntities(grants.Longview))
		d.Set("nodebalancer_grant", flattenGrantsEntities(grants.NodeBalancer))
		d.Set("stackscript_grant", flattenGrantsEntities(grants.StackScript))
		d.Set("volume_grant", flattenGrantsEntities(grants.Volume))
	}

	d.Set("username", username)
	d.Set("email", user.Email)
	d.Set("restricted", user.Restricted)
	d.Set("ssh_keys", user.SSHKeys)
	d.Set("tfa_enabled", user.TFAEnabled)
	return nil
}

func resourceLinodeUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	id := d.Id()
	username := d.Get("username").(string)
	restricted := d.Get("restricted").(bool)

	if _, err := client.UpdateUser(ctx, id, linodego.UserUpdateOptions{
		Username:   username,
		Restricted: &restricted,
	}); err != nil {
		return diag.Errorf("failed to update user (%s): %s", id, err)
	}

	if d.HasChanges(resourceLinodeUserGrantFields...) {
		if err := updateUserGrants(ctx, d, meta); err != nil {
			return diag.Errorf("failed to update user grants (%s): %s", id, err)
		}
	}

	d.SetId(username)
	return resourceLinodeUserRead(ctx, d, meta)
}

func resourceLinodeUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	username := d.Get("username").(string)
	if err := client.DeleteUser(ctx, username); err != nil {
		return diag.Errorf("failed to delete user (%s): %s", username, err)
	}
	return nil
}

func updateUserGrants(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client

	username := d.Get("username").(string)
	restricted := d.Get("restricted").(bool)

	if !restricted {
		return fmt.Errorf("user must be restricted in order to update grants")
	}

	updateOpts := linodego.UserGrantsUpdateOptions{}

	if global, ok := d.GetOk("global_grants"); ok {
		global := global.([]interface{})[0].(map[string]interface{})
		updateOpts.Global = expandGrantsGlobal(global)
	}

	updateOpts.Domain = expandGrantsEntities(d.Get("domain_grant").(*schema.Set).List())
	updateOpts.Image = expandGrantsEntities(d.Get("image_grant").(*schema.Set).List())
	updateOpts.Linode = expandGrantsEntities(d.Get("linode_grant").(*schema.Set).List())
	updateOpts.Longview = expandGrantsEntities(d.Get("longview_grant").(*schema.Set).List())
	updateOpts.NodeBalancer = expandGrantsEntities(d.Get("nodebalancer_grant").(*schema.Set).List())
	updateOpts.StackScript = expandGrantsEntities(d.Get("stackscript_grant").(*schema.Set).List())
	updateOpts.Volume = expandGrantsEntities(d.Get("volume_grant").(*schema.Set).List())

	if _, err := client.UpdateUserGrants(ctx, username, updateOpts); err != nil {
		return err
	}

	return nil
}

func expandGrantsEntities(entities []interface{}) []linodego.EntityUserGrant {
	result := make([]linodego.EntityUserGrant, len(entities))

	for i, entity := range entities {
		entity := entity.(map[string]interface{})
		result[i] = expandGrantsEntity(entity)
	}

	return result
}

func expandGrantsEntity(entity map[string]interface{}) linodego.EntityUserGrant {
	result := linodego.EntityUserGrant{}

	permissions := linodego.GrantPermissionLevel(entity["permissions"].(string))

	result.ID = entity["id"].(int)
	result.Permissions = &permissions

	return result
}

func expandGrantsGlobal(global map[string]interface{}) linodego.GlobalUserGrants {
	result := linodego.GlobalUserGrants{}

	result.AccountAccess = nil

	if accountAccess, ok := global["account_access"].(string); ok && accountAccess != "" {
		accountAccess := linodego.GrantPermissionLevel(accountAccess)

		result.AccountAccess = &accountAccess
	}

	result.AddDomains = global["add_domains"].(bool)
	result.AddImages = global["add_images"].(bool)
	result.AddLinodes = global["add_linodes"].(bool)
	result.AddLongview = global["add_longview"].(bool)
	result.AddNodeBalancers = global["add_nodebalancers"].(bool)
	result.AddStackScripts = global["add_stackscripts"].(bool)
	result.AddVolumes = global["add_volumes"].(bool)
	result.CancelAccount = global["cancel_account"].(bool)
	result.LongviewSubscription = global["longview_subscription"].(bool)

	return result
}

func flattenGrantsEntities(entities []linodego.GrantedEntity) []interface{} {
	var result []interface{}

	for _, entity := range entities {
		// Filter out entities without any permissions set.
		// This is necessary because Linode will automatically
		// create empty entities that will trigger false diffs.
		if entity.Permissions == "" {
			continue
		}

		result = append(result, flattenGrantsEntity(&entity))
	}

	return result
}

func flattenGrantsEntity(entity *linodego.GrantedEntity) map[string]interface{} {
	result := make(map[string]interface{})

	result["id"] = entity.ID
	result["permissions"] = entity.Permissions

	return result
}

func flattenGrantsGlobal(global *linodego.GlobalUserGrants) map[string]interface{} {
	result := make(map[string]interface{})

	result["account_access"] = global.AccountAccess
	result["add_domains"] = global.AddDomains
	result["add_images"] = global.AddImages
	result["add_linodes"] = global.AddLinodes
	result["add_longview"] = global.AddLongview
	result["add_nodebalancers"] = global.AddNodeBalancers
	result["add_stackscripts"] = global.AddStackScripts
	result["add_volumes"] = global.AddVolumes
	result["cancel_account"] = global.CancelAccount
	result["longview_subscription"] = global.LongviewSubscription

	return result
}

func userHasGrantsConfigured(d *schema.ResourceData) bool {
	for _, key := range resourceLinodeUserGrantFields {
		if _, ok := d.GetOk(key); ok {
			return true
		}
	}

	return false
}
