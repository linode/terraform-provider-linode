package vpcsubnet

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

// waitForDatabaseDetachmentsPropagated is waits until the deletions of attached
// databases have propagated to this subnet's assignments.
//
// This is necessary because database deletion propagation can often take a while,
// leading to errors during destruction.
//
// TODO: Make this type-agnostic.
func waitForDatabaseDetachmentsPropagated(
	ctx context.Context,
	client *linodego.Client,
	vpcID int,
	vpcSubnetID int,
) error {
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
		return fmt.Errorf("failed to get VPC subnet %d for VPC %d: %w", vpcSubnetID, vpcID, err)
	}

	if len(vpcSubnet.Databases) < 1 {
		// Nothing to do here
		return nil
	}

	tflog.Debug(
		ctx,
		"client.ListDatabases(...)",
	)
	databases, err := client.ListDatabases(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list databases: %w", err)
	}

	// Aggregate which databases still exist on the current account
	validDatabaseIDs := make(map[int]bool, len(databases))
	for _, db := range databases {
		if db.PrivateNetwork == nil || db.PrivateNetwork.VPCID != db.ID {
			// If this is a detachment that hasn't yet propagated,
			// don't set the status
			continue
		}

		validDatabaseIDs[db.ID] = true
	}

	// Aggregate which databases are still registered with the subnet
	// but do not exist on the current account
	pendingDeletePropagation := make(map[int]bool)
	for _, db := range vpcSubnet.Databases {
		if _, ok := validDatabaseIDs[db.ID]; !ok {
			pendingDeletePropagation[db.ID] = true
		}
	}

	if len(pendingDeletePropagation) != len(vpcSubnet.Databases) {
		tflog.Debug(
			ctx,
			"Databases still exist",
			map[string]any{
				"pending":  slices.Collect(maps.Keys(pendingDeletePropagation)),
				"assigned": vpcSubnet.Databases,
			})
		// There is still an active database - defer to the API for
		// error handling
		return nil
	}

	tflog.Info(
		ctx,
		"Waiting for database deletions to propagate",
		map[string]any{
			"database_ids": slices.Collect(maps.Keys(pendingDeletePropagation)),
		},
	)

	ticker := time.NewTicker(client.GetPollDelay())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tflog.Debug(
				ctx,
				"client.GetVPCSubnet(...)",
				map[string]any{
					"vpc_id":    vpcID,
					"subnet_id": vpcSubnetID,
				},
			)
			vpcSubnet, err = client.GetVPCSubnet(ctx, vpcID, vpcSubnetID)
			if err != nil {
				return fmt.Errorf("failed to get VPC subnet %d for VPC %d: %w", vpcSubnetID, vpcSubnetID, err)
			}

			tflog.Debug(
				ctx,
				"Refreshed subnet databases",
				map[string]any{
					"vpc_id":    vpcID,
					"subnet_id": vpcSubnetID,
					"dbs":       vpcSubnet.Databases,
				},
			)

			if len(vpcSubnet.Databases) < 1 {
				// We're done!
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for database deletions to propagate: %w", ctx.Err())
		}
	}
}
