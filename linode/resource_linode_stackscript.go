package linode

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeStackscript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeStackscriptCreateContext,
		ReadContext:   resourceLinodeStackscriptReadContext,
		UpdateContext: resourceLinodeStackscriptUpdateContext,
		DeleteContext: resourceLinodeStackscriptDeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"label": {
				Type:        schema.TypeString,
				Description: "The StackScript's label is for display purposes only.",
				Required:    true,
			},
			"script": {
				Type:        schema.TypeString,
				Description: "The script to execute when provisioning a new Linode with this StackScript.",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "A description for the StackScript.",
				Required:    true,
			},
			"rev_note": {
				Type:        schema.TypeString,
				Description: "This field allows you to add notes for the set of revisions made to this StackScript.",
				Optional:    true,
			},
			"is_public": {
				Type:        schema.TypeBool,
				Description: "This determines whether other users can use your StackScript. Once a StackScript is made public, it cannot be made private.",
				Default:     false,
				Optional:    true,
				ForceNew:    true,
			},
			"images": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An array of Image IDs representing the Images that this StackScript is compatible for deploying with.",
				Required:    true,
			},

			"deployments_active": {
				Type:        schema.TypeInt,
				Description: "Count of currently active, deployed Linodes created from this StackScript.",
				Computed:    true,
			},
			"user_gravatar_id": {
				Type:        schema.TypeString,
				Description: "The Gravatar ID for the User who created the StackScript.",
				Computed:    true,
			},
			"deployments_total": {
				Type:        schema.TypeInt,
				Description: "The total number of times this StackScript has been deployed.",
				Computed:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "The User who created the StackScript.",
				Computed:    true,
			},
			"created": {
				Type:        schema.TypeString,
				Description: "The date this StackScript was created.",
				Computed:    true,
			},
			"updated": {
				Type:        schema.TypeString,
				Description: "The date this StackScript was updated.",
				Computed:    true,
			},
			"user_defined_fields": {
				Description: "This is a list of fields defined with a special syntax inside this StackScript that allow for supplying customized parameters during deployment.",
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				ConfigMode:  schema.SchemaConfigModeAttr,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"example": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"one_of": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"many_of": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceLinodeStackscriptReadContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript ID %s as int: %s", d.Id(), err)
	}

	stackscript, err := client.GetStackscript(context.Background(), int(id))

	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			log.Printf("[WARN] removing StackScript ID %q from state because it no longer exists", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error finding the specified Linode Stackscript: %s", err)
	}

	d.Set("label", stackscript.Label)
	d.Set("script", stackscript.Script)
	d.Set("description", stackscript.Description)
	d.Set("is_public", stackscript.IsPublic)
	d.Set("images", stackscript.Images)
	d.Set("rev_note", stackscript.RevNote)

	// Computed
	d.Set("deployments_active", stackscript.DeploymentsActive)
	d.Set("deployments_total", stackscript.DeploymentsTotal)
	d.Set("username", stackscript.Username)
	d.Set("user_gravatar_id", stackscript.UserGravatarID)
	d.Set("created", stackscript.Created.String())
	d.Set("updated", stackscript.Updated.String())
	setStackScriptUserDefinedFields(d, stackscript)
	return nil
}

func resourceLinodeStackscriptCreateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(linodego.Client)
	if !ok {
		return diag.Errorf("Invalid Client when creating Linode Stackscript")
	}

	createOpts := linodego.StackscriptCreateOptions{
		Label:       d.Get("label").(string),
		Script:      d.Get("script").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
		RevNote:     d.Get("rev_note").(string),
	}

	for _, image := range d.Get("images").([]interface{}) {
		createOpts.Images = append(createOpts.Images, image.(string))
	}

	stackscript, err := client.CreateStackscript(context.Background(), createOpts)
	if err != nil {
		return diag.Errorf("Error creating a Linode Stackscript: %s", err)
	}
	d.SetId(fmt.Sprintf("%d", stackscript.ID))

	return resourceLinodeStackscriptReadContext(ctx, d, meta)
}

func resourceLinodeStackscriptUpdateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript id %s as int: %s", d.Id(), err)
	}

	updateOpts := linodego.StackscriptUpdateOptions{
		Label:       d.Get("label").(string),
		Script:      d.Get("script").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
		RevNote:     d.Get("rev_note").(string),
	}

	for _, image := range d.Get("images").([]interface{}) {
		updateOpts.Images = append(updateOpts.Images, image.(string))
	}

	if _, err = client.UpdateStackscript(context.Background(), int(id), updateOpts); err != nil {
		return diag.Errorf("Error updating Linode Stackscript %d: %s", int(id), err)
	}

	return resourceLinodeStackscriptReadContext(ctx, d, meta)
}

func resourceLinodeStackscriptDeleteContext(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(linodego.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error parsing Linode Stackscript id %s as int", d.Id())
	}
	err = client.DeleteStackscript(context.Background(), int(id))
	if err != nil {
		if lerr, ok := err.(*linodego.Error); ok && lerr.Code == 404 {
			return nil
		}
		return diag.Errorf("Error deleting Linode Stackscript %d: %s", id, err)
	}
	return nil
}
