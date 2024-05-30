package lke

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	k8scondition "github.com/linode/linodego/k8s/pkg/condition"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	linodediffs "github.com/linode/terraform-provider-linode/v2/linode/helper/customdiffs"
	"github.com/linode/terraform-provider-linode/v2/linode/lkenodepool"
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
		CustomizeDiff: customdiff.All(
			customDiffValidateOptionalCount,
			linodediffs.ComputedWithDefault("tags", []string{}),
			linodediffs.CaseInsensitiveSet("tags"),
		),
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

	tflog.Trace(ctx, "client.ListLKENodePools(...)")
	pools, err := client.ListLKENodePools(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get pools for LKE cluster %d: %s", id, err)
	}

	externalPoolTags := helper.ExpandStringSet(d.Get("external_pool_tags").(*schema.Set))
	if len(externalPoolTags) > 0 && len(pools) > 0 {
		pools = filterExternalPools(ctx, externalPoolTags, pools)
	}

	kubeconfig, err := client.GetLKEClusterKubeconfig(ctx, id)
	if err != nil {
		return diag.Errorf("failed to get kubeconfig for LKE cluster %d: %s", id, err)
	}

	tflog.Trace(ctx, "client.ListLKEClusterAPIEndpoints(...)")
	endpoints, err := client.ListLKEClusterAPIEndpoints(ctx, id, nil)
	if err != nil {
		return diag.Errorf("failed to get API endpoints for LKE cluster %d: %s", id, err)
	}

	acl, err := client.GetLKEClusterControlPlaneACL(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok &&
			(lerr.Code == 404 ||
				(lerr.Code == 400 && strings.Contains(lerr.Message, "Cluster does not support Control Plane ACL"))) {
			// The customer doesn't have access to LKE or the cluster does not have a Gateway. Nothing to do here.
		} else {
			return diag.Errorf("failed to get control plane ACL for LKE cluster %d: %s", id, err)
		}
	}

	flattenedControlPlane := flattenLKEClusterControlPlane(cluster.ControlPlane, acl)

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

	matchedPools, err := matchPoolsWithSchema(pools, declaredPools)
	if err != nil {
		return diag.Errorf("failed to match api pools with schema: %s", err)
	}

	p := flattenLKENodePools(matchedPools)

	d.Set("pool", p)
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
		expandedControlPlane := expandControlPlaneOptions(controlPlane[0].(map[string]interface{}))
		createOpts.ControlPlane = &expandedControlPlane
	}

	for _, nodePool := range d.Get("pool").([]interface{}) {
		poolSpec := nodePool.(map[string]interface{})

		autoscaler := expandLinodeLKEClusterAutoscalerFromPool(poolSpec)

		// If the count is not explicitly defined,
		// we should default it to the autoscaler minimum.
		count := poolSpec["count"].(int)
		if count == 0 {
			// We have validation to prevent this, but just in-case!
			if autoscaler == nil {
				return diag.Errorf(
					"Expected autoscaler for default node count, got nil. " +
						"This is always a provider issue.",
				)
			}

			count = autoscaler.Min
		}

		createOpts.NodePools = append(createOpts.NodePools, linodego.LKENodePoolCreateOptions{
			Type:       poolSpec["type"].(string),
			Count:      count,
			Autoscaler: autoscaler,
		})
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		tags := helper.ExpandStringSet(tagsRaw.(*schema.Set))
		createOpts.Tags = tags
	}

	tflog.Debug(ctx, "client.CreateLKECluster(...)", map[string]any{
		"options": createOpts,
	})
	cluster, err := client.CreateLKECluster(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create LKE cluster: %s", err)
	}
	d.SetId(strconv.Itoa(cluster.ID))

	ctx = tflog.SetField(ctx, "cluster_id", cluster.ID)
	tflog.Debug(ctx, "Waiting for a single LKE cluster node to be ready")

	// Sometimes the K8S API will raise an EOF error if polling immediately after
	// a cluster is created. We should retry accordingly.
	// NOTE: This routine has a short retry period because we want to raise
	// and meaningful errors quickly.
	diag.FromErr(retry.RetryContext(ctx, time.Second*25, func() *retry.RetryError {
		tflog.Trace(ctx, "client.WaitForLKEClusterCondition(...)", map[string]any{
			"condition": "ClusterHasReadyNode",
		})

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
		expandedControlPlane := expandControlPlaneOptions(controlPlane[0].(map[string]interface{}))
		updateOpts.ControlPlane = &expandedControlPlane
	}

	if d.HasChange("tags") {
		tags := helper.ExpandStringSet(d.Get("tags").(*schema.Set))
		updateOpts.Tags = &tags
	}
	if d.HasChanges("label", "tags", "k8s_version", "control_plane") {
		tflog.Debug(ctx, "client.UpdateLKECluster(...)", map[string]any{
			"options": updateOpts,
		})

		if _, err := client.UpdateLKECluster(ctx, id, updateOpts); err != nil {
			return diag.Errorf("failed to update LKE Cluster %d: %s", id, err)
		}
	}

	tflog.Trace(ctx, "client.ListLKENodePools(...)")

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

	oldPools, newPools := d.GetChange("pool")

	updates, err := ReconcileLKENodePoolSpecs(
		expandLinodeLKENodePoolSpecs(oldPools.([]any), false),
		expandLinodeLKENodePoolSpecs(newPools.([]any), true),
	)
	if err != nil {
		return diag.Errorf("Failed to reconcile LKE cluster node pools: %s", err)
	}

	tflog.Trace(ctx, "Reconciled LKE cluster node pool updates", map[string]any{
		"updates": updates,
	})

	updatedIds := []int{}

	for poolID, updateOpts := range updates.ToUpdate {
		tflog.Debug(ctx, "client.UpdateLKENodePool(...)", map[string]any{
			"node_pool_id": poolID,
			"options":      updateOpts,
		})

		if _, err := client.UpdateLKENodePool(ctx, id, poolID, updateOpts); err != nil {
			return diag.Errorf("failed to update LKE Cluster %d Pool %d: %s", id, poolID, err)
		}

		updatedIds = append(updatedIds, poolID)
	}

	for _, createOpts := range updates.ToCreate {
		tflog.Debug(ctx, "client.CreateLKENodePool(...)", map[string]any{
			"options": updateOpts,
		})
		pool, err := client.CreateLKENodePool(ctx, id, createOpts)
		if err != nil {
			return diag.Errorf("failed to create LKE Cluster %d Pool: %s", id, err)
		}

		updatedIds = append(updatedIds, pool.ID)
	}

	for _, poolID := range updates.ToDelete {
		tflog.Debug(ctx, "client.DeleteLKENodePool(...)", map[string]any{
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

		if _, err := lkenodepool.WaitForNodePoolReady(
			ctx,
			client,
			providerMeta.Config.LKENodeReadyPollMilliseconds,
			id,
			poolID,
		); err != nil {
			return diag.Errorf("failed to wait for LKE Cluster %d pool %d ready: %s", id, poolID, err)
		}
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx = populateLogAttributes(ctx, d)
	tflog.Debug(ctx, "Delete linode_lke_cluster")

	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("failed parsing Linode LKE Cluster ID: %s", err)
	}

	tflog.Debug(ctx, "client.DeleteLKECluster(...)")
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
	tflog.Trace(ctx, "client.WaitForLKEClusterStatus(...)", map[string]any{
		"status":  "not_ready",
		"timeout": timeoutSeconds,
	})

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

// customDiffValidateOptionalCount ensures an autoscaler must be
// defined is count is undefined.
//
// This validation logic is implemented as a custom diff because
// ValidateDiagFuncs are not currently supported directly on lists.
//
// Additionally, this validation is implemented using cty so we
// can ensure we're only validating on the user's config rather
// than state. This will prevent any possible false-negatives
// during updates.
func customDiffValidateOptionalCount(ctx context.Context, diff *schema.ResourceDiff, meta any) error {
	invalidPools := make([]string, 0)

	poolIterator := diff.GetRawConfig().GetAttr("pool").ElementIterator()

	for poolIterator.Next() {
		rawKey, rawPool := poolIterator.Element()

		// If the user has defined a count, we don't need to do anything special here
		if !rawPool.GetAttr("count").IsNull() {
			continue
		}

		// If the user hasn't defined a count but has defined an autoscaler,
		// we can assume they're deferring the count to the autoscaler.
		autoscaler := rawPool.GetAttr("autoscaler")

		if !autoscaler.IsNull() && autoscaler.LengthInt() > 0 {
			continue
		}

		// We need to use AsBigFloat to extract a number
		// value from a cty.Value
		index, _ := rawKey.AsBigFloat().Int64()

		invalidPools = append(invalidPools, fmt.Sprintf("pool.%d", index))
	}

	if len(invalidPools) > 0 {
		return fmt.Errorf(
			"%s: `count` must be defined when no autoscaler is defined",
			strings.Join(invalidPools, ", "),
		)
	}

	return nil
}
