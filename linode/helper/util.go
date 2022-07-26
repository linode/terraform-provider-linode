package helper

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseMultiSegmentID(id string, numSegments int) ([]int, error) {
	s := strings.Split(id, "/")
	result := make([]int, len(s))

	if len(s) != numSegments {
		return nil, fmt.Errorf("invalid number of id segments")
	}

	for i, seg := range s {
		segID, err := strconv.Atoi(seg)
		if err != nil {
			return nil, err
		}

		result[i] = segID
	}

	return result, nil
}

func FormatMultiSegmentID(ids ...int) string {
	idsStr := make([]string, len(ids))
	for i, id := range ids {
		idsStr[i] = strconv.Itoa(id)
	}

	return strings.Join(idsStr, "/")
}
