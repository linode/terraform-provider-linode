package sshkey

import (
	"context"
	"fmt"
	"time"

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

	reqLabel := d.Get("label").(string)

	if reqLabel == "" {
		return diag.Errorf("Error SSH Key label is required")
	}

	sshkeys, err := client.ListSSHKeys(ctx, nil)
	var sshkey linodego.SSHKey
	if err != nil {
		return diag.Errorf("Error listing sshkey: %s", err)
	}

	for _, testkey := range sshkeys {
		if testkey.Label == reqLabel {
			sshkey = testkey
			break
		}
	}

	if sshkey.ID != 0 {
		d.SetId(fmt.Sprintf("%d", sshkey.ID))
		d.Set("label", sshkey.Label)
		d.Set("ssh_key", sshkey.SSHKey)
		if sshkey.Created != nil {
			d.Set("created", sshkey.Created.Format(time.RFC3339))
		}

		return nil
	}

	return diag.Errorf("Linode SSH Key with label %s was not found", reqLabel)
}
