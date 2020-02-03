package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/linode/linodego"
)

func dataSourceLinodeLKECluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinodeLKEClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Description: "This Kubernetes cluster's unique ID.",
				Required:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "When this Kubernetes cluster was created.",
				Computed:    true,
			},
			"label": {
				Type:        schema.TypeString,
				Description: "This Kubernetes cluster's unique label for display purposes only. If no label is provided, one will be assigned automatically.",
				Computed:    true,
			},
			"tags": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "",
				Computed:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "This Kubernetes cluster's location.",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "The desired Kubernetes version for this Kubernetes cluster in the format of <major>.<minor>, and the latest supported patch version will be deployed.",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "When this Kubernetes cluster was updated.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLinodeLKEClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	reqLKEID := d.Get("id").(int)

	if reqLKEID == 0 {
		return fmt.Errorf("ID of the Kubernetes cluster is required")
	}

	var LKEcluster *linodego.LKECluster
	LKEcluster, err := client.GetLKECluster(context.Background(), reqLKEID)
	if err != nil {
		return fmt.Errorf("error listing LKE cluster: %s", err)
	}

	if LKEcluster != nil {
		d.SetId(strconv.Itoa(LKEcluster.ID))
		d.Set("created", LKEcluster.Created)
		d.Set("label", LKEcluster.Label)
		d.Set("region", LKEcluster.Region)
		d.Set("version", LKEcluster.Version)
		d.Set("status", LKEcluster.Status)
		d.Set("updated", LKEcluster.Updated)
		if err := d.Set("tags", LKEcluster.Tags); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}

		return nil
	}

	return fmt.Errorf("LKE cluster %s was not found", string(reqLKEID))
}
