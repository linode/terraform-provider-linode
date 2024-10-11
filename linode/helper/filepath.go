package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands the given path, replacing any ~'s with the user's
// home directory.
func ExpandPath(path string) (string, error) {
	homePath, homePathErr := os.UserHomeDir()

	segments := strings.Split(path, string(os.PathSeparator))
	result := make([]string, len(segments))

	for i, segment := range segments {
		if segment == "~" {
			// We don't check this error above because we don't want to raise a home resolution error
			// if the user isn't using ~'s.
			if homePathErr != nil {
				return "", fmt.Errorf("could not expand home path: %w", homePathErr)
			}

			segment = homePath
		}

		result[i] = segment
	}

	return filepath.Join(result...), nil
}
