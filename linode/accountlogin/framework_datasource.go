package accountlogin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DatasourceModel) parseAccountLogin(accountLogin *linodego.Login) {
	data.Datetime = types.StringValue(accountLogin.Datetime.Format(time.RFC3339))
	data.ID = types.Int64Value(int64(accountLogin.ID))
	data.IP = types.StringValue(accountLogin.IP)
	data.Restricted = types.BoolValue(accountLogin.Restricted)
	data.Username = types.StringValue(accountLogin.Username)
	data.Status = types.StringValue(accountLogin.Status)
}

func (d *DataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	meta := helper.GetDataSourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

type DatasourceModel struct {
	Datetime   types.String `tfsdk:"datetime"`
	ID         types.Int64  `tfsdk:"id"`
	IP         types.String `tfsdk:"ip"`
	Restricted types.Bool   `tfsdk:"restricted"`
	Username   types.String `tfsdk:"username"`
	Status     types.String `tfsdk:"status"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_account_login"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDataSourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_account_login")
	client := d.client

	var data DatasourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loginID := helper.FrameworkSafeInt64ToInt(data.ID.ValueInt64(), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Trace(ctx, "client.GetLogin(...)", map[string]any{
		"loginID": loginID,
	})
	accountlogin, err := client.GetLogin(ctx, loginID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Account Login",
			fmt.Sprintf(
				"Error finding Linode Account Login: %s",
				err.Error(),
			),
		)
		return
	}

	data.parseAccountLogin(accountlogin)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
