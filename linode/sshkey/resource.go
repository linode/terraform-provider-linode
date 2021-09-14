package sshkey

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

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
		return diag.Errorf("Error parsing Linode SSH Key ID %s as int: %s", d.Id(), err)
	}

	sshkey, err := client.GetSSHKey(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode SSH Key ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode SSH Key: %s", err)
	}

	d.Set("label", sshkey.Label)
	d.Set("ssh_key", sshkey.SSHKey)
	if sshkey.Created != nil {
		d.Set("created", sshkey.Created.Format(time.RFC3339))
	}

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.SSHKeyCreateOptions{
		Label:  d.Get("label").(string),
		SSHKey: d.Get("ssh_key").(string),
	}
	sshkey, err := client.CreateSSHKey(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode SSH Key: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", sshkey.ID))
	d.Set("label", sshkey.Label)
	d.Set("ssh_key", sshkey.SSHKey)
	if sshkey.Created != nil {
		d.Set("created", sshkey.Created.Format(time.RFC3339))
	}

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode SSH Key id %s as int: %s", d.Id(), err)
	}

	if d.HasChange("label") {
		sshkey, err := client.GetSSHKey(ctx, int(id))

		updateOpts := sshkey.GetUpdateOptions()
		updateOpts.Label = d.Get("label").(string)

		if err != nil {
			return diag.Errorf("Error fetching data about the current Linode SSH Key: %s", err)
		}

		if sshkey, err = client.UpdateSSHKey(ctx, int(id), updateOpts); err != nil {
			return diag.FromErr(err)
		}
		d.Set("label", sshkey.Label)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode SSH Key id %s as int", d.Id())
	}
	err = client.DeleteSSHKey(ctx, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode SSH Key %d: %s", id, err)
	}
	return nil
}
