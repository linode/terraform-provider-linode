// +build ignore

package linode

import (
	"fmt"
	"strconv"

	"github.com/chiefy/linodego"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLinodeTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceLinodeTemplateCreate,
		Read:   resourceLinodeTemplateRead,
		Update: resourceLinodeTemplateUpdate,
		Delete: resourceLinodeTemplateDelete,
		Exists: resourceLinodeTemplateExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"label": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The label of the Linode Template.",
				Optional:    true,
			},
			"status": &schema.Schema{
				Type:        schema.TypeInt,
				Description: "The status of the template, indicating the current readiness state.",
				Computed:    true,
			},
		},
	}
}

func resourceLinodeTemplateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return false, fmt.Errorf("Failed to parse Linode Template ID %s as int because %s", d.Id(), err)
	}

	_, err = client.GetTemplate(int(id))
	if err != nil {
		return false, fmt.Errorf("Failed to get Linode Template ID %s because %s", d.Id(), err)
	}
	return true, nil
}

func resourceLinodeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Template ID %s as int because %s", d.Id(), err)
	}

	template, err := client.GetTemplate(int(id))

	if err != nil {
		return fmt.Errorf("Failed to find the specified Linode Template because %s", err)
	}

	d.Set("label", template.Label)
	d.Set("status", template.Status)

	return nil
}

func resourceLinodeTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client, ok := meta.(linodego.Client)
	if !ok {
		return fmt.Errorf("Invalid Client when creating Linode Template")
	}
	d.Partial(true)

	createOpts := linodego.TemplateCreateOptions{
		Label:  d.Get("label").(string),
	}
	template, err := client.CreateTemplate(&createOpts)
	if err != nil {
		return fmt.Errorf("Failed to create a Linode Template in because %s", err)
	}
	d.SetId(fmt.Sprintf("%d", template.ID))
	d.Set("label", template.Label)
	d.SetPartial("label")

	return resourceLinodeTemplateRead(d, meta)
}

func resourceLinodeTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	d.Partial(true)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Template id %s as an int because %s", d.Id(), err)
	}

	template, err := client.GetTemplate(int(id))
	if err != nil {
		return fmt.Errorf("Failed to fetch data about the current linode because %s", err)
	}

	if d.HasChange("label") {
		if template, err = client.RenameTemplate(template.ID, d.Get("label").(string)); err != nil {
			return err
		}
		d.Set("label", template.Label)
		d.SetPartial("label")
	}

	return nil // resourceLinodeTemplateRead(d, meta)
}

func resourceLinodeTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse Linode Template id %s as int", d.Id())
	}
	err = client.DeleteTemplate(int(id))
	if err != nil {
		return fmt.Errorf("Failed to delete Linode Template %d because %s", id, err)
	}
	return nil
}
