package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands the given path, replacing ~'s with the user's
// home directory.
// NOTE: This does not implement feature-complete tilde expansion.
func ExpandPath(path string) (string, error) {
	segments := strings.Split(path, string(os.PathSeparator))

	if segments[0] == "~" {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not expand home path: %w", err)
		}

		segments[0] = homePath
	}

	return filepath.Join(segments...), nil
}
