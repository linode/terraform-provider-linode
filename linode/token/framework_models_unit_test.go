//go:build unit

package token

import (
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/helper/customtypes"
	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	expiryDate := time.Date(2050, time.August, 17, 12, 0, 0, 0, time.UTC)

	sampleToken := linodego.Token{
		ID:      123,
		Scopes:  "*",
		Label:   "Test Token",
		Token:   "test-token-value",
		Created: &createdTime,
		Expiry:  &expiryDate,
	}

	model := &ResourceModel{}

	model.parseToken(&sampleToken, false)

	assert.Equal(t, types.StringValue(sampleToken.Label), model.Label)
	assert.Equal(t, customtypes.LinodeScopesStringValue{StringValue: types.StringValue(sampleToken.Scopes)}, model.Scopes)
	assert.Equal(t, types.StringValue(sampleToken.Token), model.Token)
	assert.Equal(t, types.StringValue(strconv.Itoa(sampleToken.ID)), model.ID)
}

func TestParseTokenRefreshTrue(t *testing.T) {
	createdTime := time.Date(2023, time.August, 17, 12, 0, 0, 0, time.UTC)
	expiryDate := time.Date(2050, time.August, 17, 12, 0, 0, 0, time.UTC)

	sampleToken := linodego.Token{
		ID:      456,
		Scopes:  "linodes",
		Label:   "Another Token",
		Token:   "another-token-value",
		Created: &createdTime,
		Expiry:  &expiryDate, // Set to nil for testing purposes
	}

	rm := &ResourceModel{}

	rm.parseToken(&sampleToken, true)

	assert.Empty(t, rm.Token)
}
