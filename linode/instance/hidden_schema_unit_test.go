//go:build unit

package instance

import (
	"os"
	"testing"
)

func TestHiddenHostID_absentByDefault(t *testing.T) {
	os.Unsetenv("LINODE_ENABLE_HIDDEN_HOST_ID")
	if _, ok := resourceSchemaWithHiddenAttributes()["host_id"]; ok {
		t.Fatal("host_id MUST NOT be present when env var is unset")
	}
	if _, ok := Resource().Schema["host_id"]; ok {
		t.Fatal("host_id MUST NOT be in Resource() schema when env var unset")
	}
}

func TestHiddenHostID_presentWhenEnabled(t *testing.T) {
	t.Setenv("LINODE_ENABLE_HIDDEN_HOST_ID", "1")
	if _, ok := resourceSchemaWithHiddenAttributes()["host_id"]; !ok {
		t.Fatal("host_id MUST be present when env var is set")
	}
	// base schema must remain unmutated
	if _, ok := resourceSchema["host_id"]; ok {
		t.Fatal("base resourceSchema was mutated!")
	}
}
