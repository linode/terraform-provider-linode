package token

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
		return diag.Errorf("Error parsing Linode Token ID %s as int: %s", d.Id(), err)
	}

	token, err := client.GetToken(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode Token ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Token: %s", err)
	}

	d.Set("label", token.Label)
	d.Set("scopes", token.Scopes)
	d.Set("created", token.Created.Format(time.RFC3339))
	d.Set("expiry", token.Expiry.Format(time.RFC3339))

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.TokenCreateOptions{
		Label:  d.Get("label").(string),
		Scopes: d.Get("scopes").(string),
	}

	if expiryRaw, ok := d.GetOk("expiry"); ok {
		if expiry, ok := expiryRaw.(string); !ok {
			return diag.Errorf("expected expiry to be a string, got %s", expiryRaw)
		} else if dt, err := time.Parse("2006-01-02T15:04:05Z", expiry); err != nil {
			return diag.Errorf("expected expiry to be a datetime, got %s", expiry)
		} else {
			createOpts.Expiry = &dt
		}
	}

	token, err := client.CreateToken(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Token: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", token.ID))
	d.Set("token", token.Token)

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Token id %s as int: %s", d.Id(), err)
	}

	token, err := client.GetToken(ctx, int(id))
	if err != nil {
		return diag.Errorf("Error fetching data about the current linode: %s", err)
	}

	updateOpts := token.GetUpdateOptions()
	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)

		if token, err = client.UpdateToken(ctx, token.ID, updateOpts); err != nil {
			return diag.FromErr(err)
		}
		d.Set("label", token.Label)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Token id %s as int", d.Id())
	}
	err = client.DeleteToken(ctx, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode Token %d: %s", id, err)
	}
	// a settling cooldown to avoid expired tokens from being returned in listings
	time.Sleep(3 * time.Second)
	return nil
}
