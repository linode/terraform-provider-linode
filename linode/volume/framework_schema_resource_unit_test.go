//go:build unit

package volume

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/stretchr/testify/require"
)

// test that guards the schema for the encryption attribute.
func TestEncryptionAttribute_HasDefaultAndRequiresReplace(t *testing.T) {
	t.Helper()

	attrRaw, ok := frameworkResourceSchema.Attributes["encryption"]
	require.True(t, ok, "encryption attribute must exist in schema")

	attr, ok := attrRaw.(schema.StringAttribute)
	require.True(t, ok, "encryption must be a StringAttribute")

	// Should be Optional + Computed with no schema default; provider preserves state when omitted.
	require.True(t, attr.Optional, "encryption should be Optional")
	require.True(t, attr.Computed, "encryption should be Computed")
	require.Nil(t, attr.Default, "encryption should not have a schema default")

	// Must preserve state when config omits the field.
	expectedUseStateType := reflect.TypeOf(stringplanmodifier.UseStateForUnknown())
	foundUseState := false
	for _, pm := range attr.PlanModifiers {
		if reflect.TypeOf(pm) == expectedUseStateType {
			foundUseState = true
			break
		}
	}
	require.True(t, foundUseState, "encryption should have UseStateForUnknown plan modifier")

	// Must require replacement when changed.
	expectedReplaceType := reflect.TypeOf(stringplanmodifier.RequiresReplace())
	foundReplace := false
	for _, pm := range attr.PlanModifiers {
		if reflect.TypeOf(pm) == expectedReplaceType {
			foundReplace = true
			break
		}
	}
	require.True(t, foundReplace, "encryption should have a RequiresReplace plan modifier")

	// Should show default "enabled" on create when omitted.
	expectedDefaultOnCreateType := reflect.TypeOf(DefaultEnabledOnCreate())
	foundDefaultOnCreate := false
	for _, pm := range attr.PlanModifiers {
		if reflect.TypeOf(pm) == expectedDefaultOnCreateType {
			foundDefaultOnCreate = true
			break
		}
	}
	require.True(t, foundDefaultOnCreate, "encryption should have DefaultEnabledOnCreate plan modifier")

	// Should have validators (e.g., OneOf("enabled","disabled")). We don't assert exact type, just presence.
	require.NotEmpty(t, attr.Validators, "encryption should have validators (e.g., OneOf)")
}
