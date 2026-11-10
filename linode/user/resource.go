package user

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-set/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

var resourceLinodeUserGrantFields = []string{
	"global_grants", "domain_grant", "firewall_grant", "image_grant",
	"linode_grant", "longview_grant", "nodebalancer_grant", "stackscript_grant",
	"vpc_grant", "volume_grant",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func createResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tflog.Debug(ctx, "Create linode_user")

	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.UserCreateOptions{
		Email:      d.Get("email").(string),
		Username:   d.Get("username").(string),
		Restricted: d.Get("restricted").(bool),
	}

	tflog.Debug(ctx, "client.CreateUser(...)", map[string]any{
		"options": createOpts,
	})
	user, err := client.CreateUser(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create user: %s", err)
	}

	d.SetId(user.Username)

	ctx = populateLogAttributes(ctx, d)

	if userHasGrantsConfigured(d) {
		if err := updateUserGrants(ctx, d, meta); err != nil {
			return diag.Errorf("failed to set user grants (%s): %s", user.Username, err)
		}
	}

	return readResource(ctx, d, meta)
}

func readResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_user")

	client := meta.(*helper.ProviderMeta).Client

	username := d.Id()

	user, err := client.GetUser(ctx, username)
	if err != nil {
		if linodego.IsNotFound(err) {
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

		d.Set("global_grants", []any{flattenGrantsGlobal(&grants.Global)})

		d.Set("domain_grant", flattenGrantsEntities(grants.Domain))
		d.Set("firewall_grant", flattenGrantsEntities(grants.Firewall))
		d.Set("image_grant", flattenGrantsEntities(grants.Image))
		d.Set("linode_grant", flattenGrantsEntities(grants.Linode))
		d.Set("longview_grant", flattenGrantsEntities(grants.Longview))
		d.Set("nodebalancer_grant", flattenGrantsEntities(grants.NodeBalancer))
		d.Set("stackscript_grant", flattenGrantsEntities(grants.StackScript))
		d.Set("volume_grant", flattenGrantsEntities(grants.Volume))
		d.Set("vpc_grant", flattenGrantsEntities(grants.VPC))
	}

	d.Set("username", username)
	d.Set("email", user.Email)
	d.Set("restricted", user.Restricted)
	d.Set("user_type", user.UserType)
	d.Set("ssh_keys", user.SSHKeys)
	d.Set("tfa_enabled", user.TFAEnabled)
	return nil
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_user")

	client := meta.(*helper.ProviderMeta).Client

	id := d.Id()
	username := d.Get("username").(string)
	restricted := d.Get("restricted").(bool)

	updateOpts := linodego.UserUpdateOptions{
		Username:   username,
		Restricted: &restricted,
	}

	tflog.Debug(ctx, "client.UpdateUser(...)", map[string]any{
		"options": updateOpts,
	})

	if _, err := client.UpdateUser(ctx, id, updateOpts); err != nil {
		return diag.Errorf("failed to update user (%s): %s", id, err)
	}

	d.SetId(username)

	if d.HasChanges(resourceLinodeUserGrantFields...) {
		if err := updateUserGrants(ctx, d, meta); err != nil {
			return diag.Errorf("failed to update user grants (%s): %s", username, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_user")

	client := meta.(*helper.ProviderMeta).Client

	username := d.Id()

	tflog.Debug(ctx, "client.DeleteUser(...)")
	if err := client.DeleteUser(ctx, username); err != nil {
		return diag.Errorf("failed to delete user (%s): %s", username, err)
	}
	return nil
}

func updateUserGrants(ctx context.Context, d *schema.ResourceData, meta any) error {
	client := meta.(*helper.ProviderMeta).Client

	username := d.Id()
	restricted := d.Get("restricted").(bool)

	// TODO: Implement this validation at plan-time
	if !restricted {
		return fmt.Errorf("user must be restricted in order to update grants")
	}

	updateOpts := linodego.UserGrantsUpdateOptions{}

	if global, ok := d.GetOk("global_grants"); ok {
		global := global.([]any)[0].(map[string]any)
		updateOpts.Global = expandGrantsGlobal(global)
	}

	updateOpts.Domain = expandGrantsEntityChanges(d, "domain_grant")
	updateOpts.Firewall = expandGrantsEntityChanges(d, "firewall_grant")
	updateOpts.Image = expandGrantsEntityChanges(d, "image_grant")
	updateOpts.Linode = expandGrantsEntityChanges(d, "linode_grant")
	updateOpts.Longview = expandGrantsEntityChanges(d, "longview_grant")
	updateOpts.NodeBalancer = expandGrantsEntityChanges(d, "nodebalancer_grant")
	updateOpts.StackScript = expandGrantsEntityChanges(d, "stackscript_grant")
	updateOpts.VPC = expandGrantsEntityChanges(d, "vpc_grant")
	updateOpts.Volume = expandGrantsEntityChanges(d, "volume_grant")

	tflog.Debug(ctx, "client.UpdateUserGrants(...)", map[string]any{
		"options": updateOpts,
	})

	if _, err := client.UpdateUserGrants(ctx, username, updateOpts); err != nil {
		return err
	}

	return nil
}

func expandGrantsEntities(entities []any) []linodego.EntityUserGrant {
	result := make([]linodego.EntityUserGrant, len(entities))

	for i, entity := range entities {
		entity := entity.(map[string]any)
		result[i] = expandGrantsEntity(entity)
	}

	return result
}

func expandGrantsEntityChanges(d *schema.ResourceData, field string) []linodego.EntityUserGrant {
	oldValue, newValue := d.GetChange(field)

	oldEntities := oldValue.(*schema.Set).List()
	newEntities := newValue.(*schema.Set).List()

	return expandGrantsEntitiesWithRemovals(oldEntities, newEntities)
}

func expandGrantsEntitiesWithRemovals(oldEntities, newEntities []any) []linodego.EntityUserGrant {
	result := expandGrantsEntities(newEntities)

	newIDs := set.New[int](len(result))
	for _, grant := range result {
		newIDs.Insert(grant.ID)
	}

	for _, entity := range oldEntities {
		id := entity.(map[string]any)["id"].(int)

		if !newIDs.Contains(id) {
			result = append(result, linodego.EntityUserGrant{
				ID:          id,
				Permissions: nil,
			})
		}
	}

	return result
}

func expandGrantsEntity(entity map[string]any) linodego.EntityUserGrant {
	result := linodego.EntityUserGrant{}

	permissions := linodego.GrantPermissionLevel(entity["permissions"].(string))

	result.ID = entity["id"].(int)
	result.Permissions = &permissions

	return result
}

func expandGrantsGlobal(global map[string]any) linodego.GlobalUserGrants {
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
	result.AddVPCs = global["add_vpcs"].(bool)
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

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return tflog.SetField(ctx, "username", d.Id())
}
