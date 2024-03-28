package account

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"

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

func (data *DataSourceModel) ParseAccount(account *linodego.Account) {
	data.EUUID = types.StringValue(account.EUUID)
	data.Email = types.StringValue(account.Email)
	data.FirstName = types.StringValue(account.FirstName)
	data.LastName = types.StringValue(account.LastName)
	data.Company = types.StringValue(account.Company)
	data.Address1 = types.StringValue(account.Address1)
	data.Address2 = types.StringValue(account.Address2)
	data.Phone = types.StringValue(account.Phone)
	data.City = types.StringValue(account.City)
	data.State = types.StringValue(account.State)
	data.Country = types.StringValue(account.Country)
	data.Zip = types.StringValue(account.Zip)
	data.Balance = types.Float64Value(float64(account.Balance))
	data.Capabilities = helper.StringSliceToFramework(account.Capabilities)
	data.ActiveSince = timetypes.NewRFC3339TimePointerValue(account.ActiveSince)
	data.ID = types.StringValue(account.Email)
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

type DataSourceModel struct {
	EUUID        types.String      `tfsdk:"euuid"`
	Email        types.String      `tfsdk:"email"`
	FirstName    types.String      `tfsdk:"first_name"`
	LastName     types.String      `tfsdk:"last_name"`
	Company      types.String      `tfsdk:"company"`
	Address1     types.String      `tfsdk:"address_1"`
	Address2     types.String      `tfsdk:"address_2"`
	Phone        types.String      `tfsdk:"phone"`
	City         types.String      `tfsdk:"city"`
	State        types.String      `tfsdk:"state"`
	Country      types.String      `tfsdk:"country"`
	Zip          types.String      `tfsdk:"zip"`
	Balance      types.Float64     `tfsdk:"balance"`
	Capabilities []types.String    `tfsdk:"capabilities"`
	ActiveSince  timetypes.RFC3339 `tfsdk:"active_since"`
	ID           types.String      `tfsdk:"id"`
}

func (d *DataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_account"
}

func (d *DataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = DataSourceSchema()
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data.linode_account")
	client := d.client

	var data DataSourceModel

	tflog.Trace(ctx, "client.GetAccount(...)")

	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Account: %s", err.Error(),
		)
		return
	}

	data.ParseAccount(account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
