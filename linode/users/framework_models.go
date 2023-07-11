package users

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper/frameworkfilter"
	"github.com/linode/terraform-provider-linode/linode/user"
)

// UserFilterModel describes the Terraform resource data model to match the
// resource schema.
type UserFilterModel struct {
	ID      types.String                     `tfsdk:"id"`
	Filters frameworkfilter.FiltersModelType `tfsdk:"filter"`
	Order   types.String                     `tfsdk:"order"`
	OrderBy types.String                     `tfsdk:"order_by"`
	Users   []user.DataSourceModel           `tfsdk:"users"`
}

func (data *UserFilterModel) parseUsers(
	ctx context.Context,
	client *linodego.Client,
	users []linodego.User,
) diag.Diagnostics {
	result := make([]user.DataSourceModel, len(users))
	for i := range users {
		var userModel user.DataSourceModel
		diags := userModel.ParseUser(ctx, &users[i])
		if diags.HasError() {
			return diags
		}

		if users[i].Restricted {
			grants, err := client.GetUserGrants(ctx, userModel.Username.ValueString())
			if err != nil {
				diags.AddError(
					fmt.Sprintf("Failed to get User Grants (%s): ", userModel.Username.ValueString()),
					err.Error(),
				)
				return diags
			}
			diags := userModel.ParseUserGrants(ctx, grants)
			if diags != nil {
				return diags
			}
		} else {
			userModel.ParseNonUserGrants()
		}
		result[i] = userModel
	}

	data.Users = result
	return nil
}
