package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
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

	reqUsername := d.Get("username").(string)

	if reqUsername == "" {
		return diag.Errorf("Error User username is required")
	}

	users, err := client.ListUsers(ctx, nil)
	var user linodego.User
	if err != nil {
		return diag.Errorf("Error listing user: %s", err)
	}

	for _, testuser := range users {
		if testuser.Username == reqUsername {
			user = testuser
			break
		}
	}

	if user.Username != "" {
		d.SetId(fmt.Sprintf("%s", user.Username))
		d.Set("username", user.Username)
		d.Set("email", user.Email)
		d.Set("ssh_keys", user.SSHKeys)
		d.Set("restricted", user.Restricted)

		return nil
	}

	return diag.Errorf("Linode User with username %s was not found", reqUsername)
}
