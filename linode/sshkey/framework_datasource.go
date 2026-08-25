package sshkey

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
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

func (data *DataSourceModel) ParseSSHKey(ssh *linodego.SSHKey) {
	data.ID = types.StringValue(strconv.Itoa(ssh.ID))
	data.Label = types.StringValue(ssh.Label)
	data.SSHKey = types.StringValue(ssh.SSHKey)
	data.Created = timetypes.NewRFC3339TimePointerValue(ssh.Created)
}

type DataSourceModel struct {
	Label   types.String      `tfsdk:"label"`
	SSHKey  types.String      `tfsdk:"ssh_key"`
	Created timetypes.RFC3339 `tfsdk:"created"`
	ID      types.String      `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasID := !data.ID.IsNull() && !data.ID.IsUnknown()
	hasLabel := !data.Label.IsNull() && !data.Label.IsUnknown()

	var sshkey *linodego.SSHKey
	var ddiag diag.Diagnostic

	switch {
	case hasID:
		id := helper.FrameworkSafeStringToInt(data.ID.ValueString(), &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		ctx = tflog.SetField(ctx, "sshkey_id", id)
		sshkey, ddiag = d.getSSHKeyByID(ctx, id)
	case hasLabel:
		ctx = tflog.SetField(ctx, "sshkey_label", data.Label.ValueString())
		sshkey, ddiag = d.getSSHKeyByLabel(ctx, data.Label.ValueString())
	default:
		resp.Diagnostics.AddError(
			"Missing required argument",
			"Either id or label must be specified.",
		)
		return
	}

	if ddiag != nil {
		resp.Diagnostics.Append(ddiag)
		return
	}

	data.ParseSSHKey(sshkey)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *DataSource) getSSHKeyByID(ctx context.Context, id int) (*linodego.SSHKey, diag.Diagnostic) {
	tflog.Trace(ctx, "client.GetSSHKey(...)")

	sshkey, err := d.Meta.Client.GetSSHKey(ctx, id)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed to get SSH Key with id %d", id),
			err.Error(),
		)
	}

	return sshkey, nil
}

func (d *DataSource) getSSHKeyByLabel(ctx context.Context, label string) (*linodego.SSHKey, diag.Diagnostic) {
	tflog.Trace(ctx, "client.ListSSHKeys(...)")

	keys, err := d.Meta.Client.ListSSHKeys(ctx, nil)
	if err != nil {
		return nil, diag.NewErrorDiagnostic(
			"Failed to list SSH Keys",
			err.Error(),
		)
	}

	for i := range keys {
		if keys[i].Label == label {
			return &keys[i], nil
		}
	}

	return nil, diag.NewErrorDiagnostic(
		"Failed to retrieve Linode SSH Key",
		fmt.Sprintf("SSH Key with label %s was not found", label),
	)
}
