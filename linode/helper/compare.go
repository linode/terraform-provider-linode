package helper

import (
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
