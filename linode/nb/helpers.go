package nb

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
)

func validateIPv4RangeAutoAssignConflict(
	ctx context.Context,
	vpcs types.List,
	attrName string,
	diags *diag.Diagnostics,
) {
	if vpcs.IsNull() || vpcs.IsUnknown() {
		return
	}

	var models []ResourceBackendVPCModel
	diags.Append(vpcs.ElementsAs(ctx, &models, false)...)
	if diags.HasError() {
		return
	}

	for i, model := range models {
		if model.IPv4RangeAutoAssign.IsNull() ||
			model.IPv4RangeAutoAssign.IsUnknown() ||
			!model.IPv4RangeAutoAssign.ValueBool() {
			continue
		}

		if model.IPv4Range.IsNull() || model.IPv4Range.IsUnknown() {
			continue
		}

		diags.AddAttributeError(
			path.Root(attrName).AtListIndex(i).AtName("ipv4_range"),
			"Invalid Attribute Combination",
			fmt.Sprintf(
				"`ipv4_range` cannot be set when `ipv4_range_auto_assign` is true in %s[%d].",
				attrName,
				i,
			),
		)
	}
}

func safeListVPCConfigs(
	ctx context.Context,
	client *linodego.Client,
	nodeBalancerID int,
	listOptions *linodego.ListOptions,
	diagnostics diag.Diagnostics,
) []linodego.NodeBalancerVPCConfig {
	tflog.Trace(ctx, "client.ListNodeBalancerVPCConfigs(...)")

	result, err := helper.NotFoundDefault(
		func() ([]linodego.NodeBalancerVPCConfig, error) {
			return client.ListNodeBalancerVPCConfigs(
				ctx,
				nodeBalancerID,
				listOptions,
			)
		},
		nil,
	)
	if err != nil {
		diagnostics.AddError(
			"Failed to list NodeBalancer VPC configurations",
			err.Error(),
		)
	}

	return result
}

func safeListFirewalls(
	ctx context.Context,
	client *linodego.Client,
	nodeBalancerID int,
	listOptions *linodego.ListOptions,
	diagnostics diag.Diagnostics,
) []linodego.Firewall {
	tflog.Trace(ctx, "client.ListNodeBalancerFirewalls(...)")

	result, err := helper.NotFoundDefault(
		func() ([]linodego.Firewall, error) {
			return client.ListNodeBalancerFirewalls(
				ctx,
				nodeBalancerID,
				listOptions,
			)
		},
		nil,
	)
	if err != nil {
		diagnostics.AddError(
			"Failed to list NodeBalancer Firewalls",
			err.Error(),
		)
	}

	return result
}
