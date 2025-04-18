package profile

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func NewDataSource() datasource.DataSource {
	return &DataSource{
		BaseDataSource: helper.NewBaseDataSource(
			helper.BaseDataSourceConfig{
				Name:   "linode_profile",
				Schema: &frameworkDatasourceSchema,
			},
		),
	}
}

type DataSource struct {
	helper.BaseDataSource
}

func (data *DataSourceModel) parseProfile(ctx context.Context, profile *linodego.Profile) diag.Diagnostics {
	data.Email = types.StringValue(profile.Email)
	data.Timezone = types.StringValue(profile.Timezone)
	data.EmailNotifications = types.BoolValue(profile.EmailNotifications)
	data.Username = types.StringValue(profile.Username)
	data.IPWhitelistEnabled = types.BoolValue(profile.IPWhitelistEnabled)
	data.LishAuthMethod = types.StringValue(string(profile.LishAuthMethod))

	authorized_keys, diags := types.ListValueFrom(ctx, types.StringType, profile.AuthorizedKeys)
	if diags.HasError() {
		return diags
	}

	data.AuthorizedKeys = authorized_keys

	data.TwoFactorAuth = types.BoolValue(profile.TwoFactorAuth)
	data.Restricted = types.BoolValue(profile.Restricted)

	referrals, diags := flattenReferral(ctx, profile.Referrals)
	if diags.HasError() {
		return diags
	}

	data.Referrals = *referrals

	id, err := json.Marshal(profile)
	if err != nil {
		diags.AddError("Error marshalling json: %s", err.Error())
		return diags
	}

	data.ID = types.StringValue(string(id))

	return nil
}

type DataSourceModel struct {
	Email              types.String `tfsdk:"email"`
	Timezone           types.String `tfsdk:"timezone"`
	EmailNotifications types.Bool   `tfsdk:"email_notifications"`
	Username           types.String `tfsdk:"username"`
	IPWhitelistEnabled types.Bool   `tfsdk:"ip_whitelist_enabled"`
	LishAuthMethod     types.String `tfsdk:"lish_auth_method"`
	AuthorizedKeys     types.List   `tfsdk:"authorized_keys"`
	TwoFactorAuth      types.Bool   `tfsdk:"two_factor_auth"`
	Restricted         types.Bool   `tfsdk:"restricted"`
	Referrals          types.List   `tfsdk:"referrals"`
	ID                 types.String `tfsdk:"id"`
}

func (d *DataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	tflog.Debug(ctx, "Read data."+d.Config.Name)

	client := d.Meta.Client

	var data DataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	profile, err := client.GetProfile(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get profile", err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.parseProfile(ctx, profile)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenReferral(ctx context.Context,
	referral linodego.ProfileReferrals,
) (*basetypes.ListValue, diag.Diagnostics) {
	result := make(map[string]attr.Value)

	result["total"] = types.Int64Value(int64(referral.Total))
	result["completed"] = types.Int64Value(int64(referral.Completed))
	result["pending"] = types.Int64Value(int64(referral.Pending))
	result["credit"] = types.Float64Value(float64(referral.Credit))
	result["code"] = types.StringValue(referral.Code)
	result["url"] = types.StringValue(referral.URL)

	obj, diags := types.ObjectValue(referralObjectType.AttrTypes, result)
	if diags.HasError() {
		return nil, diags
	}

	objList := []attr.Value{obj}

	resultList, diag := basetypes.NewListValue(
		referralObjectType,
		objList,
	)
	if diag.HasError() {
		return nil, diag
	}

	return &resultList, nil
}
