package lke

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/lkenodepool"
)

type NodePoolSpec struct {
	ID                int
	Tags              []string
	Type              string
	Count             int
	Taints            []map[string]any
	Labels            map[string]string
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
) (NodePoolUpdates, error) {
	result := NodePoolUpdates{
		ToCreate: make([]linodego.LKENodePoolCreateOptions, 0),
		ToUpdate: make(map[int]linodego.LKENodePoolUpdateOptions),
		ToDelete: make([]int, 0),
	}

	createPool := func(spec NodePoolSpec) error {
		createOpts := linodego.LKENodePoolCreateOptions{
			Count:  spec.Count,
			Type:   spec.Type,
			Tags:   spec.Tags,
			Labels: linodego.LKENodePoolLabels(spec.Labels),
		}

		if spec.Taints != nil {
			createOpts.Taints = expandNodePoolTaints(spec.Taints)
		}

		if createOpts.Count == 0 {
			if !spec.AutoScalerEnabled {
				return fmt.Errorf("count was 0 without an autoscaler. This is always a provider issue")
			}
			createOpts.Count = spec.AutoScalerMin
		}

		if spec.AutoScalerEnabled {
			createOpts.Autoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: true,
				Min:     spec.AutoScalerMin,
				Max:     spec.AutoScalerMax,
			}
		}

		result.ToCreate = append(result.ToCreate, createOpts)

		return nil
	}

	deletePool := func(id int) {
		result.ToDelete = append(result.ToDelete, id)
	}

	// If there are fewer node pools than expected
	// we can assume the rest have been deleted
	if len(newSpecs) < len(oldSpecs) {
		for _, v := range oldSpecs[len(newSpecs):] {
			deletePool(v.ID)
		}
	}

	// If there are more node pools then there were previously
	// we can assume new ones have been created
	if len(newSpecs) > len(oldSpecs) {
		for _, v := range newSpecs[len(oldSpecs):] {
			if err := createPool(v); err != nil {
				return result, err
			}
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
			if err := createPool(newSpec); err != nil {
				return result, err
			}

			deletePool(oldSpec.ID)
			continue
		}

		updateOpts := linodego.LKENodePoolUpdateOptions{
			Count: newSpec.Count,
			Tags:  &newSpecs[i].Tags,
		}

		if !helper.CompareSets(helper.TypedSliceToAny(newSpec.Taints), helper.TypedSliceToAny(oldSpec.Taints)) {
			taints := expandNodePoolTaints(newSpec.Taints)
			updateOpts.Taints = &taints
		}

		if !reflect.DeepEqual(newSpec.Labels, oldSpec.Labels) && !(len(newSpec.Labels) == 0 && len(oldSpec.Labels) == 0) {
			labels := linodego.LKENodePoolLabels(newSpecs[i].Labels)
			updateOpts.Labels = &labels
		}

		// Only include the autoscaler if the autoscaler has updated
		// This isn't stricly necessary but it makes unit testing easier
		if newSpec.AutoScalerEnabled != oldSpec.AutoScalerEnabled ||
			newSpec.AutoScalerMin != oldSpec.AutoScalerMin ||
			newSpec.AutoScalerMax != oldSpec.AutoScalerMax {
			updateOpts.Autoscaler = &linodego.LKENodePoolAutoscaler{
				Enabled: newSpec.AutoScalerEnabled,
				Min:     newSpec.AutoScalerMin,
				Max:     newSpec.AutoScalerMax,
			}
		}

		result.ToUpdate[oldSpec.ID] = updateOpts
	}

	return result, nil
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
		if _, err := lkenodepool.WaitForNodePoolReady(ctx, client, meta.Config.EventPollMilliseconds, id, pool.ID); err != nil {
			return fmt.Errorf("failed to wait for pool %d ready: %w", pool.ID, err)
		}
	}

	tflog.Debug(ctx, "All node pools have entered ready status; recycle operation completed")

	return nil
}

