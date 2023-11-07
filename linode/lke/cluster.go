package lke

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type NodePoolSpec struct {
	Type              string
	Count             int
	AutoScalerEnabled bool
	AutoScalerMin     int
	AutoScalerMax     int
}

type NodePoolUpdates struct {
	ToDelete []int
	ToCreate []linodego.LKENodePoolCreateOptions
	ToUpdate map[int]linodego.LKENodePoolUpdateOptions
}

type nodePoolAssignRequest struct {
	Spec, State NodePoolSpec
	PoolID      int
	SpecIndex   int
}

func (r nodePoolAssignRequest) Diff() int {
	return int(math.Abs(float64(r.State.Count - r.Spec.Count)))
}

func getLKENodePoolProvisionedSpecs(pools []linodego.LKENodePool) map[NodePoolSpec]map[int]struct{} {
	provisioned := make(map[NodePoolSpec]map[int]struct{})
	for _, pool := range pools {
		spec := NodePoolSpec{
			Type:              pool.Type,
			Count:             pool.Count,
			AutoScalerEnabled: pool.Autoscaler.Enabled,
			AutoScalerMin:     pool.Autoscaler.Min,
			AutoScalerMax:     pool.Autoscaler.Max,
		}
		if _, ok := provisioned[spec]; !ok {
			provisioned[spec] = make(map[int]struct{})
		}
		provisioned[spec][pool.ID] = struct{}{}
	}
	return provisioned
}

func ReconcileLKENodePoolSpecs(
	poolSpecs []NodePoolSpec, pools []linodego.LKENodePool,
) (updates NodePoolUpdates) {
	provisionedPools := getLKENodePoolProvisionedSpecs(pools)
	poolSpecsToAssign := make(map[int]struct{})
	assignedPools := make(map[int]struct{})
	updates.ToUpdate = make(map[int]linodego.LKENodePoolUpdateOptions)

	// find exact pool matches and filter out
	for i, spec := range poolSpecs {
		poolSpecsToAssign[i] = struct{}{}
		if ids, ok := provisionedPools[spec]; ok {
			for id := range ids {
				assignedPools[i] = struct{}{}
				delete(ids, id)
				break
			}

			if len(provisionedPools[spec]) == 0 {
				delete(provisionedPools, spec)
			}

			delete(poolSpecsToAssign, i)
		}
	}

	// calculate diffs for assigning remaining provisioned pools to remaining pool specs
	poolAssignRequests := []nodePoolAssignRequest{}
	for i := range poolSpecsToAssign {
		poolSpec := poolSpecs[i]
		for pool := range provisionedPools {
			if pool.Type != poolSpec.Type {
				continue
			}

			for id := range provisionedPools[pool] {
				poolAssignRequests = append(poolAssignRequests, nodePoolAssignRequest{
					Spec:      poolSpec,
					State:     pool,
					PoolID:    id,
					SpecIndex: i,
				})
			}
		}
	}

	// order poolAssignRequests by smallest diffs for smallest updates needed
	sort.Slice(poolAssignRequests, func(x, y int) bool {
		return poolAssignRequests[x].Diff() < poolAssignRequests[y].Diff()
	})

	for _, request := range poolAssignRequests {
		if _, ok := poolSpecsToAssign[request.SpecIndex]; !ok {
			// pool spec was already assigned to a provisioned pool
			continue
		}
		if _, ok := assignedPools[request.PoolID]; ok {
			// pool was already assigned to a pool spec
			continue
		}

		var newAutoscaler *linodego.LKENodePoolAutoscaler

		if request.Spec.AutoScalerEnabled {
			newAutoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: request.Spec.AutoScalerEnabled,
				Min:     request.Spec.AutoScalerMin,
				Max:     request.Spec.AutoScalerMax,
			}
		}

		// Only disable if already enabled
		if !request.Spec.AutoScalerEnabled && request.State.AutoScalerEnabled {
			newAutoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: request.Spec.AutoScalerEnabled,
				Min:     request.Spec.Count,
				Max:     request.Spec.Count,
			}
		}

		updates.ToUpdate[request.PoolID] = linodego.LKENodePoolUpdateOptions{
			Count:      request.Spec.Count,
			Autoscaler: newAutoscaler,
		}

		assignedPools[request.PoolID] = struct{}{}
		delete(poolSpecsToAssign, request.SpecIndex)
		delete(provisionedPools[request.State], request.PoolID)
		if len(provisionedPools[request.State]) == 0 {
			delete(provisionedPools, request.State)
		}
	}

	for i := range poolSpecsToAssign {
		poolSpec := poolSpecs[i]

		var newAutoscaler *linodego.LKENodePoolAutoscaler

		if poolSpec.AutoScalerEnabled {
			newAutoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: poolSpec.AutoScalerEnabled,
				Min:     poolSpec.AutoScalerMin,
				Max:     poolSpec.AutoScalerMax,
			}
		}

		updates.ToCreate = append(updates.ToCreate, linodego.LKENodePoolCreateOptions{
			Count:      poolSpec.Count,
			Type:       poolSpec.Type,
			Autoscaler: newAutoscaler,
		})
	}

	for spec := range provisionedPools {
		for id := range provisionedPools[spec] {
			updates.ToDelete = append(updates.ToDelete, id)
		}
	}

	return
}

