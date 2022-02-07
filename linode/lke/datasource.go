package lke

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id := d.Get("id").(int)

	cluster, err := client.GetLKECluster(ctx, id)
	if err != nil {
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

	d.SetId(strconv.Itoa(id))
	d.Set("label", cluster.Label)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("region", cluster.Region)
	d.Set("tags", cluster.Tags)
	d.Set("status", cluster.Status)
	d.Set("kubeconfig", kubeconfig.KubeConfig)
	d.Set("pools", flattenLKENodePools(pools))
	d.Set("api_endpoints", flattenLKEClusterAPIEndpoints(endpoints))
	d.Set("control_plane", []interface{}{flattenedControlPlane})
	return nil
}