// This cannot currently be handled efficiently by a DiffSuppressFunc
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchPoolsWithSchema(ctx context.Context, pools []linodego.LKENodePool, declaredPools []interface{}) ([]linodego.LKENodePool, error) {
	tflog.Info(ctx, "Enter matchPoolsWithSchema helper function")
	result := make([]linodego.LKENodePool, len(declaredPools))

	// Contains all unpaired pools returned by the API
	apiPools := make(map[int]linodego.LKENodePool, len(pools))
	for _, pool := range pools {
		apiPools[pool.ID] = pool
	}

	// Tracks which local pools have been processed
	pairedDeclaredPools := make(map[int]bool)

	// First let's match any pools in state with an ID
	for i, declaredPool := range declaredPools {
		declaredPool := declaredPool.(map[string]any)

		poolID, ok := declaredPool["id"].(int)
		if !ok {
			return nil, fmt.Errorf("declared pool ID was not of type int")
		}

		apiPool, ok := apiPools[poolID]
		if !ok {
			continue
		}

		// Pair the found pool with the declared pool
		result[i] = apiPool
		delete(apiPools, poolID)
		pairedDeclaredPools[i] = true
	}

	// Second, let's match pools that have all matching attributes.
	// This is necessary because declared pools will not be populated with
	// an ID on first apply but still have matching node pools.
	for i, declaredPool := range declaredPools {
		declaredPool := declaredPool.(map[string]any)
		declaredAutoscaler := expandLinodeLKEClusterAutoscalerFromPool(declaredPool)

		if _, ok := pairedDeclaredPools[i]; ok {
			// This apiPool has already been handled in the previous step,
			// we can skip it
			continue
		}

		for _, apiPool := range apiPools {
			if declaredPool["type"] != apiPool.Type {
				continue
			}

			declaredCount := declaredPool["count"].(int)
			if declaredCount == 0 {
				if declaredAutoscaler == nil {
					return nil, fmt.Errorf("autoscaler is null when count is 0. This is always a provider issue")
				}
				declaredCount = declaredAutoscaler.Min
			}

			if declaredCount != apiPool.Count {
				continue
			}

			if (declaredAutoscaler != nil && declaredAutoscaler.Enabled) != apiPool.Autoscaler.Enabled {
				continue
			}

			// Only compare autoscalers if the declared autoscaler is enabled
			if declaredAutoscaler != nil && !reflect.DeepEqual(
				*declaredAutoscaler, apiPool.Autoscaler,
			) {
				continue
			}

			if !helper.CompareStringSets(helper.ExpandStringSet(declaredPool["tags"].(*schema.Set)), apiPool.Tags) {
				continue
			}

			declaredTaints := expandNodePoolTaints(helper.ExpandObjectSet(declaredPool["taint"].(*schema.Set)))

			if !helper.CompareSets(helper.TypedSliceToAny(declaredTaints), helper.TypedSliceToAny(apiPool.Taints)) {
				continue
			}

			declaredLabels := helper.StringAnyMapToTyped[string](declaredPool["labels"].(map[string]any))

			// - Length comparison is for handling the case of nil vs empty slice
			// - Converting `apiPool.Labels` back to original (non-alias) type to make `reflect.DeepEqual` to really compare them
			if !reflect.DeepEqual(declaredLabels, map[string]string(apiPool.Labels)) &&
				!(len(declaredLabels) == 0 && len(apiPool.Labels) == 0) {
				continue
			}

			// Pair the API pool with the declared pool
			result[i] = apiPool
			delete(apiPools, apiPool.ID)
			break
		}
	}

	// Append any unresolved pools to the end
	// These are typically pools planned to be deleted
	for _, pool := range apiPools {
		//nolint:makezero
		result = append(result, pool)
	}

	return result, nil
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
			Tags:              helper.ExpandStringSet(specMap["tags"].(*schema.Set)),
			Taints:            helper.ExpandObjectSet(specMap["taint"].(*schema.Set)),
			Labels:            helper.StringAnyMapToTyped[string](specMap["labels"].(map[string]any)),
			Count:             specMap["count"].(int),
			AutoScalerEnabled: autoscaler.Enabled,
			AutoScalerMin:     autoscaler.Min,
			AutoScalerMax:     autoscaler.Max,
		})
	}
	return
}

func flattenLKENodePools(pools []linodego.LKENodePool) []map[string]interface{} {
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
			"id":              pool.ID,
			"count":           pool.Count,
			"type":            pool.Type,
			"tags":            pool.Tags,
			"disk_encryption": pool.DiskEncryption,
			"taint":           flattenNodePoolTaints(pool.Taints),
			"labels":          pool.Labels,
			"nodes":           nodes,
			"autoscaler":      autoscaler,
		}
	}
	return flattened
}

func flattenNodePoolTaints(taints []linodego.LKENodePoolTaint) []map[string]string {
	result := make([]map[string]string, len(taints))

	for i, t := range taints {
		result[i] = map[string]string{
			"effect": string(t.Effect),
			"key":    t.Key,
			"value":  t.Value,
		}
	}

	return result
}

func flattenLKEClusterControlPlane(controlPlane linodego.LKEClusterControlPlane, aclResp *linodego.LKEClusterControlPlaneACLResponse) map[string]interface{} {
	flattened := make(map[string]any)
	if aclResp != nil {
		acl := aclResp.ACL

		flattenedACL := make(map[string]any)
		flattenedACL["enabled"] = acl.Enabled

		if acl.Addresses != nil {
			flattenedAddresses := map[string]any{
				"ipv4": acl.Addresses.IPv4,
				"ipv6": acl.Addresses.IPv6,
			}
			flattenedACL["addresses"] = []map[string]any{flattenedAddresses}
		}

		flattened["acl"] = []map[string]any{flattenedACL}
	}

	flattened["high_availability"] = controlPlane.HighAvailability

	return flattened
}

