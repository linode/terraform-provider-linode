package lke

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type ClusterPoolSpec struct {
	Type  string
	Count int
}

type linodelkeClusterPoolUpdates struct {
	ToDelete []int
	ToCreate []linodego.LKEClusterPoolCreateOptions
	ToUpdate map[int]linodego.LKEClusterPoolUpdateOptions
}

type clusterPoolAssignRequest struct {
	Spec, State ClusterPoolSpec
	PoolID      int
	SpecIndex   int
}

func (r clusterPoolAssignRequest) Diff() int {
	return int(math.Abs(float64(r.State.Count - r.Spec.Count)))
}

func expandLinodeLKEClusterPoolSpecs(pool []interface{}) (poolSpecs []ClusterPoolSpec) {
	for _, spec := range pool {
		specMap := spec.(map[string]interface{})
		poolSpecs = append(poolSpecs, ClusterPoolSpec{
			Type:  specMap["type"].(string),
			Count: specMap["count"].(int),
		})
	}
	return
}

func getLKEClusterPoolProvisionedSpecs(pools []linodego.LKEClusterPool) map[ClusterPoolSpec]map[int]struct{} {
	provisioned := make(map[ClusterPoolSpec]map[int]struct{})
	for _, pool := range pools {
		spec := ClusterPoolSpec{
			Type:  pool.Type,
			Count: pool.Count,
		}
		if _, ok := provisioned[spec]; !ok {
			provisioned[spec] = make(map[int]struct{})
		}
		provisioned[spec][pool.ID] = struct{}{}
	}
	return provisioned
}

