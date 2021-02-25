package linode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
)

func resourceLinodeUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLinodeUserCreate,
		ReadContext:   resourceLinodeUserRead,
		UpdateContext: resourceLinodeUserUpdate,
		DeleteContext: resourceLinodeUserDelete,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email of the user.",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The username of the user.",
			},
			"restricted": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, the user must be explicitly granted access to platform actions and entities.",
			},
			"ssh_keys": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "SSH keys to add to the user profile.",
			},
			"tfa_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If the User has Two Factor Authentication (TFA) enabled.",
			},
		},
	}
}

func resourceLinodeUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client
	createOpts := linodego.UserCreateOptions{
		Email:      d.Get("email").(string),
		Username:   d.Get("username").(string),
		Restricted: d.Get("restricted").(bool),
	}

	user, err := client.CreateUser(ctx, createOpts)
	if err != nil {
		return diag.Errorf("failed to create user: %s", err)
	}
	d.SetId(user.Username)

	return resourceLinodeUserRead(ctx, d, meta)
}

func resourceLinodeUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	username := d.Get("username").(string)
	user, err := client.GetUser(ctx, username)
	if err != nil {
		return diag.Errorf("failed to get user (%s): %s", username, err)
	}

	d.Set("username", username)
	d.Set("email", user.Email)
	d.Set("restricted", user.Restricted)
	d.Set("ssh_keys", user.SSHKeys)
	d.Set("tfa_enabled", user.TFAEnabled)
	return nil
}

func resourceLinodeUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	id := d.Id()
	username := d.Get("username").(string)
	restricted := d.Get("restricted").(bool)
	if _, err := client.UpdateUser(ctx, id, linodego.UserUpdateOptions{
		Username:   username,
		Restricted: &restricted,
	}); err != nil {
		return diag.Errorf("failed to update user (%s): %s", id, err)
	}

	d.SetId(username)
	return resourceLinodeUserRead(ctx, d, meta)
}

func resourceLinodeUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderMeta).Client

	username := d.Get("username").(string)
	if err := client.DeleteUser(ctx, username); err != nil {
		return diag.Errorf("failed to delete user (%s): %s", username, err)
	}
	return nil
}
