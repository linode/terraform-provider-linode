package linodeinterface

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DefaultRouteAttrModel struct {
	IPv4 types.Bool `tfsdk:"ipv4"`
	IPv6 types.Bool `tfsdk:"ipv6"`
}

func (plan *DefaultRouteAttrModel) GetCreateOrUpdateOptions(
	ctx context.Context,
	state *DefaultRouteAttrModel,
) (opts linodego.InterfaceDefaultRoute, shouldUpdate bool) {
	tflog.Trace(ctx, "Enter DefaultRouteAttrModel.GetCreateOrUpdateOptions")
	if !plan.IPv4.IsUnknown() && (state == nil || !state.IPv4.Equal(plan.IPv4)) {
		opts.IPv4 = plan.IPv4.ValueBoolPointer()
		shouldUpdate = true
	}

	if !plan.IPv6.IsUnknown() && (state == nil || !state.IPv6.Equal(plan.IPv6)) {
		opts.IPv6 = plan.IPv6.ValueBoolPointer()
		shouldUpdate = true
	}

	return opts, shouldUpdate
}

func (data *DefaultRouteAttrModel) FlattenInterfaceDefaultRoute(
	ctx context.Context, defaultRoute linodego.InterfaceDefaultRoute, preserveKnown bool,
) {
	tflog.Trace(ctx, "Enter DefaultRouteAttrModel.FlattenInterfaceDefaultRoute")

	data.IPv4 = helper.KeepOrUpdateBoolPointer(data.IPv4, defaultRoute.IPv4, preserveKnown)
	data.IPv6 = helper.KeepOrUpdateBoolPointer(data.IPv6, defaultRoute.IPv6, preserveKnown)
}