func waitForNodePoolReady(
	ctx context.Context, client linodego.Client, pollMs, clusterID, poolID int,
) error {
	eventTicker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for LKE Cluster (%d) Pool (%d) to be ready", clusterID, poolID)

		case <-eventTicker.C:
			pool, err := client.GetLKENodePool(ctx, clusterID, poolID)
			if err != nil {
				return fmt.Errorf("failed to get LKE Cluster (%d) Pool (%d): %w", clusterID, poolID, err)
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

			return nil
		}
	}
}

func waitForNodesDeleted(
	ctx context.Context,
	client linodego.Client,
	intervalMS int,
	nodes []linodego.LKENodePoolLinode,
) error {
	ticker := time.NewTicker(time.Duration(intervalMS) * time.Millisecond)
	defer ticker.Stop()

	// Let's track which nodes still haven't been deleted
	// using a pseudo-set
	remainingNodes := make(map[int]bool, len(nodes))
	for _, node := range nodes {
		remainingNodes[node.InstanceID] = true
	}

	// Filter down to only instance deletion events
	f := linodego.Filter{
		OrderBy: "created",
		Order:   linodego.Descending,
	}
	f.AddField(linodego.Eq, "entity.type", linodego.EntityLinode)
	f.AddField(linodego.Eq, "action", linodego.ActionLinodeDelete)

	filterBytes, err := f.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal filter: %w", err)
	}

	listOpts := linodego.ListOptions{
		Filter:      string(filterBytes),
		PageOptions: &linodego.PageOptions{Page: 1},
	}

	for {
		select {
		case <-ticker.C:
			events, err := client.ListEvents(ctx, &listOpts)
			if err != nil {
				return fmt.Errorf("failed to list events: %w", err)
			}

			for _, event := range events {
				var instID int

				// Sometimes go will parse entity.id as float,
				// we should convert accordingly
				switch event.Entity.ID.(type) {
				case int:
					instID = event.Entity.ID.(int)
				case float64:
					instID = int(event.Entity.ID.(float64))
				case float32:
					instID = int(event.Entity.ID.(float32))
				default:
					// This shouldn't happen, but let's handle it gracefully just in case
					tflog.Trace(ctx, "Invalid entity.id type detected", map[string]any{
						"value": event.Entity.ID,
						"type":  fmt.Sprintf("%T", event.Entity.ID),
					})
					continue
				}

				if _, ok := remainingNodes[instID]; ok {
					delete(remainingNodes, instID)
					tflog.Trace(ctx, "Node detected as deleted", map[string]any{
						"instance_id":     instID,
						"nodes_remaining": len(remainingNodes),
					})
				}
			}

			// All nodes have been deleted
			if len(remainingNodes) < 1 {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("failed to wait for node deletion: %w", ctx.Err())
		}
	}
}

