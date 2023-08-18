//go:build unit

package accountlogins

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
	"testing"
	"time"
)

func TestParseLogins(t *testing.T) {
	model := &AccountLoginFilterModel{}

	// Sample input login data
	login1 := linodego.Login{
		ID:         1,
		Datetime:   &time.Time{},
		IP:         "127.0.0.1",
		Restricted: false,
		Username:   "user1",
		Status:     "success",
	}

	login2 := linodego.Login{
		ID:         2,
		Datetime:   &time.Time{},
		IP:         "192.168.1.1",
		Restricted: true,
		Username:   "user2",
		Status:     "failure",
	}

	model.parseLogins([]linodego.Login{login1, login2})

	if len(model.Logins) != 2 {
		t.Errorf("Expected %d logins, but got %d", 2, len(model.Logins))
	}

	// Check if the fields of the first login in the model have been populated correctly
	if model.Logins[0].ID != types.Int64Value(1) {
		t.Errorf("Expected ID to be %d, but got %d", 1, model.Logins[0].ID)
	}
}
