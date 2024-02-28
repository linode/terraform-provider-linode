package nodepool

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
)

func WaitForNodePoolReady(
	ctx context.Context, client linodego.Client, pollMs, clusterID, poolID int,
) (*linodego.LKENodePool, error) {
	ctx = tflog.SetField(ctx, "node_pool_id", poolID)
	eventTicker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out waiting for LKE Cluster (%d) Pool (%d) to be ready", clusterID, poolID)

		case <-eventTicker.C:
			tflog.Trace(ctx, "client.GetLKENodePool(...)")
			pool, err := client.GetLKENodePool(ctx, clusterID, poolID)
			if err != nil {
				return nil, fmt.Errorf("failed to get LKE Cluster (%d) Pool (%d): %w", clusterID, poolID, err)
			}

			allNodesReady := true

			for _, instance := range pool.Linodes {
				if instance.Status == linodego.LKELinodeNotReady {
					allNodesReady = false
					tflog.Trace(ctx, "Node detected as unready", map[string]any{
						"node_id":     instance.ID,
						"instance_id": instance.InstanceID,
					})
					break
				}
			}

			if !allNodesReady {
				continue
			}

			// We're finished!
			tflog.Trace(ctx, "All nodes ready!")

			return pool, nil
		}
	}
}
