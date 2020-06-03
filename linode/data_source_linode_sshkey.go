package linode

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeSSHKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLinodeSSHKeyRead,

		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode SSH Key.",
				Required:    true,
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Description: "The public SSH Key, which is used to authenticate to the root user of the Linodes you deploy.",
				Computed:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "The date this key was added.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeSSHKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)

	reqLabel := d.Get("label").(string)

	sshkeys, err := client.ListSSHKeys(context.Background(), nil)
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
