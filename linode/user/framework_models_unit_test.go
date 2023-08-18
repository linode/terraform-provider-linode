package user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"testing"
	"time"

	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
)

func TestParseUser(t *testing.T) {
	// Sample data
	phoneNumber := "+5555555555"
	userData := linodego.User{
		Username:            "example_user",
		Email:               "example_user@linode.com",
		Restricted:          true,
		TFAEnabled:          true,
		SSHKeys:             []string{"home-pc", "laptop"},
		PasswordCreated:     &time.Time{},
		VerifiedPhoneNumber: &phoneNumber,
	}

	dataModel := &DataSourceModel{}

	diags := dataModel.ParseUser(context.Background(), &userData)

	assert.False(t, diags.HasError())

	assert.Equal(t, types.StringValue("example_user"), dataModel.Username)
	assert.Equal(t, types.StringValue("example_user@linode.com"), dataModel.Email)
	assert.Equal(t, types.BoolValue(true), dataModel.Restricted)
	assert.Equal(t, types.BoolValue(true), dataModel.TFAEnabled)
	for _, key := range userData.SSHKeys {
		assert.Contains(t, dataModel.SSHKeys.String(), key)
	}
	assert.Contains(t, dataModel.VerifiedPhoneNumber.String(), "+5555555555")
}

func TestParseUserGrants(t *testing.T) {
	permissionLevel := linodego.AccessLevelReadOnly
	userGrantsData := linodego.UserGrants{
		Database: []linodego.GrantedEntity{
			{
				ID:          123,
				Label:       "example-database",
				Permissions: "read_only",
			},
		},
		Domain: []linodego.GrantedEntity{
			{
				ID:          456,
				Label:       "example-domain",
				Permissions: "read_write",
			},
		},
		Firewall: []linodego.GrantedEntity{
			{
				ID:          789,
				Label:       "example-firewall",
				Permissions: "read_only",
			},
		},
		Global: linodego.GlobalUserGrants{
			AccountAccess:        &permissionLevel,
			AddDatabases:         true,
			AddDomains:           true,
			AddFirewalls:         true,
			AddImages:            true,
			AddLinodes:           true,
			AddLongview:          true,
			AddNodeBalancers:     true,
			AddStackScripts:      true,
			AddVolumes:           true,
			CancelAccount:        false,
			LongviewSubscription: true,
		},
		Image: []linodego.GrantedEntity{
			{
				ID:          101,
				Label:       "example-image",
				Permissions: "read_write",
			},
		},
		Linode: []linodego.GrantedEntity{
			{
				ID:          102,
				Label:       "example-linode",
				Permissions: "read_write",
			},
		},
		Longview: []linodego.GrantedEntity{
			{
				ID:          103,
				Label:       "example-longview",
				Permissions: "read_only",
			},
		},
		NodeBalancer: []linodego.GrantedEntity{
			{
				ID:          104,
				Label:       "example-nodebalancer",
				Permissions: "read_only",
			},
		},
		StackScript: []linodego.GrantedEntity{
			{
				ID:          105,
				Label:       "example-stackscript",
				Permissions: "read_write",
			},
		},
		Volume: []linodego.GrantedEntity{
			{
				ID:          106,
				Label:       "example-volume",
				Permissions: "read_only",
			},
		},
	}

	dataModel := &DataSourceModel{}

	diags := dataModel.ParseUserGrants(context.Background(), &userGrantsData)

	assert.False(t, diags.HasError())
	assert.Contains(t, dataModel.DatabaseGrant.String(), strconv.Itoa(userGrantsData.Database[0].ID))
	assert.Contains(t, dataModel.FirewallGrant.String(), strconv.Itoa(userGrantsData.Firewall[0].ID))
	assert.Contains(t, dataModel.ImageGrant.String(), strconv.Itoa(userGrantsData.Image[0].ID))
	assert.Contains(t, dataModel.LinodeGrant.String(), strconv.Itoa(userGrantsData.Linode[0].ID))
	assert.Contains(t, dataModel.LongviewGrant.String(), strconv.Itoa(userGrantsData.Longview[0].ID))
	assert.Contains(t, dataModel.NodebalancerGrant.String(), strconv.Itoa(userGrantsData.NodeBalancer[0].ID))
	assert.Contains(t, dataModel.StackscriptGrant.String(), strconv.Itoa(userGrantsData.StackScript[0].ID))
	assert.Contains(t, dataModel.VolumeGrant.String(), strconv.Itoa(userGrantsData.Volume[0].ID))

	assert.Contains(t, dataModel.GlobalGrants.String(), "\"account_access\":\"read_only\"")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_databases\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_domains\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_firewalls\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_images\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_linodes\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_longview\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_nodebalancers\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_stackscripts\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"add_volumes\":true")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"cancel_account\":false")
	assert.Contains(t, dataModel.GlobalGrants.String(), "\"longview_subscription\":true")

}
