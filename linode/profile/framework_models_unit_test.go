//go:build unit

package profile

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseProfile(t *testing.T) {
	// Prepare a mock Linode Profile
	mockProfile := &linodego.Profile{
		Email:              "test@example.com",
		Timezone:           "UTC",
		EmailNotifications: true,
		Username:           "testuser",
		IPWhitelistEnabled: false,
		LishAuthMethod:     linodego.AuthMethodPasswordKeys,
		AuthorizedKeys:     []string{"ssh-rsa AAAAB3NzaC1yc2E...", "ssh-rsa AAAAB3NzaCNUMBER2..."},
		TwoFactorAuth:      true,
		Restricted:         false,
		Referrals: linodego.ProfileReferrals{
			Total:     10,
			Completed: 5,
			Pending:   5,
			Credit:    50,
			Code:      "SECRETREFCODEOMG",
			URL:       "https://example.com/referral",
		},
	}

	profileModel := &DataSourceModel{}

	profileModel.parseProfile(context.Background(), mockProfile)

	assert.Equal(t, types.StringValue("test@example.com"), profileModel.Email)
	assert.Equal(t, types.StringValue("UTC"), profileModel.Timezone)
	assert.Equal(t, types.BoolValue(true), profileModel.EmailNotifications)
	assert.Equal(t, types.StringValue("testuser"), profileModel.Username)
	assert.Equal(t, types.StringValue("password_keys"), profileModel.LishAuthMethod)

	for _, authKey := range mockProfile.AuthorizedKeys {
		assert.Contains(t, profileModel.AuthorizedKeys.String(), authKey)
	}

	assert.Equal(t, types.StringValue("password_keys"), profileModel.LishAuthMethod)
	assert.Equal(t, types.BoolValue(false), profileModel.IPWhitelistEnabled)
	assert.Equal(t, types.BoolValue(true), profileModel.TwoFactorAuth)
	assert.Equal(t, types.BoolValue(false), profileModel.Restricted)

	assert.Equal(t, "SECRETREFCODEOMG", profileModel.Referrals.Elements()[0].(types.Object).Attributes()["code"].(types.String).ValueString())
	assert.Equal(t, 50.000000, profileModel.Referrals.Elements()[0].(types.Object).Attributes()["credit"].(types.Float64).ValueFloat64())
	assert.Equal(t, int64(10), profileModel.Referrals.Elements()[0].(types.Object).Attributes()["total"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(5), profileModel.Referrals.Elements()[0].(types.Object).Attributes()["pending"].(types.Int64).ValueInt64())
	assert.Equal(t, int64(5), profileModel.Referrals.Elements()[0].(types.Object).Attributes()["completed"].(types.Int64).ValueInt64())
	assert.Equal(t, "https://example.com/referral", profileModel.Referrals.Elements()[0].(types.Object).Attributes()["url"].(types.String).ValueString())
}
