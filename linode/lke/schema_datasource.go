package lke

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var dataSourceSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of this LKE Cluster.",
		Required:    true,
	},

	"label": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The unique label for the cluster.",
	},
	"k8s_version": {
		Type:     schema.TypeString,
		Computed: true,
		Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>. " +
			"The latest supported patch version will be deployed.",
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
				"autoscaler": {
					Type:     schema.TypeList,
					Computed: true,
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
		Computed:    true,
		Description: "A node pool in the cluster.",
	},
	"control_plane": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"high_availability": {
					Type:        schema.TypeBool,
					Description: "Defines whether High Availability is enabled for the Control Plane Components of the cluster.",
					Computed:    true,
				},
			},
		},
		Computed:    true,
		Description: "Defines settings for the Kubernetes Control Plane.",
	},
}
