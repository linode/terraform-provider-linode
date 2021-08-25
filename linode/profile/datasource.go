package profile

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Schema:      dataSourceSchema,
		ReadContext: readDataSource,
	}
}

func readDataSource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*helper.ProviderMeta).Client

	profile, err := client.GetProfile(ctx)
	if err != nil {
		return diag.Errorf("Error getting profile: %s", err)
	}

	d.SetId(fmt.Sprintf("%d", profile.UID))
	d.Set("referrals", flattenProfileReferrals(profile.Referrals))
	d.Set("email", profile.Email)
	d.Set("timezone", profile.Timezone)
	d.Set("email_notifications", profile.EmailNotifications)
	d.Set("username", profile.Username)
	d.Set("ip_whitelist_enabled", profile.IPWhitelistEnabled)
	d.Set("lish_auth_method", profile.LishAuthMethod)
	d.Set("authorized_keys", profile.AuthorizedKeys)
	d.Set("two_factor_auth", profile.TwoFactorAuth)
	d.Set("restricted", profile.Restricted)

	return nil
}

type flattenedProfileReferrals map[string]interface{}

func flattenProfileReferrals(referrals linodego.ProfileReferrals) []flattenedProfileReferrals {
	return []flattenedProfileReferrals{{
		"code":      referrals.Code,
		"url":       referrals.URL,
		"total":     referrals.Total,
		"completed": referrals.Completed,
		"pending":   referrals.Pending,
		"credit":    referrals.Credit,
	}}
}
