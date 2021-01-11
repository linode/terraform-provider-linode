package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeVLAN() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeVLANCreate,
		Read:   resourceLinodeVLANRead,
		Update: resourceLinodeVLANUpdate,
		Delete: resourceLinodeVLANDelete,
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the vlan for display purposes only.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region where the vlan is deployed.",
			},
			"attached_linodes": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipv4_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed:    true,
				Description: "The Linodes attached to this vlan.",
			},
			"linodes": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Description: "The IDs of the Linodes to attach to this vlan.",
			},
			"cidr_block": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceLinodeVLANCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	createOpts := linodego.VLANCreateOptions{
		Description: d.Get("description").(string),
		Region:      d.Get("region").(string),
		CIDRBlock:   d.Get("cidr_block").(string),
		Linodes:     expandIntSet(d.Get("linodes").(*schema.Set)),
	}

	vlan, err := client.CreateVLAN(context.Background(), createOpts)
	if err != nil {
		return fmt.Errorf("failed to create vlan: %s", err)
	}
	d.SetId(strconv.Itoa(vlan.ID))
	return resourceLinodeVLANRead(d, meta)
}

func resourceLinodeVLANRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to parse vlan ID %s as int: %s", d.Id(), err)
	}

	vlan, err := client.GetVLAN(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to get vlan: %s", err)
	}

	d.Set("description", vlan.Description)
	d.Set("region", vlan.Region)
	d.Set("cidr_block", vlan.CIDRBlock)
	d.Set("attached_linodes", flattenLinodeVLANLinodes(vlan.Linodes))
	return nil
}

func resourceLinodeVLANUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to parse vlan ID %s as int: %s", d.Id(), err)
	}

	vlan, err := client.GetVLAN(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to get vlan %d: %s", id, err)
	}

	var toAttach, toDetach []int
	attached := make(map[int]struct{})

	for _, linode := range vlan.Linodes {
		attached[linode.ID] = struct{}{}
	}

	for _, linodeID := range expandIntSet(d.Get("linodes").(*schema.Set)) {
		if _, ok := attached[linodeID]; !ok {
			toAttach = append(toAttach, linodeID)
		} else {
			delete(attached, linodeID)
		}
	}
	for linodeID := range attached {
		toDetach = append(toDetach, linodeID)
	}

	if len(toDetach) != 0 {
		if _, err := client.DetachVLAN(context.Background(), id, linodego.VLANDetachOptions{
			Linodes: toDetach,
		}); err != nil {
			return fmt.Errorf("failed to detach linodes %v from vlan %d: %s", toDetach, id, err)
		}
	}

	if len(toAttach) != 0 {
		if _, err := client.AttachVLAN(context.Background(), id, linodego.VLANAttachOptions{
			Linodes: toAttach,
		}); err != nil {
			return fmt.Errorf("failed to attach linodes %v from vlan %d: %s", toDetach, id, err)
		}
	}
	return resourceLinodeVLANRead(d, meta)
}

func resourceLinodeVLANDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("failed to parse vlan ID %s as int: %s", d.Id(), err)
	}

	vlan, err := client.GetVLAN(context.Background(), id)
	if err != nil {
		return fmt.Errorf("failed to get vlan %d: %s", id, err)
	}

	linodes := make([]int, len(vlan.Linodes))
	for i, linode := range vlan.Linodes {
		linodes[i] = linode.ID
	}

	if len(linodes) != 0 {
		if _, err := client.DetachVLAN(context.Background(), vlan.ID, linodego.VLANDetachOptions{
			Linodes: linodes,
		}); err != nil {
			return fmt.Errorf("failed to detach linodes %v from vlan %d: %s", linodes, id, err)
		}
	}

	if err := client.DeleteVLAN(context.Background(), id); err != nil {
		return fmt.Errorf("failed to delete vlan: %s", err)
	}
	return nil
}

func flattenLinodeVLANLinodes(linodes []linodego.VLANLinode) []map[string]interface{} {
	attachedLinodes := make([]map[string]interface{}, len(linodes))
	for i, linode := range linodes {
		attachedLinodes[i] = map[string]interface{}{
			"id":           linode.ID,
			"mac_address":  linode.MacAddress,
			"ipv4_address": linode.IPv4Address,
		}
	}
	return attachedLinodes
}
