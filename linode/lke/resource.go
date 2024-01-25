package lke

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	k8scondition "github.com/linode/linodego/k8s/pkg/condition"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

const (
	createLKETimeout = 35 * time.Minute
	updateLKETimeout = 40 * time.Minute
	deleteLKETimeout = 15 * time.Minute
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(createLKETimeout),
			Update: schema.DefaultTimeout(updateLKETimeout),
			Delete: schema.DefaultTimeout(deleteLKETimeout),
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Read linode_lke_cluster")

	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode LKE Cluster ID: %s", err)
	}

	declaredPools, ok := d.Get("pool").([]interface{})
	if !ok {
		return diag.Errorf("failed to parse linode lke cluster pools: %d", id)
	}

	cluster, err := client.GetLKECluster(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing LKE Cluster ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("failed to get LKE cluster %d: %s", id, err)
	}

	pools, err := client.ListLKENodePools(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get pools for LKE cluster %d: %s", id, err)
	}

	kubeconfig, err := client.GetLKEClusterKubeconfig(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get kubeconfig for LKE cluster %d: %s", id, err)
	}

	endpoints, err := client.ListLKEClusterAPIEndpoints(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get API endpoints for LKE cluster %d: %s", id, err)
	}

	flattenedControlPlane := flattenLKEClusterControlPlane(cluster.ControlPlane)

	dashboard, err := client.GetLKEClusterDashboard(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get dashboard URL for LKE cluster %d: %s", id, err)
	}

	d.Set("label", cluster.Label)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("region", cluster.Region)
	d.Set("tags", cluster.Tags)
	d.Set("status", cluster.Status)
	d.Set("kubeconfig", kubeconfig.KubeConfig)
	d.Set("dashboard_url", dashboard.URL)
	d.Set("api_endpoints", flattenLKEClusterAPIEndpoints(endpoints))
	d.Set("pool", flattenLKENodePools(matchPoolsWithSchema(pools, declaredPools)))
	d.Set("control_plane", []map[string]interface{}{flattenedControlPlane})

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Create linode_lke_cluster")

	client := meta.(*helper.ProviderMeta).Client

	controlPlane := d.Get("control_plane").([]interface{})

	createOpts := linodego.LKEClusterCreateOptions{
		Label:      d.Get("label").(string),
		Region:     d.Get("region").(string),
		K8sVersion: d.Get("k8s_version").(string),
	}

	if len(controlPlane) > 0 {
		expandedControlPlane := expandLKEClusterControlPlane(controlPlane[0].(map[string]interface{}))
		createOpts.ControlPlane = &expandedControlPlane
	}

	for _, nodePool := range d.Get("pool").([]interface{}) {
		poolSpec := nodePool.(map[string]interface{})

		createOpts.NodePools = append(createOpts.NodePools, linodego.LKENodePoolCreateOptions{
			Type:       poolSpec["type"].(string),
			Count:      poolSpec["count"].(int),
			Autoscaler: expandLinodeLKEClusterAutoscalerFromPool(poolSpec),
		})
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	tflog.Debug(ctx, "Creating LKE cluster", map[string]any{
		"create_opts": createOpts,
	})

	cluster, err := client.CreateLKECluster(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create LKE cluster: %s", err)
	}
	d.SetId(strconv.Itoa(cluster.ID))

	tflog.Debug(ctx, "Waiting for a single LKE cluster node to be ready")

	// Sometimes the K8S API will raise an EOF error if polling immediately after
	// a cluster is created. We should retry accordingly.
	// NOTE: This routine has a short retry period because we want to raise
	// and meaningful errors quickly.
	diag.FromErr(retry.RetryContext(ctx, time.Second*25, func() *retry.RetryError {
		err := client.WaitForLKEClusterConditions(ctx, cluster.ID, linodego.LKEClusterPollOptions{
			TimeoutSeconds: 15 * 60,
		}, k8scondition.ClusterHasReadyNode)
		if err != nil {
			return retry.RetryableError(err)
		}

		return nil
	}))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Update linode_lke_cluster")

	providerMeta := meta.(*helper.ProviderMeta)
	client := providerMeta.Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed parsing Linode LKE Cluster ID: %s", err)
	}

	updateOpts := linodego.LKEClusterUpdateOptions{}

	if d.HasChange("label") {
		updateOpts.Label = d.Get("label").(string)
	}

	if d.HasChange("k8s_version") {
		updateOpts.K8sVersion = d.Get("k8s_version").(string)
	}

	controlPlane := d.Get("control_plane").([]interface{})
	if len(controlPlane) > 0 {
		expandedControlPlane := expandLKEClusterControlPlane(controlPlane[0].(map[string]interface{}))
		updateOpts.ControlPlane = &expandedControlPlane
	}

	if d.HasChange("tags") {
		tags := []string{}
		for _, tag := range d.Get("tags").(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}

		updateOpts.Tags = &tags
	}
	if d.HasChanges("label", "tags", "k8s_version", "control_plane") {
		tflog.Debug(ctx, "Updating LKE cluster", map[string]any{
			"update_opts": updateOpts,
		})

		if _, err := client.UpdateLKECluster(ctx, id, updateOpts); err != nil {
			return diag.Errorf("failed to update LKE Cluster %d: %s", id, err)
		}
	}

	pools, err := client.ListLKENodePools(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get Pools for LKE Cluster %d: %s", id, err)
	}

	if d.HasChange("k8s_version") {
		tflog.Debug(ctx, "Implicitly recycling LKE cluster to apply Kubernetes version upgrade")

		if err := recycleLKECluster(ctx, providerMeta, id, pools); err != nil {
			return diag.FromErr(err)
		}
	}

	poolSpecs := expandLinodeLKENodePoolSpecs(d.Get("pool").([]interface{}))
	updates := ReconcileLKENodePoolSpecs(poolSpecs, pools)

	tflog.Trace(ctx, "Reconciled LKE cluster node pool updates", map[string]any{
		"updates": updates,
	})

	updatedIds := []int{}

	for poolID, updateOpts := range updates.ToUpdate {
		tflog.Debug(ctx, "Updating LKE cluster node pool", map[string]any{
			"node_pool_id": poolID,
			"update_opts":  updateOpts,
		})

		if _, err := client.UpdateLKENodePool(ctx, id, poolID, updateOpts); err != nil {
			return diag.Errorf("failed to update LKE Cluster %d Pool %d: %s", id, poolID, err)
		}

		updatedIds = append(updatedIds, poolID)
	}

	for _, createOpts := range updates.ToCreate {
		tflog.Debug(ctx, "Creating LKE cluster node pool", map[string]any{
			"update_opts": updateOpts,
		})

		pool, err := client.CreateLKENodePool(ctx, id, createOpts)
		if err != nil {
			return diag.Errorf("failed to create LKE Cluster %d Pool: %s", id, err)
		}

		updatedIds = append(updatedIds, pool.ID)
	}

	for _, poolID := range updates.ToDelete {
		tflog.Debug(ctx, "Deleting LKE cluster node pool", map[string]any{
			"node_pool_id": poolID,
		})

		if err := client.DeleteLKENodePool(ctx, id, poolID); err != nil {
			return diag.Errorf("failed to delete LKE Cluster %d Pool %d: %s", id, poolID, err)
		}
	}

	tflog.Debug(ctx, "Waiting for all updated node pools to be ready")

	for _, poolID := range updatedIds {
		tflog.Trace(ctx, "Waiting for node pool to be ready", map[string]any{
			"node_pool_id": poolID,
		})

		if err := waitForNodePoolReady(
			ctx,
			client,
			providerMeta.Config.LKENodeReadyPollMilliseconds,
			id,
			poolID,
		); err != nil {
			return diag.Errorf("failed to wait for LKE Cluster %d pool %d ready: %s", id, poolID, err)
		}
	}

	return nil
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_lke_cluster")

	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed parsing Linode LKE Cluster ID: %s", err)
	}

	err = client.DeleteLKECluster(ctx, id)
	if err != nil {
		return diag.Errorf("failed to delete Linode LKE cluster %d: %s", id, err)
	}
	timeoutSeconds, err := helper.SafeFloat64ToInt(
		d.Timeout(schema.TimeoutCreate).Seconds(),
	)
	if err != nil {
		return diag.Errorf("failed to convert float64 creation timeout to int: %s", err)
	}

	tflog.Debug(ctx, "Deleted LKE cluster, waiting for all nodes deleted...")

	_, err = client.WaitForLKEClusterStatus(ctx, id, "not_ready", timeoutSeconds)
	if err != nil {
		// If we're getting a 404, it's safe to say the cluster has been
		// deleted.
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return nil
		}

		return diag.FromErr(err)
	}

	return nil
}

func flattenLKEClusterAPIEndpoints(apiEndpoints []linodego.LKEClusterAPIEndpoint) []string {
	flattened := make([]string, len(apiEndpoints))
	for i, endpoint := range apiEndpoints {
		flattened[i] = endpoint.Endpoint
	}
	return flattened
}

func populateLogAttributes(ctx context.Context, d *schema.ResourceData) context.Context {
	return helper.SetLogFieldBulk(ctx, map[string]any{
		"cluster_id": d.Id(),
	})
}
