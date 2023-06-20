package domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

type Resource struct {
	client *linodego.Client
}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client := r.client

	var data DomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueInt64()
	if resp.Diagnostics.HasError() {
		return
	}

	domain, err := client.GetDomain(ctx, int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			resp.Diagnostics.AddWarning(
				"Domain No Longer Exists",
				fmt.Sprintf(
					"Removing Domain with ID %v from state because it no longer exists",
					data.ID,
				),
			)
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error finding the specified Domain",
			err.Error(),
		)
		return
	}

	data.parseDomain(domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client := r.client
	var data DomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := linodego.DomainCreateOptions{
		Domain:      data.Domain.ValueString(),
		Type:        linodego.DomainType(data.Type.ValueString()),
		Group:       data.Group.ValueString(),
		Description: data.Description.ValueString(),
		SOAEmail:    data.SOAEmail.ValueString(),
		RetrySec:    int(data.RetrySec.ValueInt64()),
		ExpireSec:   int(data.ExpireSec.ValueInt64()),
		RefreshSec:  int(data.RefreshSec.ValueInt64()),
		TTLSec:      int(data.TTLSec.ValueInt64()),
		MasterIPs:   helper.FrameworkToStringSlice(data.MasterIPs),
		AXfrIPs:     helper.FrameworkToStringSlice(data.AXFRIPs),
		Tags:        helper.FrameworkToStringSlice(data.Tags),
	}

	domain, err := client.CreateDomain(ctx, createOpts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Domain creation error",
			err.Error(),
		)
		return
	}

	data.parseDomain(domain)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
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

	_, err = client.UpdateDomain(ctx, int(id), updateOpts)
	if err != nil {
		return diag.Errorf("Error updating Linode Domain %d: %s", id, err)
	}
	return readResource(ctx, d, meta)
}

func Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data DomainModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id := data.ID.ValueInt64
	if resp.Diagnostics.HasError() {
		return
	}

	client := r.client
	err := client.DeleteDomain(ctx, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); (ok && lerr.Code != 404) || !ok {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Failed to delete domain with id %v", id),
				err.Error(),
			)
		}
		return
	}
}
