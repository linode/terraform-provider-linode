package helper

import (
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"
)

func CompareTimeStrings(t1, t2, timeFormat string) bool {
	parsedT1, err := time.Parse(timeFormat, t1)
	if err != nil {
		return false
	}

	parsedT2, err := time.Parse(timeFormat, t2)
	if err != nil {
		return false
	}

	return parsedT1.Equal(parsedT2)
}

func CompareRFC3339TimeStrings(t1, t2 string) bool {
	return CompareTimeStrings(t1, t2, time.RFC3339)
}

func CompareTimeWithTimeString(t1 *time.Time, t2 string, timeFormat string) bool {
	parsedT2, err := time.Parse(timeFormat, t2)
	if err != nil {
		return false
	}

	return t1.Equal(parsedT2)
}

// Check if two string lists contain the same elements regardless of order
func StringListElementsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	freq := make(map[string]int)
	for _, v := range a {
		freq[v]++
	}
	for _, v := range b {
		freq[v]--
		if freq[v] < 0 {
			return false
		}
	}
	return true
}

// Check if `subset` is a subset of `superset`, or in other words, assuming no duplicated items are in the sets,
// whether slice `superset` contains all elements of slice `subset`.
func ValidateStringSubset(superset, subset []string) bool {
	return ValidateSubset(TypedSliceToAny(superset), TypedSliceToAny(subset))
}

// Check if `subset` is a subset of `superset`, or in other words, whether slice `superset` contains all elements of slice `subset`,
// assuming no duplicated items are in the sets.
func ValidateSubset(superset, subset []any) bool {
	for _, v := range subset {
		if !slices.Contains(superset, v) {
			return false
		}
	}

	return true
}

// Check if two slices are equivalent without considering ordering,
// assuming no duplicated items are in the sets.
func CompareSets(a, b []any) bool {
	return CompareSlices(true, true, a, b)
}

// Check if two string slices are equivalent without considering ordering,
// assuming no duplicated items are in the sets.
func CompareStringSets(a, b []string) bool {
	return CompareSets(TypedSliceToAny(a), TypedSliceToAny(b))
}

func CompareScopes(s1, s2 string) bool {
	s1AccountScope := s1 == "*"
	s2AccountScope := s2 == "*"
	if s1AccountScope != s2AccountScope {
		return false
	}
	if s1AccountScope && s2AccountScope {
		return true
	}

	s1List := strings.Split(s1, " ")
	s2List := strings.Split(s2, " ")
	return StringListElementsEqual(s1List, s2List)
}

func CompareSlices(ignoreNil, unordered bool, a, b []any) bool {
	if ignoreNil && len(a) == 0 && len(b) == 0 {
		return true
	}

	if unordered {
		less := func(i, j int) bool {
			return fmt.Sprintf("%v", a[i]) < fmt.Sprintf("%v", a[j])
		}
		sort.Slice(a, less)
		sort.Slice(b, less)
	}

	return reflect.DeepEqual(a, b)
}
