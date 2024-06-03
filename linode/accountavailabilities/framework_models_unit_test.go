//go:build unit

package accountavailabilities

import (
	"context"
	"reflect"
	"testing"

	"github.com/linode/linodego"
)

func TestParseAvailabilities(t *testing.T) {
	model := &AccountAvailabilityFilterModel{}

	// Sample input avail data
	avail1 := linodego.AccountAvailability{
		Region: "us-mia",
		Unavailable: []string{
			"something",
			"something else",
		},
		Available: []string{
			"something available",
		},
	}

	avail2 := linodego.AccountAvailability{
		Region: "us-iad",
		Unavailable: []string{
			"another thing",
		},
		Available: []string{
			"something else",
		},
	}

	model.parseAvailabilities(context.Background(), []linodego.AccountAvailability{avail1, avail2})

	if len(model.Availabilities) != 2 {
		t.Errorf("Expected %d avail entries, but got %d", 2, len(model.Availabilities))
	}

	// Check if the fields of the first entry in the model have been populated correctly
	if model.Availabilities[0].Region.ValueString() != "us-mia" {
		t.Errorf(
			"Expected region to be %s, but got %s",
			"us-mia", model.Availabilities[0].Region.ValueString(),
		)
	}

	var unavailList []string
	model.Availabilities[0].Unavailable.ElementsAs(context.Background(), &unavailList, false)

	if !reflect.DeepEqual(
		unavailList,
		[]string{"something", "something else"},
	) {
		t.Error(
			"mismatch when comparing unavailable list",
		)
	}

	var availList []string
	model.Availabilities[0].Available.ElementsAs(context.Background(), &availList, false)
	if !reflect.DeepEqual(
		availList,
		[]string{"something available"},
	) {
		t.Error(
			"mismatch when comparing available list",
		)
	}
}
