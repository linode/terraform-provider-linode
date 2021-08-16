package account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataSource,
		Schema:      dataSourceSchema,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	account, err := client.GetAccount(ctx)
	if err != nil {
		return diag.Errorf("Error getting account: %s", err)
	}

	d.SetId(account.Email)
	d.Set("email", account.Email)
	d.Set("first_name", account.FirstName)
	d.Set("last_name", account.LastName)
	d.Set("company", account.Company)

	d.Set("address_1", account.Address1)
	d.Set("address_2", account.Address2)
	d.Set("phone", account.Phone)
	d.Set("city", account.City)
	d.Set("state", account.State)
	d.Set("country", account.Country)
	d.Set("zip", account.Zip)

	d.Set("balance", account.Balance)

	// We exclude the credit_card and tax_id fields because they are too sensitive

	return nil
}
