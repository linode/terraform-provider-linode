package sshkey

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_sshkey",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseSSHKey(ssh *linodego.SSHKey) diag.Diagnostics {
	diags := diag.Diagnostics{}

	if ssh.ID == 0 {
		diags.AddError(
			fmt.Sprintf("Linode SSH Key with label %s was not found", data.Label.ValueString()), "",
		)
		return diags
	}

	data.Label = types.StringValue(ssh.Label)
	data.SSHKey = types.StringValue(ssh.SSHKey)
	data.Created = types.StringValue(ssh.Created.Format(time.RFC3339))

	id, err := json.Marshal(ssh)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

type DataSourceModel struct {
	Label   types.String `tfsdk:"label"`
	SSHKey  types.String `tfsdk:"ssh_key"`
	Created types.String `tfsdk:"created"`
	ID      types.String `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keys, err := client.ListSSHKeys(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error listing SSH keys: %s", err.Error(),
		)
		return
	}

	var sshkey linodego.SSHKey

	for _, testkey := range keys {
		if testkey.Label == data.Label.ValueString() {
			sshkey = testkey
			break
		}
	}

	resp.Diagnostics.Append(data.parseSSHKey(&sshkey)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
