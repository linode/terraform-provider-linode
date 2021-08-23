package balancernode

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Schema:        resourceSchema,
		ReadContext:   readResource,
		CreateContext: createResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerNode ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	node, err := client.GetNodeBalancerNode(ctx, nodebalancerID, configID, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing NodeBalancer Node ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error finding the specified Linode NodeBalancerNode: %s", err)
	}

	d.Set("label", node.Label)
	d.Set("weight", node.Weight)
	d.Set("mode", node.Mode)
	d.Set("address", node.Address)
	d.Set("status", node.Status)
	return nil
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		// Validate that this is an ID by making sure it can be converted into an int
		_, err := strconv.Atoi(s[2])
		if err != nil {
			return nil, fmt.Errorf("invalid nodebalancer_node ID: %v", err)
		}

		configID, err := strconv.Atoi(s[1])
		if err != nil {
			return nil, fmt.Errorf("invalid config ID: %v", err)
		}

		nodebalancerID, err := strconv.Atoi(s[0])
		if err != nil {
			return nil, fmt.Errorf("invalid nodebalancer ID: %v", err)
		}

		d.SetId(s[2])
		d.Set("nodebalancer_id", nodebalancerID)
		d.Set("config_id", configID)
	}

	err := readResource(ctx, d, meta)
	if err != nil {
		return nil, fmt.Errorf("unable to import %v as nodebalancer_node: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	createOpts := linodego.NodeBalancerNodeCreateOptions{
		Address: d.Get("address").(string),
		Label:   d.Get("label").(string),
		Mode:    linodego.NodeMode(d.Get("mode").(string)),
		Weight:  d.Get("weight").(int),
	}
	node, err := client.CreateNodeBalancerNode(ctx, int(nodebalancerID), int(configID), createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode NodeBalancerNode: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", node.ID))
	d.Set("config_id", configID)
	d.Set("nodebalancer_id", nodebalancerID)

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerConfig ID %v as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}

	updateOpts := linodego.NodeBalancerNodeUpdateOptions{
		Address: d.Get("address").(string),
		Label:   d.Get("label").(string),
		Mode:    linodego.NodeMode(d.Get("mode").(string)),
		Weight:  d.Get("weight").(int),
	}

	if _, err = client.UpdateNodeBalancerNode(ctx, nodebalancerID, configID, int(id), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode Nodebalancer %d Config %d Node %d: %s",
			nodebalancerID, configID, int(id), err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode NodeBalancerConfig ID %s as int: %s", d.Id(), err)
	}
	nodebalancerID, ok := d.Get("nodebalancer_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("nodebalancer_id"))
	}
	configID, ok := d.Get("config_id").(int)
	if !ok {
		return diag.Errorf("Error parsing Linode NodeBalancer ID %v as int", d.Get("config_id"))
	}
	err = client.DeleteNodeBalancerNode(ctx, nodebalancerID, configID, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode NodeBalancerNode %d: %s", id, err)
	}
	return nil
}
