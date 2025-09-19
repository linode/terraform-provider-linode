package linodeinterface

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

type DefaultRouteAttrModel struct {
	IPv4 types.Bool `tfsdk:"ipv4"`
	IPv6 types.Bool `tfsdk:"ipv6"`
}

func (plan *DefaultRouteAttrModel) GetCreateOrUpdateOptions(state *DefaultRouteAttrModel) (opts linodego.InterfaceDefaultRoute, shouldUpdate bool) {
	if !plan.IPv4.IsUnknown() && (state == nil || !state.IPv4.Equal(plan.IPv4)) {
		opts.IPv4 = plan.IPv4.ValueBoolPointer()
		shouldUpdate = true
	}

	if !plan.IPv6.IsUnknown() && (state == nil || !state.IPv6.Equal(plan.IPv6)) {
		opts.IPv6 = plan.IPv6.ValueBoolPointer()
		shouldUpdate = true
	}

	return
}

func (data *DefaultRouteAttrModel) FlattenInterfaceDefaultRoute(
	defaultRoute linodego.InterfaceDefaultRoute, preserveKnown bool,
) {
	data.IPv4 = helper.KeepOrUpdateBoolPointer(data.IPv4, defaultRoute.IPv4, preserveKnown)
	data.IPv6 = helper.KeepOrUpdateBoolPointer(data.IPv6, defaultRoute.IPv6, preserveKnown)
}