func expandControlPlaneOptions(controlPlane map[string]interface{}) (
	result linodego.LKEClusterControlPlaneOptions,
	diags diag.Diagnostics,
) {
	if value, ok := controlPlane["high_availability"]; ok {
		v := value.(bool)
		result.HighAvailability = &v
	}

	// default to disabled
	enabled := false
	result.ACL = &linodego.LKEClusterControlPlaneACLOptions{Enabled: &enabled}

	if value, ok := controlPlane["acl"]; ok {
		v := value.([]any)
		if len(v) > 0 {
			result.ACL, diags = expandACLOptions(v[0].(map[string]interface{}))
			if diags.HasError() {
				return
			}
		}
	}

	return
}

func expandACLOptions(aclOptions map[string]interface{}) (*linodego.LKEClusterControlPlaneACLOptions, diag.Diagnostics) {
	var result linodego.LKEClusterControlPlaneACLOptions

	if value, ok := aclOptions["enabled"]; ok {
		v := value.(bool)
		result.Enabled = &v
	}

	if value, ok := aclOptions["addresses"]; ok {
		v := value.([]interface{})
		if len(v) > 0 {
			result.Addresses = expandACLAddressOptions(v[0].(map[string]interface{}))
		}
	}

	if (result.Enabled != nil && !*result.Enabled) &&
		(result.Addresses != nil &&
			((result.Addresses.IPv4 != nil && len(*result.Addresses.IPv4) > 0) ||
				(result.Addresses.IPv6 != nil && len(*result.Addresses.IPv6) > 0))) {
		return nil, diag.Errorf("addresses are not acceptable when ACL is disabled")
	}

	return &result, nil
}

func expandACLAddressOptions(addressOptions map[string]interface{}) *linodego.LKEClusterControlPlaneACLAddressesOptions {
	var result linodego.LKEClusterControlPlaneACLAddressesOptions

	if value, ok := addressOptions["ipv4"]; ok {
		ipv4 := helper.ExpandStringSet(value.(*schema.Set))
		result.IPv4 = &ipv4
	}

	if value, ok := addressOptions["ipv6"]; ok {
		ipv6 := helper.ExpandStringSet(value.(*schema.Set))
		result.IPv6 = &ipv6
	}

	return &result
}

func filterExternalPools(ctx context.Context, externalPoolTags []string, pools []linodego.LKENodePool) []linodego.LKENodePool {
	var filteredPools []linodego.LKENodePool
	if len(externalPoolTags) == 0 {
		return pools
	}
	tagSet := make(map[string]bool, len(externalPoolTags))
	for _, tag := range externalPoolTags {
		tagSet[tag] = true
	}
	for _, pool := range pools {
		tag := poolHasAnyOfTags(pool, tagSet)
		if tag != nil {
			tflog.Info(ctx, "Excluding pool from management by this resource", map[string]interface{}{
				"pool_id": pool.ID,
				"tag":     tag,
				"reason":  "Pool tagged to be managed by a separate linode_lke_node_pool resource",
			})
			continue
		}
		filteredPools = append(filteredPools, pool)
	}
	return filteredPools
}

func poolHasAnyOfTags(pool linodego.LKENodePool, tagSet map[string]bool) *string {
	for _, poolTag := range pool.Tags {
		if _, exists := tagSet[poolTag]; exists {
			result := poolTag
			return &result
		}
	}
	return nil
}

func expandNodePoolTaints(poolTaints []map[string]any) []linodego.LKENodePoolTaint {
	taints := make([]linodego.LKENodePoolTaint, len(poolTaints))
	for i, taint := range poolTaints {
		taints[i] = linodego.LKENodePoolTaint{
			Key:    taint["key"].(string),
			Value:  taint["value"].(string),
			Effect: linodego.LKENodePoolTaintEffect(taint["effect"].(string)),
		}
	}
	return taints
}

func waitForLKEKubeConfig(ctx context.Context, client linodego.Client, intervalMS int, clusterID int) error {
	ticker := time.NewTicker(time.Duration(intervalMS) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := client.GetLKEClusterKubeconfig(ctx, clusterID)
			if err != nil {
				if strings.Contains(err.Error(), "Cluster kubeconfig is not yet available") {
					continue
				} else {
					return fmt.Errorf("failed to get Kubeconfig for LKE cluster %d: %w", clusterID, err)
				}
			} else {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("Error waiting for Cluster %d kubeconfig: %w", clusterID, ctx.Err())
		}
	}
}
