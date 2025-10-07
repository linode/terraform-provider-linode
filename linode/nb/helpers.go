package nb

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

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
