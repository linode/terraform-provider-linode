package lke

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

type NodePoolSpec struct {
	ID                int
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

func ReconcileLKENodePoolSpecs(
	ctx context.Context, oldSpecs []NodePoolSpec, newSpecs []NodePoolSpec,
) (updates NodePoolUpdates) {
	updates.ToCreate = make([]linodego.LKENodePoolCreateOptions, 0)
	updates.ToDelete = make([]int, 0)
	updates.ToUpdate = make(map[int]linodego.LKENodePoolUpdateOptions)

	// If there are fewer node pools than expected
	// we can assume the rest have been deleted
	if len(newSpecs) < len(oldSpecs) {
		for _, v := range oldSpecs[len(newSpecs):] {
			updates.ToDelete = append(updates.ToDelete, v.ID)
		}
	}

	// If there are more node pools then there were previously
	// we can assume new ones have been created
	if len(newSpecs) > len(oldSpecs) {
		for _, v := range newSpecs[len(oldSpecs):] {
			createOpts := linodego.LKENodePoolCreateOptions{
				Count: v.Count,
				Type:  v.Type,
			}

			if v.AutoScalerEnabled {
				createOpts.Autoscaler = &linodego.LKENodePoolAutoscaler{
					Enabled: true,
					Min:     v.AutoScalerMin,
					Max:     v.AutoScalerMax,
				}
			}

			updates.ToCreate = append(updates.ToCreate, createOpts)
		}
	}

	maxUpdateIndex := len(oldSpecs)
	if maxUpdateIndex > len(newSpecs) {
		maxUpdateIndex = len(newSpecs)
	}

	for i, newSpec := range newSpecs[:maxUpdateIndex] {
		oldSpec := oldSpecs[i]

		if reflect.DeepEqual(newSpec, oldSpec) {
			continue
		}

		// Types cannot be updated on node pools
		// so we should delete the old one and create a new one
		if newSpec.Type != oldSpec.Type {
			createOpts := linodego.LKENodePoolCreateOptions{
				Count: newSpec.Count,
				Type:  newSpec.Type,
			}

			if newSpec.AutoScalerEnabled {
				createOpts.Autoscaler = &linodego.LKENodePoolAutoscaler{
					Enabled: true,
					Min:     newSpec.AutoScalerMin,
					Max:     newSpec.AutoScalerMax,
				}
			}

			updates.ToDelete = append(updates.ToDelete, oldSpec.ID)
			updates.ToCreate = append(updates.ToCreate, createOpts)
			continue
		}

		updateOpts := linodego.LKENodePoolUpdateOptions{
			Count: newSpec.Count,
		}

		updateOpts.Autoscaler = &linodego.LKENodePoolAutoscaler{
			Enabled: newSpec.AutoScalerEnabled,
			Min:     newSpec.AutoScalerMin,
			Max:     newSpec.AutoScalerMax,
		}

		updates.ToUpdate[oldSpec.ID] = updateOpts
	}

	return
}

func waitForNodePoolReady(
	ctx context.Context, client linodego.Client, pollMs, clusterID, poolID int,
) error {
	ctx = tflog.SetField(ctx, "node_pool_id", poolID)
	eventTicker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for LKE Cluster (%d) Pool (%d) to be ready", clusterID, poolID)

		case <-eventTicker.C:
			tflog.Trace(ctx, "client.GetLKENodePool(...)")
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
			tflog.Trace(ctx, "client.ListEvents(...)", map[string]any{
				"options": listOpts,
			})

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
	tflog.Trace(ctx, "client.RecycleLKEClusterNodes(...)")

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
func matchPoolsWithSchema(ctx context.Context, pools []linodego.LKENodePool, declaredPools []interface{}) []linodego.LKENodePool {
	result := make([]linodego.LKENodePool, len(declaredPools))

	poolMap := make(map[int]linodego.LKENodePool, len(pools))
	for _, pool := range pools {
		poolMap[pool.ID] = pool
	}

	declaredPoolMap := make(map[int]map[string]any)
	for i, pool := range declaredPools {
		declaredPoolMap[i] = pool.(map[string]any)
	}

	// First, let's match any pools in state with an ID
	// TODO: Fix use of undefined behavior
	for i, declaredPool := range declaredPoolMap {
		poolID, ok := declaredPool["id"].(int)
		if !ok {
			continue
		}

		pool, ok := poolMap[poolID]
		if !ok {
			continue
		}

		result[i] = pool
		delete(poolMap, poolID)
		delete(declaredPoolMap, i)
	}

	// Second, let's match pools that have all matching attributes.
	// This is necessary because declared pools will not be populated with
	// an ID on first apply but still have matching node pools.
	for i, declaredPool := range declaredPoolMap {
		for _, pool := range poolMap {
			declaredAutoscaler := expandLinodeLKEClusterAutoscalerFromPool(declaredPool)

			if declaredPool["type"] != pool.Type {
				continue
			}

			if declaredPool["count"] != 0 && declaredPool["count"] != pool.Count {
				continue
			}

			// TODO: Make this less bad
			autoScalerEnabled := declaredAutoscaler != nil && declaredAutoscaler.Enabled
			if autoScalerEnabled != pool.Autoscaler.Enabled {
				continue
			}

			if declaredAutoscaler != nil && !reflect.DeepEqual(
				*declaredAutoscaler, pool.Autoscaler,
			) {
				continue
			}

			result[i] = pool
			delete(poolMap, pool.ID)
			break
		}
	}

	// Populate any additional pools
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

func expandLinodeLKENodePoolSpecs(pool []interface{}, preserveNoTarget bool) (poolSpecs []NodePoolSpec) {
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

		if !preserveNoTarget && specMap["id"].(int) == 0 {
			continue
		}

		poolSpecs = append(poolSpecs, NodePoolSpec{
			ID:                specMap["id"].(int),
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
