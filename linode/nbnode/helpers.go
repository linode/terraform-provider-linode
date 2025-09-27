package nbnode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

func safeGetVPCConfig(
	ctx context.Context,
	client *linodego.Client,
	nodeBalancerID int,
	vpcConfigID int,
	diagnostics diag.Diagnostics,
) *linodego.NodeBalancerVPCConfig {
	result, err := helper.NotFoundDefault(
		func() (*linodego.NodeBalancerVPCConfig, error) {
			return client.GetNodeBalancerVPCConfig(
				ctx,
				nodeBalancerID,
				vpcConfigID,
			)
		},
		nil,
	)
	if err != nil {
		diagnostics.AddError(
			"Failed to get NodeBalancer VPC configuration",
			err.Error(),
		)
	}

	return result
}
