package stackscript

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id := d.Get("id").(int)

	ss, err := client.GetStackscript(ctx, id)
	if err != nil {
		return diag.Errorf("Error getting Staakscript: %s", err)
	}

	if ss != nil {
		d.SetId(strconv.Itoa(id))
		d.Set("label", ss.Label)
		d.Set("script", ss.Script)
		d.Set("description", ss.Description)
		d.Set("rev_note", ss.RevNote)
		d.Set("is_public", ss.IsPublic)
		d.Set("images", ss.Images)
		d.Set("user_gravatar_id", ss.UserGravatarID)
		d.Set("deployments_active", ss.DeploymentsActive)
		d.Set("deployments_total", ss.DeploymentsTotal)
		d.Set("username", ss.Username)
		d.Set("created", ss.Created.Format(time.RFC3339))
		d.Set("updated", ss.Created.Format(time.RFC3339))
		setStackScriptUserDefinedFields(d, ss)
		return nil
	}

	return diag.Errorf("StackScript %d not found", id)
}
