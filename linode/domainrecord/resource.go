package domainrecord

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
		CreateContext: createResource,
		ReadContext:   readResource,
		UpdateContext: updateResource,
		DeleteContext: deleteResource,
		Importer: &schema.ResourceImporter{
			StateContext: importResource,
		},
	}
}

func importResource(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if !strings.Contains(d.Id(), ",") {
		return nil, fmt.Errorf("failed to parse argument: %v", d.Id())
	}

	s := strings.Split(d.Id(), ",")

	if len(s) != 2 {
		return nil, fmt.Errorf("invalid number of arguments: %v", len(s))
	}
	// Validate that this is an ID by making sure it can be converted into an int
	_, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, fmt.Errorf("invalid domain_record ID: %v", err)
	}

	domainID, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, fmt.Errorf("invalid domain ID: %v", err)
	}

	d.SetId(s[1])
	d.Set("domain_id", domainID)

	if err := readResource(ctx, d, meta); err != nil {
		return nil, fmt.Errorf("unable to import %v as domain_record: %v", d.Id(), err)
	}

	results := make([]*schema.ResourceData, 0)
	results = append(results, d)

	return results, nil
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode DomainRecord ID %s as int: %s", d.Id(), err)
	}
	domainID := d.Get("domain_id").(int)
	record, err := client.GetDomainRecord(ctx, domainID, id)
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing Linode Domain Record ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode DomainRecord: %s", err)
	}

	d.Set("name", record.Name)
	d.Set("port", record.Port)
	d.Set("priority", record.Priority)
	d.Set("protocol", record.Protocol)
	d.Set("service", record.Service)
	d.Set("tag", record.Tag)
	d.Set("target", record.Target)
	d.Set("ttl_sec", record.TTLSec)
	d.Set("record_type", record.Type)
	d.Set("weight", record.Weight)

	return nil
}

func resourceDataStringOrNil(d *schema.ResourceData, name string) *string {
	if val, ok := d.GetOkExists(name); ok {
		i := val.(string)
		if len(i) == 0 {
			return nil
		}
		return &i
	}
	return nil
}

func resourceDataIntOrNil(d *schema.ResourceData, name string) *int {
	if val, ok := d.GetOkExists(name); ok {
		i := val.(int)
		return &i
	}
	return nil
}

func domainRecordFromResourceData(d *schema.ResourceData) *linodego.DomainRecord {
	return &linodego.DomainRecord{
		Name:     d.Get("name").(string),
		Type:     linodego.DomainRecordType(d.Get("record_type").(string)),
		Target:   d.Get("target").(string),
		Priority: d.Get("priority").(int),
		Weight:   d.Get("weight").(int),
		Port:     d.Get("port").(int),
		Service:  resourceDataStringOrNil(d, "service"),
		Protocol: resourceDataStringOrNil(d, "protocol"),
		TTLSec:   d.Get("ttl_sec").(int),
		Tag:      resourceDataStringOrNil(d, "tag"),
	}
}

func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	domainID := d.Get("domain_id").(int)
	rec := domainRecordFromResourceData(d)

	createOpts := linodego.DomainRecordCreateOptions{
		Type:     rec.Type,
		Name:     rec.Name,
		Target:   rec.Target,
		Priority: resourceDataIntOrNil(d, "priority"),
		Weight:   resourceDataIntOrNil(d, "weight"),
		Port:     resourceDataIntOrNil(d, "port"),
		Service:  resourceDataStringOrNil(d, "service"),
		Protocol: resourceDataStringOrNil(d, "protocol"),
		TTLSec:   rec.TTLSec,
		Tag:      resourceDataStringOrNil(d, "tag"),
	}

	domainRecord, err := client.CreateDomainRecord(ctx, domainID, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode DomainRecord: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", domainRecord.ID))

	return readResource(ctx, d, meta)
}

func updateResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	domainID := d.Get("domain_id").(int)
	rec := domainRecordFromResourceData(d)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode DomainRecord id %s as int: %s", d.Id(), err)
	}
	updateOpts := linodego.DomainRecordUpdateOptions{
		Type:     rec.Type,
		Name:     rec.Name,
		Target:   rec.Target,
		Priority: resourceDataIntOrNil(d, "priority"),
		Weight:   resourceDataIntOrNil(d, "weight"),
		Port:     resourceDataIntOrNil(d, "port"),
		Service:  resourceDataStringOrNil(d, "service"),
		Protocol: resourceDataStringOrNil(d, "protocol"),
		TTLSec:   rec.TTLSec,
		Tag:      resourceDataStringOrNil(d, "tag"),
	}

	_, err = client.UpdateDomainRecord(ctx, domainID, id, updateOpts)
	if err != nil {
		return diag.Errorf("Error updating Domain Record: %s", err)
	}

	return readResource(ctx, d, meta)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client
	domainID := d.Get("domain_id").(int)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error parsing Linode DomainRecord id %s as int", d.Id())
	}
	err = client.DeleteDomainRecord(ctx, domainID, id)
	if err != nil {
		return diag.Errorf("Error deleting Linode DomainRecord %d: %s", id, err)
	}
	d.SetId("")

	return nil
}

func domainRecordTargetSuppressor(k, provisioned, declared string, d *schema.ResourceData) bool {
	return len(strings.Split(declared, ".")) == 1 &&
		strings.Contains(provisioned, declared)
}
