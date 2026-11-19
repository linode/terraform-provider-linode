package accounttfaoptionalusers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
)

const singletonID = "account"

type ResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Usernames types.Set    `tfsdk:"usernames"`
}

type DataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Usernames types.Set    `tfsdk:"usernames"`
}

func readUsernames(ctx context.Context, client *linodego.Client, diags *diag.Diagnostics) types.Set {
	tflog.Trace(ctx, "client.ListTFAOptionalUsers(...)")
	users, err := client.ListTFAOptionalUsers(ctx, nil)
	if err != nil {
		diags.AddError("Failed to list TFA optional users", err.Error())
		return types.SetNull(types.StringType)
	}

	usernames := make([]string, len(users))
	for i, user := range users {
		usernames[i] = user.Username
	}

	result, d := types.SetValueFrom(ctx, types.StringType, usernames)
	diags.Append(d...)
	return result
}