func ReconcileLKEClusterPoolSpecs(
	poolSpecs []ClusterPoolSpec, pools []linodego.LKEClusterPool) (updates linodelkeClusterPoolUpdates) {
	provisionedPools := getLKEClusterPoolProvisionedSpecs(pools)
	poolSpecsToAssign := make(map[int]struct{})
	assignedPools := make(map[int]struct{})
	updates.ToUpdate = make(map[int]linodego.LKEClusterPoolUpdateOptions)

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
	poolAssignRequests := []clusterPoolAssignRequest{}
	for i := range poolSpecsToAssign {
		poolSpec := poolSpecs[i]
		for pool := range provisionedPools {
			if pool.Type != poolSpec.Type {
				continue
			}

			for id := range provisionedPools[pool] {
				poolAssignRequests = append(poolAssignRequests, clusterPoolAssignRequest{
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

		updates.ToUpdate[request.PoolID] = linodego.LKEClusterPoolUpdateOptions{
			Count: request.Spec.Count,
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
		updates.ToCreate = append(updates.ToCreate, linodego.LKEClusterPoolCreateOptions{
			Count: poolSpec.Count,
			Type:  poolSpec.Type,
		})
	}

	for spec := range provisionedPools {
		for id := range provisionedPools[spec] {
			updates.ToDelete = append(updates.ToDelete, id)
		}
	}

	return
}

func waitForClusterPoolReady(
	ctx context.Context, client *linodego.Client, errCh chan<- error, wg *sync.WaitGroup, pollMs, clusterID, poolID int) {
	eventTicker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)

main:
	for {
		select {
		case <-ctx.Done():
			log.Printf("[ERROR] timed out waiting for LKE Cluster (%d) Pool (%d) to be ready", clusterID, poolID)
			return

		case <-eventTicker.C:
			pool, err := client.GetLKEClusterPool(ctx, clusterID, poolID)
			if err != nil {
				errCh <- fmt.Errorf("failed to get LKE Cluster (%d) Pool (%d): %w", clusterID, poolID, err)
			}

			for _, instance := range pool.Linodes {
				if instance.Status == linodego.LKELinodeNotReady {
					continue main
				}
			}

			log.Printf("[DEBUG] finished waiting for LKE Cluster (%d) Pool (%d) to be ready", clusterID, poolID)
			wg.Done()
			return
		}
	}
}

func waitForClusterPoolsToStartRecycle(
	ctx context.Context, client *linodego.Client, pollMs, clusterID int, pools []linodego.LKEClusterPool,
) (<-chan int, <-chan error) {
	clusterInstances := make(map[int]int)
	poolInstances := make(map[int]map[int]struct{}, len(pools))
	for _, pool := range pools {
		poolInstances[pool.ID] = make(map[int]struct{}, len(pool.Linodes))
		for _, instance := range pool.Linodes {
			poolInstances[pool.ID][instance.InstanceID] = struct{}{}
			clusterInstances[instance.InstanceID] = pool.ID
		}
	}

	poolRecyclesCh := make(chan int)
	errCh := make(chan error)

	eventTicker := time.NewTicker(time.Duration(pollMs) * time.Millisecond)

	go func() {
		defer eventTicker.Stop()
		defer close(poolRecyclesCh)
		defer close(errCh)

		lastEventID := 0

		for len(clusterInstances) != 0 {
			select {
			case <-ctx.Done():
				log.Printf("[ERROR] timed out waiting for all original nodes of LKE Cluster (%d) to be deleted (%d remaining)\n",
					clusterID, len(clusterInstances))
				return

			case <-eventTicker.C:
				filterBytes, _ := json.Marshal(map[string]interface{}{
					"+order_by":   "created",
					"+gte":        lastEventID,
					"entity.type": string(linodego.EntityLinode),
				})

				events, err := client.ListEvents(ctx, linodego.NewListOptions(1, string(filterBytes)))
				if err != nil {
					errCh <- err
					return
				}

				if len(events) != 0 {
					lastEventID = events[0].ID
				}

				for _, event := range events {
					if event.Action != linodego.ActionLinodeDelete {
						continue
					}
					id := int(event.Entity.ID.(float64))
					poolID, ok := clusterInstances[id]
					if !ok {
						continue
					}

					delete(clusterInstances, id)
					delete(poolInstances[poolID], id)
					log.Printf("[DEBUG] finished waiting for LKE Cluster (%d) Pool (%d) Node (%d) to be deleted\n",
						clusterID, poolID, id)

					if len(poolInstances[poolID]) == 0 {
						// all original instances for this pool have been deleted
						delete(poolInstances, poolID)
						log.Printf("[DEBUG] finished waiting for all nodes in LKE Cluster (%d) Pool (%d) to be recreated\n",
							clusterID, poolID)
						poolRecyclesCh <- poolID
					}
				}
			}
		}
	}()
	return poolRecyclesCh, errCh
}

func recycleLKECluster(ctx context.Context, meta *helper.ProviderMeta, id int, pools []linodego.LKEClusterPool) error {
	client := meta.Client

	if err := client.RecycleLKEClusterNodes(ctx, id); err != nil {
		return fmt.Errorf("failed to recycle LKE Cluster (%d): %s", id, err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	poolRecyclesCh, errCh := waitForClusterPoolsToStartRecycle(
		ctx, &client, meta.Config.LKEEventPollMilliseconds, id, pools)

	var wg sync.WaitGroup
	wg.Add(len(pools))
	poolsRecycledCh := waitGroupCh(&wg)

	readyErrCh := make(chan error)
	defer close(readyErrCh)

	go func() {
		for poolID := range poolRecyclesCh {
			go waitForClusterPoolReady(ctx, &client, readyErrCh, &wg, meta.Config.LKENodeReadyPollMilliseconds, id, poolID)
		}
	}()

	for {
		select {
		case <-poolsRecycledCh:
			return nil

		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("failed to wait for all LKE Cluster (%d) nodes to start recycle: %w", id, err)
			}

		case err := <-readyErrCh:
			if err != nil {
				return err
			}
		}
	}
}

func waitGroupCh(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

// This cannot currently be handled efficiently by a DiffSuppressFunc
// See: https://github.com/hashicorp/terraform-plugin-sdk/issues/477
func matchPoolsWithSchema(pools []linodego.LKEClusterPool, declaredPools []interface{}) []linodego.LKEClusterPool {
	result := make([]linodego.LKEClusterPool, len(declaredPools))

	poolMap := make(map[int]linodego.LKEClusterPool, len(declaredPools))
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
		result = append(result, pool)
	}

	return result
}
