package linode

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLinodeLKECluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceLinodeLKEClusterRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique label for the cluster.",
			},
			"k8s_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>. The latest supported patch version will be deployed.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "This cluster's location.",
			},
			"api_endpoints": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "The API endpoints for the cluster.",
			},
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The Base64-encoded Kubeconfig for the cluster.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the cluster.",
			},
			"pools": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of the Node Pool.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of nodes in the Node Pool.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A Linode Type for all of the nodes in the Node Pool.",
						},
						"nodes": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the node.",
									},
									"instance_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The ID of the underlying Linode instance.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The status of the node.`,
									},
								},
							},
							Computed:    true,
							Description: "The nodes in the node pool.",
						},
					},
				},
				Computed:    true,
				Description: "A node pool in the cluster.",
			},
		},
	}
}

func datasourceLinodeLKEClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	id := d.Get("id").(int)

	cluster, err := client.GetLKECluster(context.Background(), id)
	if err != nil {
		return diag.Errorf("failed to get LKE cluster %d: %s", id, err)
	}

	pools, err := client.ListLKEClusterPools(context.Background(), id, nil)
	if err != nil {
		return diag.Errorf("failed to get pools for LKE cluster %d: %s", id, err)
	}

	kubeconfig, err := client.GetLKEClusterKubeconfig(context.Background(), id)
	if err != nil {
		return diag.Errorf("failed to get kubeconfig for LKE cluster %d: %s", id, err)
	}

	endpoints, err := client.ListLKEClusterAPIEndpoints(context.Background(), id, nil)
	if err != nil {
		return diag.Errorf("failed to get API endpoints for LKE cluster %d: %s", id, err)
	}

	d.SetId(strconv.Itoa(id))
	d.Set("label", cluster.Label)
	d.Set("k8s_version", cluster.K8sVersion)
	d.Set("region", cluster.Region)
	d.Set("tags", cluster.Tags)
	d.Set("status", cluster.Status)
	d.Set("kubeconfig", kubeconfig.KubeConfig)
	d.Set("pools", flattenLinodeLKEClusterPools(pools))
	d.Set("api_endpoints", flattenLinodeLKEClusterAPIEndpoints(endpoints))
	return nil
}
