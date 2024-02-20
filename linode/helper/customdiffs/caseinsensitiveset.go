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

		// Map all lowered entries to their original case
		oldEntriesMap := make(map[string]string)
		for _, oldTag := range oldEntriesSet.List() {
			oldTag := oldTag.(string)
			oldEntriesMap[strings.ToLower(oldTag)] = oldTag
		}

		// Check if there is a corresponding old entry for the lowered
		// version of each new entry
		for _, newTag := range newEntriesSet.List() {
			newTag := newTag.(string)

			oldTagWithCase, ok := oldEntriesMap[strings.ToLower(newTag)]
			if !ok {
				continue
			}

			// If we found a match, update the entry in the
			// plan to match the old case
			newEntriesSet.Remove(newTag)
			newEntriesSet.Add(oldTagWithCase)
		}

		// Apply the updated plan
		return diff.SetNew(field, newEntriesSet)
	}
}
