package helper

import "log"

func Must[T any](result T, err error) T {
	if err != nil {
		log.Fatalf("helper.Must failed: %s", err)
	}

	return result
}

// StringSet is a set of strings implemented as a map.
// Use ExistsInSet as the value to indicate membership.
type StringSet = map[string]struct{}

// ExistsInSet is used as a value in set maps to indicate membership.
// Example: set[key] = ExistsInSet
var ExistsInSet = struct{}{}
