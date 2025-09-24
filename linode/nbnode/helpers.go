package nbnode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/linode/linodego"
)

func safeFetchVPCConfig(
	ctx context.Context,
	client *linodego.Client,
	nodeBalancerID int,
	vpcConfigID int,
) (*linodego.NodeBalancerVPCConfig, diag.Diagnostics) {
	var d diag.Diagnostics

	if vpcConfigID == 0 {
		return nil, d
	}

	result, err := client.GetNodeBalancerVPCConfig(
		ctx,
		nodeBalancerID,
		vpcConfigID,
	)
	if err != nil {
		if linodego.IsNotFound(err) {
			// The user might not have access to NB/VPC
			return nil, nil
		}

		d.AddError(
			"Failed to get NodeBalancer VPC configuration",
			err.Error(),
		)
		return nil, d
	}

	return result, d
}
