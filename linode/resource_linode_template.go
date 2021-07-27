// +build ignore

/**
 * Using this template:
 * - Copy resource_linode_template.go and resource_linode_template_test.go
 *   - Remove "// +build ignore"
 *   - Replace "Template" with Linode Resource Name
 *   - Replace "template" with Linode resource name
 * - Add Resource to provider.go
 */
package linode

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeTemplateCreate,
		ReadContext:   resourceLinodeTemplateRead,
		UpdateContext: resourceLinodeTemplateUpdate,
		DeleteContext: resourceLinodeTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The label of the Linode Template.",
				Optional:    true,
			},
			"status": {
				Type:        schema.TypeInt,
				Description: "The status of the template, indicating the current readiness state.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Template ID %s as int: %s", d.Id(), err)
	}

	template, err := client.GetTemplate(ctx, int(id))

	if err != nil {
		return diag.Errorf("Error finding the specified Linode Template: %s", err)
	}

	d.Set("label", template.Label)
	d.Set("status", template.Status)

	return nil
}

func resourceLinodeTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*ProviderMeta).Client
	if !ok {
		return diag.Errorf("Invalid Client when creating Linode Template")
	}

	createOpts := linodego.TemplateCreateOptions{
		Label: d.Get("label").(string),
	}
	template, err := client.CreateTemplate(ctx, createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Template: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", template.ID))
	d.Set("label", template.Label)

	return resourceLinodeTemplateRead(ctx, d, meta)
}

func resourceLinodeTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Template id %s as int: %s", d.Id(), err)
	}

	template, err := client.GetTemplate(int(id))
	if err != nil {
		return diag.Errorf("Error fetching data about the current Linode Template: %s", err)
	}

	if d.HasChange("label") {
		if template, err = client.RenameTemplate(ctx, template.ID, d.Get("label").(string)); err != nil {
			return diag.FromErr(err)
		}
		d.Set("label", template.Label)
	}

	return resourceLinodeTemplateRead(ctx, d, meta)
}

func resourceLinodeTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Template id %s as int", d.Id())
	}
	err = client.DeleteTemplate(ctx, int(id))
	if err != nil {
		return diag.Errorf("Error deleting Linode Template %d: %s", id, err)
	}
	return nil
}
