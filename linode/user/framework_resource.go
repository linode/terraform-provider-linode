package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "linode_user",
				IDType: types.StringType,
				Schema: &frameworkResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

// var resourceLinodeUserGrantFields = []string{
// 	"global_grants", "domain_grant", "firewall_grant", "image_grant",
// 	"linode_grant", "longview_grant", "nodebalancer_grant", "stackscript_grant", "volume_grant",
// }

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data UserModel
	client := r.Meta.Client

	fmt.Printf("req.Config.Raw: %v\n", req.Config.Raw)
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.UserCreateOptions{
		Email:      data.Email.ValueString(),
		Username:   data.Username.ValueString(),
		Restricted: data.Restricted.ValueBool(),
	}

	user, err := client.CreateUser(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create user.",
			err.Error(),
		)
		return
	}
	// TODO: implement userHasGrantsConfigured
	if userHasGrantsConfigured(data) {
		if diag := data.updateUserGrants(ctx, client); diag != nil {
			resp.Diagnostics.Append(diag)
			return
		}
	}

	data.ID = types.StringValue(user.Username)
	// TODO: parse computed
	resp.Diagnostics.Append(data.parseComputedAttrs(ctx, user)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data UserModel
	client := r.Meta.Client

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := client.GetUser(ctx, data.Username.ValueString())
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"User does not exist.",
				fmt.Sprintf("Removing Linode User %v from state because it no longer exists", data.ID.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to get user %v.", data.Username.ValueString()),
			err.Error(),
		)
		return
	}
	resp.Diagnostics.Append(data.parseComputedAttrs(ctx, user)...)
	// TODO: parse non computed
	if resp.Diagnostics.HasError() {
		return
	}

	if user.Restricted {
		_, err := client.GetUserGrants(ctx, data.Username.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to get User Grants (%s): ", data.Username.ValueString()), err.Error(),
			)
			return
		}
		// TODO: PARSE resp.Diagnostics.Append(data.ParseUserGrants(ctx, grants)...)
	} else {
		// TODO: parse no grants
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {

}

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
}

func (d *UserModel) updateUserGrants(
	ctx context.Context,
	client *linodego.Client,
) diag.Diagnostic {
	if !d.Restricted.ValueBool() {
		return diag.NewErrorDiagnostic(
			"Failed to set user grants.",
			"User must be restricted in order to update grants")
	}

	updateOpts := linodego.UserGrantsUpdateOptions{}

	if len(d.GlobalGrants) > 0 {
		updateOpts.Global = expandGlobalGrant(d.GlobalGrants[0])
	}

	updateOpts.Domain = expandUserGrantsEntities(d.DomainGrant)
	updateOpts.Firewall = expandUserGrantsEntities(d.FirewallGrant)
	updateOpts.Image = expandUserGrantsEntities(d.ImageGrant)
	updateOpts.Linode = expandUserGrantsEntities(d.LinodeGrant)
	updateOpts.Longview = expandUserGrantsEntities(d.LongviewGrant)
	updateOpts.NodeBalancer = expandUserGrantsEntities(d.NodebalancerGrant)
	updateOpts.StackScript = expandUserGrantsEntities(d.StackscriptGrant)
	updateOpts.Volume = expandUserGrantsEntities(d.VolumeGrant)

	if _, err := client.UpdateUserGrants(ctx, d.Username.ValueString(), updateOpts); err != nil {
		return diag.NewErrorDiagnostic(
			fmt.Sprintf("Failed to set user grants %v.", d.Username.ValueString()),
			err.Error(),
		)
	}

	return nil
}

func userHasGrantsConfigured(data UserModel) bool {
	// data.DatabaseGrant
	// if !data.GlobalGrants.IsNull() && !data.GlobalGrants.IsUnknown() {
	// 	return true
	// }
	// if !data.DatabaseGrant.IsNull() && !data.DatabaseGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.DomainGrant.IsNull() && !data.DomainGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.FirewallGrant.IsNull() && !data.FirewallGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.ImageGrant.IsNull() && !data.ImageGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.LinodeGrant.IsNull() && !data.LinodeGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.LongviewGrant.IsNull() && !data.LongviewGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.NodebalancerGrant.IsNull() && !data.NodebalancerGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.StackscriptGrant.IsNull() && !data.StackscriptGrant.IsUnknown() {
	// 	return true
	// }
	// if !data.VolumeGrant.IsNull() && !data.VolumeGrant.IsUnknown() {
	// 	return true
	// }
	// data.GlobalGrants.IsNull()

	// values := reflect.ValueOf(data)
	// dataType := values.Type()

	// for i := 0; i < values.NumField(); i++ {
	// 	if slices.Contains(resourceLinodeUserGrantFields, dataType.Field(i).Name) &&
	// 		!values.Field(i).Interface() {
	// 		return true
	// 	}
	// }

	// for _, key := range resourceLinodeUserGrantFields {
	// 	if d.key
	// 	// if _, ok := d.GetOk(key); ok {
	// 	// 	return true
	// 	// }
	// }

	return false
}
