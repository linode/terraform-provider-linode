package domain

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Domain ID %s as int: %s", d.Id(), err)
	}

	domain, err := client.GetDomain(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode Domain ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Domain: %s", err)
	}

	d.Set("domain", domain.Domain)
	d.Set("type", domain.Type)
	d.Set("group", domain.Group)
	d.Set("status", domain.Status)
	d.Set("description", domain.Description)
	d.Set("master_ips", domain.MasterIPs)
	if len(domain.AXfrIPs) > 0 {
		d.Set("axfr_ips", domain.AXfrIPs)
	}
	d.Set("ttl_sec", domain.TTLSec)
	d.Set("retry_sec", domain.RetrySec)
	d.Set("expire_sec", domain.ExpireSec)
	d.Set("refresh_sec", domain.RefreshSec)
	d.Set("soa_email", domain.SOAEmail)
	d.Set("tags", domain.Tags)

	return nil
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	createOpts := linodego.DomainCreateOptions{
		Domain:      d.Get("domain").(string),
		Type:        linodego.DomainType(d.Get("type").(string)),
		Group:       d.Get("group").(string),
		Description: d.Get("description").(string),
		SOAEmail:    d.Get("soa_email").(string),
		RetrySec:    d.Get("retry_sec").(int),
		ExpireSec:   d.Get("expire_sec").(int),
		RefreshSec:  d.Get("refresh_sec").(int),
		TTLSec:      d.Get("ttl_sec").(int),
	}

	if tagsRaw, tagsOk := d.GetOk("tags"); tagsOk {
		for _, tag := range tagsRaw.(*schema.Set).List() {
			createOpts.Tags = append(createOpts.Tags, tag.(string))
		}
	}

	if v, ok := d.GetOk("master_ips"); ok {
		v := v.(*schema.Set).List()

		createOpts.MasterIPs = make([]string, len(v))
		for i, ip := range v {
			createOpts.MasterIPs[i] = ip.(string)
		}
	}

	if v, ok := d.GetOk("axfr_ips"); ok {
		v := v.(*schema.Set).List()

		createOpts.AXfrIPs = make([]string, len(v))
		for i, ip := range v {
			createOpts.AXfrIPs[i] = ip.(string)
		}
	}

	domain, err := client.CreateDomain(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Domain: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", domain.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Domain id %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.DomainUpdateOptions{
		Domain:      d.Get("domain").(string),
		Status:      linodego.DomainStatus(d.Get("status").(string)),
		Group:       d.Get("group").(string),
		Description: d.Get("description").(string),
		SOAEmail:    d.Get("soa_email").(string),
		RetrySec:    d.Get("retry_sec").(int),
		ExpireSec:   d.Get("expire_sec").(int),
		RefreshSec:  d.Get("refresh_sec").(int),
		TTLSec:      d.Get("ttl_sec").(int),
	}

	if d.HasChange("master_ips") {
		v := d.Get("master_ips").(*schema.Set).List()

		updateOpts.MasterIPs = make([]string, len(v))
		for i, ip := range v {
			updateOpts.MasterIPs[i] = ip.(string)
		}
	}

	if d.HasChange("axfr_ips") {
		v := d.Get("axfr_ips").(*schema.Set).List()

		updateOpts.AXfrIPs = make([]string, len(v))
		for i, ip := range v {
			updateOpts.AXfrIPs[i] = ip.(string)
		}
	}

	if d.HasChange("tags") {
		tags := []string{}
		for _, tag := range d.Get("tags").(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}

		updateOpts.Tags = tags
	}

	_, err = client.UpdateDomain(ctx, id, updateOpts)
	if err != nil {
		return diag.Errorf("Error updating Linode Domain %d: %s", id, err)
	}
	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode Domain id %s as int", d.Id())
	}
	err = client.DeleteDomain(ctx, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode Domain %d: %s", id, err)
	}
	d.SetId("")

	return nil
}
