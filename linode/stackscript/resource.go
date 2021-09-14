package stackscript

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

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

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript ID %s as int: %s", d.Id(), err)
	}

	stackscript, err := client.GetStackscript(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 401 {
			log.Printf("[WARN] removing StackScript ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Stackscript: %s", err)
	}

	d.Set("label", stackscript.Label)
	d.Set("script", stackscript.Script)
	d.Set("description", stackscript.Description)
	d.Set("is_public", stackscript.IsPublic)
	d.Set("images", stackscript.Images)
	d.Set("rev_note", stackscript.RevNote)

	// Computed
	d.Set("deployments_active", stackscript.DeploymentsActive)
	d.Set("deployments_total", stackscript.DeploymentsTotal)
	d.Set("username", stackscript.Username)
	d.Set("user_gravatar_id", stackscript.UserGravatarID)
	d.Set("created", stackscript.Created.String())
	d.Set("updated", stackscript.Updated.String())
	setStackScriptUserDefinedFields(d, stackscript)
	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.StackscriptCreateOptions{
		Label:       d.Get("label").(string),
		Script:      d.Get("script").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
		RevNote:     d.Get("rev_note").(string),
	}

	for _, image := range d.Get("images").([]interface{}) {
		createOpts.Images = append(createOpts.Images, image.(string))
	}

	stackscript, err := client.CreateStackscript(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Stackscript: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", stackscript.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript id %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.StackscriptUpdateOptions{
		Label:       d.Get("label").(string),
		Script:      d.Get("script").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
		RevNote:     d.Get("rev_note").(string),
	}

	for _, image := range d.Get("images").([]interface{}) {
		updateOpts.Images = append(updateOpts.Images, image.(string))
	}

	if _, err = client.UpdateStackscript(ctx, int(id), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode Stackscript %d: %s", int(id), err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript id %s as int", d.Id())
	}
	err = client.DeleteStackscript(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return nil
		}
		return diag.Errorf("Error deleting Linode Stackscript %d: %s", id, err)
	}
	return nil
}
