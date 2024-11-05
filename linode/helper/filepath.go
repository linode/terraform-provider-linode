package helper

import (
	"fmt"
	"os"
	"strings"
)

const pathSeparatorString = string(os.PathSeparator)

// ExpandPath expands the given path, replacing ~'s with the user's
// home directory.
// NOTE: This does not implement feature-complete tilde expansion.
func ExpandPath(path string) (string, error) {
	segments := strings.Split(path, pathSeparatorString)

	if segments[0] == "~" {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not expand home path: %w", err)
		}

		segments[0] = homePath
	}

	// We don't use filepath.Join(...) here because it does not
	// support rebuilding paths starting with `/`.
	return strings.Join(segments, pathSeparatorString), nil
}
