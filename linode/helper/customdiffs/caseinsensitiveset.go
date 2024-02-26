package customdiffs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// CaseInsensitiveSet allows us to ignore case diffs on case-insensitive
// set values (e.g. tags).
//
// This function could not be implemented as a DiffSuppressFunc because DiffSuppressFuncs
// are run per-entry rather than on the set as a whole.
//
// NOTE: The target field must be marked as computed.
func CaseInsensitiveSet(field string) schema.CustomizeDiffFunc {
	return func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
		if !diff.HasChange(field) {
			return nil
		}

		oldEntries, newEntries := diff.GetChange(field)
		oldEntriesSet, newEntriesSet := oldEntries.(*schema.Set), newEntries.(*schema.Set)

		// Apply the updated plan
		return diff.SetNew(field, computeCaseInsensitivePlannedSet(oldEntriesSet, newEntriesSet))
	}
}

// computeCaseInsensitivePlannedSet computes a new set to plan that ignores
// any case changes made in the configuration.
func computeCaseInsensitivePlannedSet(oldSet, newSet *schema.Set) *schema.Set {
	result := schema.NewSet(newSet.F, []any{})

	// Map all lowered entries to their original case
	oldEntriesMap := make(map[string]string)
	for _, oldTag := range oldSet.List() {
		oldTag := oldTag.(string)
		oldEntriesMap[strings.ToLower(oldTag)] = oldTag
	}

	// Check if there is a corresponding old entry for the lowered
	// version of each new entry
	for _, newTag := range newSet.List() {
		newTag := newTag.(string)

		oldTagWithCase, ok := oldEntriesMap[strings.ToLower(newTag)]
		if !ok {
			result.Add(newTag)
			continue
		}

		// If we found a match, update the entry in the
		// plan to match the old case
		result.Add(oldTagWithCase)
	}

	return result
}
