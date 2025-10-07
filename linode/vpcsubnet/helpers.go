package vpcsubnet

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// shouldRetryOn400sForDatabasePropagation returns whether all attached databases
// are in the process of propagating a detachment/deletion.
//
// This is necessary because database deletion propagation can often take a while,
// leading to errors during destruction.
func shouldRetryOn400sForDatabasePropagation(
	ctx context.Context,
	client *linodego.Client,
	vpcID int,
	vpcSubnetID int,
) (bool, error) {
	tflog.Debug(
		ctx,
		"client.GetVPCSubnet(...)",
		map[string]any{
			"vpc_id":    vpcID,
			"subnet_id": vpcSubnetID,
		},
	)
	vpcSubnet, err := client.GetVPCSubnet(ctx, vpcID, vpcSubnetID)
	if err != nil {
		return false, fmt.Errorf("failed to get VPC subnet %d for VPC %d: %w", vpcSubnetID, vpcID, err)
	}

	if len(vpcSubnet.Databases) < 1 {
		// Nothing to do here
		return false, nil
	}

	tflog.Debug(
		ctx,
		"client.ListDatabases(...)",
	)
	databases, err := client.ListDatabases(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to list databases: %w", err)
	}

	// Aggregate which databases still exist on the current account
	validDatabaseIDs := make(map[int]bool, len(databases))
	for _, db := range databases {
		if db.PrivateNetwork == nil || db.PrivateNetwork.VPCID != vpcID {
			// If this is a detachment that hasn't yet propagated,
			// don't set the status
			continue
		}

		validDatabaseIDs[db.ID] = true
	}

	// Aggregate which databases are still registered with the subnet
	// but do not exist on the current account
	pendingDetachPropagation := make(map[int]bool)
	for _, db := range vpcSubnet.Databases {
		if _, ok := validDatabaseIDs[db.ID]; !ok {
			pendingDetachPropagation[db.ID] = true
		}
	}

	return len(pendingDetachPropagation) >= len(vpcSubnet.Databases), nil
}
