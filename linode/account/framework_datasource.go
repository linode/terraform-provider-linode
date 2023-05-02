package account

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *linodego.Client
}

func (data *DataSourceModel) parseAccount(account *linodego.Account) {
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

	meta := helper.GetMetaFromProviderDataDatasource(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client = meta.Client
}

type DataSourceModel struct {
	Email     types.String  `tfsdk:"email"`
	FirstName types.String  `tfsdk:"first_name"`
	LastName  types.String  `tfsdk:"last_name"`
	Company   types.String  `tfsdk:"company"`
	Address1  types.String  `tfsdk:"address_1"`
	Address2  types.String  `tfsdk:"address_2"`
	Phone     types.String  `tfsdk:"phone"`
	City      types.String  `tfsdk:"city"`
	State     types.String  `tfsdk:"state"`
	Country   types.String  `tfsdk:"country"`
	Zip       types.String  `tfsdk:"zip"`
	Balance   types.Float64 `tfsdk:"balance"`
	ID        types.String  `tfsdk:"id"`
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
	resp.Schema = frameworkDatasourceSchema
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.client

	var data DataSourceModel

	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Account",
			fmt.Sprintf(
				"Error finding Linode Account: %s",
				err.Error(),
			),
		)
		return
	}

	data.parseAccount(account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
