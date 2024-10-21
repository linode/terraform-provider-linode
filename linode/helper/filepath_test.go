//go:build unit

package helper

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExpandPath(t *testing.T) {
	homePath, err := os.UserHomeDir()
	require.NoError(t, err)

	expandedPath, err := ExpandPath(filepath.Join("~", "foo", "bar"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(homePath, "foo", "bar"), expandedPath)

	expandedPath, err = ExpandPath("")
	require.NoError(t, err)
	require.Equal(t, "", expandedPath)

	// /foo/bar
	absPath := strings.Join([]string{"", "foo", "bar"}, string(os.PathSeparator))

	expandedPath, err = ExpandPath(absPath)
	require.NoError(t, err)
	require.Equal(t, absPath, expandedPath)
}
