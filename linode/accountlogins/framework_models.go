package accountlogins

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/frameworkfilter"
)

type AccountLoginModel struct {
	Datetime   types.String `tfsdk:"datetime"`
	ID         types.Int64  `tfsdk:"id"`
	IP         types.String `tfsdk:"ip"`
	Restricted types.Bool   `tfsdk:"restricted"`
	Username   types.String `tfsdk:"username"`
	Status     types.String `tfsdk:"status"`
}

type AccountLoginFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Logins  []AccountLoginModel              `tfsdk:"logins"`
}

func (model *AccountLoginFilterModel) parseLogins(logins []linodego.Login) {
	parseAccountLogin := func(login linodego.Login) AccountLoginModel {
		var m AccountLoginModel

		m.Datetime = types.StringValue(login.Datetime.Format(time.RFC3339))
		m.ID = types.Int64Value(int64(login.ID))
		m.IP = types.StringValue(login.IP)
		m.Restricted = types.BoolValue(login.Restricted)
		m.Username = types.StringValue(login.Username)
		m.Status = types.StringValue(login.Status)

		return m
	}

	result := make([]AccountLoginModel, len(logins))

	for i, login := range logins {
		result[i] = parseAccountLogin(login)
	}

	model.Logins = result
}
