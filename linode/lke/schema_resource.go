package lke

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

var resourceSchema = map[string]*schema.Schema{
	"label": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The unique label for the cluster.",
	},
	"k8s_version": {
		Type:     schema.TypeString,
		Required: true,
		Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>. " +
			"The latest supported patch version will be deployed.",
	},
	"tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Computed:    true,
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
	},
	"external_pool_tags": {
		Type:        schema.TypeSet,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Description: "An array of tags indicating that node pools having those tags are defined with a separate nodepool resource, rather than inside the current cluster resource.",
	},
	"region": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
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
	"dashboard_url": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The dashboard URL of the cluster.",
	},
	"status": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The status of the cluster.",
	},
	"pool": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeInt,
					Computed:    true,
					Description: "The ID of the Node Pool.",
				},
				"count": {
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
					Description:  "The number of nodes in the Node Pool.",
					Optional:     true,
					Computed:     true,
				},
				"type": {
					Type:        schema.TypeString,
					Description: "A Linode Type for all of the nodes in the Node Pool.",
					Required:    true,
				},
				"tags": {
					Type:        schema.TypeSet,
					Description: "A set of tags applied to this node pool.",

					Elem:                  &schema.Schema{Type: schema.TypeString},
					Optional:              true,
					DiffSuppressOnRefresh: true,
					DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
						setAttrName := k[0:strings.LastIndex(k, ".")]

						index, err := strconv.Atoi(k[strings.LastIndex(k, "."):])
						if err != nil {
							return false
						}

						tags := helper.ExpandStringSet(d.Get(setAttrName).(*schema.Set))
						for i, item := range tags {
							if i == index {
								continue
							}
							if strings.EqualFold(item, oldValue) || strings.EqualFold(item, newValue) {
								return true
							}
						}
						return false
					},
				},
				"disk_encryption": {
					Type: schema.TypeString,
					Description: "The disk encryption policy for the nodes in this pool. " +
						"NOTE: Disk encryption may not currently be available to all users.",
					Computed: true,
				},
				"nodes": {
					Type: schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:        schema.TypeString,
								Description: "The ID of the node.",
								Computed:    true,
							},
							"instance_id": {
								Type:        schema.TypeInt,
								Description: "The ID of the underlying Linode instance.",
								Computed:    true,
							},
							"status": {
								Type:        schema.TypeString,
								Description: `The status of the node.`,
								Computed:    true,
							},
						},
					},
					Computed:    true,
					Description: "The nodes in the node pool.",
				},
				"autoscaler": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"min": {
								Type:        schema.TypeInt,
								Description: "The minimum number of nodes to autoscale to.",
								Required:    true,
							},
							"max": {
								Type:        schema.TypeInt,
								Description: "The maximum number of nodes to autoscale to.",
								Required:    true,
							},
						},
					},
					Description: "When specified, the number of nodes autoscales within " +
						"the defined minimum and maximum values.",
				},
			},
		},
		MinItems:    1,
		Required:    true,
		Description: "A node pool in the cluster.",
	},
	"control_plane": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Computed:    true,
		Description: "Defines settings for the Kubernetes Control Plane.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"high_availability": {
					Type:        schema.TypeBool,
					Description: "Defines whether High Availability is enabled for the Control Plane Components of the cluster.",
					Optional:    true,
					Computed:    true,
				},
				"acl": {
					Type:        schema.TypeList,
					Description: "Defines the ACL configuration for an LKE cluster's control plane.",
					Optional:    true,
					Computed:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enabled": {
								Type:        schema.TypeBool,
								Description: "Defines default policy. A value of true results in a default policy of DENY. A value of false results in default policy of ALLOW, and has the same effect as delete the ACL configuration.",
								Computed:    true,
								Optional:    true,
							},
							"addresses": {
								Type:        schema.TypeList,
								Description: "A list of ip addresses to allow.",
								Optional:    true,
								Computed:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"ipv4": {
											Type:        schema.TypeSet,
											Description: "A set of individual ipv4 addresses or CIDRs to ALLOW.",
											Optional:    true,
											Computed:    true,
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
										"ipv6": {
											Type:        schema.TypeSet,
											Description: "A set of individual ipv6 addresses or CIDRs to ALLOW.",
											Optional:    true,
											Computed:    true,
											Elem:        &schema.Schema{Type: schema.TypeString},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