func recycleLKECluster(ctx context.Context, meta *helper.ProviderMeta, id int, pools []linodego.LKENodePool) error {
	client := meta.Client

	ctx = helper.SetLogFieldBulk(ctx, map[string]any{
		"cluster_id": id,
		"pools":      pools,
	})

	tflog.Info(ctx, "Recycling LKE cluster")

	if err := client.RecycleLKEClusterNodes(ctx, id); err != nil {
		return fmt.Errorf("failed to recycle LKE Cluster (%d): %s", id, err)
	}

	// Aggregate all nodes to be recycled
	oldNodes := make([]linodego.LKENodePoolLinode, 0)
	for _, pool := range pools {
		oldNodes = append(oldNodes, pool.Linodes...)
	}

	tflog.Debug(ctx, "Waiting for all nodes to be deleted", map[string]any{
		"nodes": oldNodes,
	})

	// Wait for the old nodes to be deleted
	if err := waitForNodesDeleted(ctx, client, meta.Config.EventPollMilliseconds, oldNodes); err != nil {
		return fmt.Errorf("failed to wait for old nodes to be recycled: %w", err)
	}

	tflog.Debug(ctx, "All old nodes detected as deleted, waiting for all node pools to enter ready status")

	// Wait for all node pools to be ready
	for _, pool := range pools {
		if err := waitForNodePoolReady(ctx, client, meta.Config.EventPollMilliseconds, id, pool.ID); err != nil {
			return fmt.Errorf("failed to wait for pool %d ready: %w", pool.ID, err)
		}
	}

	tflog.Debug(ctx, "All node pools have entered ready status; recycle operation completed")

	return nil
}

// This cannot currently be handled efficiently by a DiffSuppressFunc
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchPoolsWithSchema(pools []linodego.LKENodePool, declaredPools []interface{}) []linodego.LKEClusterPool {
	result := make([]linodego.LKENodePool, len(declaredPools))

	poolMap := make(map[int]linodego.LKENodePool, len(declaredPools))
	for _, pool := range pools {
		poolMap[pool.ID] = pool
	}

	for i, declaredPool := range declaredPools {
		declaredPool := declaredPool.(map[string]interface{})

		for key, pool := range poolMap {
			if pool.Count != declaredPool["count"] || pool.Type != declaredPool["type"] {
				continue
			}

			result[i] = pool
			delete(poolMap, key)
			break
		}
	}

	for _, pool := range poolMap {
		//nolint:makezero
		result = append(result, pool)
	}

	return result
}

func expandLinodeLKEClusterAutoscalerFromPool(pool map[string]interface{}) *linodego.LKENodePoolAutoscaler {
	scalersSpec, ok := pool["autoscaler"].([]interface{})

	// Return nil if the autoscaler isn't defined
	if !ok || len(scalersSpec) < 1 {
		return nil
	}

	scalerSpec := scalersSpec[0].(map[string]interface{})
	return &linodego.LKENodePoolAutoscaler{
		Enabled: true,
		Min:     scalerSpec["min"].(int),
		Max:     scalerSpec["max"].(int),
	}
}

func expandLinodeLKENodePoolSpecs(pool []interface{}) (poolSpecs []NodePoolSpec) {
	for _, spec := range pool {
		specMap := spec.(map[string]interface{})
		autoscaler := expandLinodeLKEClusterAutoscalerFromPool(specMap)
		if autoscaler == nil {
			autoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: false,
				Min:     specMap["count"].(int),
				Max:     specMap["count"].(int),
			}
		}

		poolSpecs = append(poolSpecs, NodePoolSpec{
			Type:              specMap["type"].(string),
			Count:             specMap["count"].(int),
			AutoScalerEnabled: autoscaler.Enabled,
			AutoScalerMin:     autoscaler.Min,
			AutoScalerMax:     autoscaler.Max,
		})
	}
	return
}

func flattenLKENodePools(pools []linodego.LKEClusterPool) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(pools))
	for i, pool := range pools {

		nodes := make([]map[string]interface{}, len(pool.Linodes))
		for i, node := range pool.Linodes {
			nodes[i] = map[string]interface{}{
				"id":          node.ID,
				"instance_id": node.InstanceID,
				"status":      node.Status,
			}
		}

		var autoscaler []map[string]interface{}

		if pool.Autoscaler.Enabled {
			autoscaler = []map[string]interface{}{
				{
					"min": pool.Autoscaler.Min,
					"max": pool.Autoscaler.Max,
				},
			}
		}

		flattened[i] = map[string]interface{}{
			"id":         pool.ID,
			"count":      pool.Count,
			"type":       pool.Type,
			"nodes":      nodes,
			"autoscaler": autoscaler,
		}
	}
	return flattened
}

func flattenLKEClusterControlPlane(controlPlane linodego.LKEClusterControlPlane) map[string]interface{} {
	flattened := make(map[string]interface{})

	flattened["high_availability"] = controlPlane.HighAvailability

	return flattened
}

func expandLKEClusterControlPlane(controlPlane map[string]interface{}) linodego.LKEClusterControlPlane {
	var result linodego.LKEClusterControlPlane

	if value, ok := controlPlane["high_availability"]; ok {
		result.HighAvailability = value.(bool)
	}

	return result
}
