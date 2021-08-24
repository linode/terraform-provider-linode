package lke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
		Description: "An array of tags applied to this object. Tags are for organizational purposes only.",
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
					Required:     true,
				},
				"type": {
					Type:        schema.TypeString,
					Description: "A Linode Type for all of the nodes in the Node Pool.",
					Required:    true,
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
			},
		},
		MinItems:    1,
		Required:    true,
		Description: "A node pool in the cluster.",
	},
}
