package account

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

func NewDatasource() datasource.DataSource {
	return &Datasource{}
}

type Datasource struct {
	client *linodego.Client
}

func (data *DatasourceModel) parseAccount(account *linodego.Account) {
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
}

// ResourceModel describes the Terraform resource data model to match the
// resource schema.
type DatasourceModel struct {
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
}

func (d *Datasource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "linode_account"
}

func (d *Datasource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = frameworkDatasourceSchema
}

func (d *Datasource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	client := d.client

	var data DatasourceModel

	account, err := client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get Account",
			fmt.Sprintf(
				"Error finding Linode Account: %s",
				err.Error(),
			),
		)
	}

	data.parseAccount(account)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
