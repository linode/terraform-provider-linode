package user

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var resourceLinodeUserGrantFields = []string{
	"global_grants", "domain_grant", "firewall_grant", "image_grant",
	"linode_grant", "longview_grant", "nodebalancer_grant", "stackscript_grant", "volume_grant",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
	}
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

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

	return readResource(ctx, d, meta)
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

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
		d.Set("firewall_grant", flattenGrantsEntities(grants.Firewall))
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

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

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
	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	username := d.Get("username").(string)
	if err := client.DeleteUser(ctx, username); err != nil {
		return diag.Errorf("failed to delete user (%s): %s", username, err)
	}
	return nil
}

func updateUserGrants(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*helper.ProviderMeta).Client

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
	updateOpts.Firewall = expandGrantsEntities(d.Get("firewall_grant").(*schema.Set).List())
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
	result.AddDatabases = global["add_databases"].(bool)
	result.AddFirewalls = global["add_firewalls"].(bool)
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

func userHasGrantsConfigured(d *schema.ResourceData) bool {
	for _, key := range resourceLinodeUserGrantFields {
		if _, ok := d.GetOk(key); ok {
			return true
		}
	}

	return false
}
