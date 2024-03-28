//go:build unit

package childaccounts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/linodego"
)

func TestParseAccounts(t *testing.T) {
	model := &ChildAccountFilterModel{}

	// Sample input login data
	account1 := linodego.ChildAccount{
		EUUID:     "12345",
		FirstName: "foo",
		LastName:  "bar",
	}

	account2 := linodego.ChildAccount{
		EUUID:     "54321",
		FirstName: "bar",
		LastName:  "foo",
	}

	model.parseAccounts([]linodego.ChildAccount{account1, account2})

	if len(model.ChildAccounts) != 2 {
		t.Errorf("Expected %d logins, but got %d", 2, len(model.ChildAccounts))
	}

	// Check if the fields of the first login in the model have been populated correctly
	if !model.ChildAccounts[0].EUUID.Equal(types.StringValue("12345")) {
		t.Errorf("Expected ID to be 12345, but got %s", model.ChildAccounts[0].EUUID.ValueString())
	}
}
