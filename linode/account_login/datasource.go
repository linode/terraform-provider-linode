package account_login

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

	loginInfo, err := client.GetLogin(ctx, id)
	if err != nil {
		return diag.Errorf("Error getting login %d: %s", id, err)
	}

	d.Set("datetime", loginInfo.Datetime.Format(time.RFC3339))
	d.Set("id", loginInfo.ID)
	d.Set("ip", loginInfo.IP)
	d.Set("restricted", loginInfo.Restricted)
	d.Set("username", loginInfo.Username)

	d.SetId(strconv.Itoa(loginInfo.ID))
	return nil
}
