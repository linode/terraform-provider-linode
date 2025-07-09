package acceptance

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func SweeperListOptions(prefix, field string) *linodego.ListOptions {
	filterFmt := "{ %q : {\"+contains\": %q }}"

	filter := fmt.Sprintf(filterFmt, field, prefix)
	listOpts := linodego.NewListOptions(0, filter)
	return listOpts
}

func ShouldSweep(prefix, name string) bool {
	loweredName := strings.ToLower(name)
	if len(prefix) < 3 {
		log.Printf("Ignoring Resource %q because sweeper prefix is too short %q", name, prefix)
		return false
	}

	if !strings.HasPrefix(loweredName, prefix) && !strings.HasPrefix(loweredName, "renamed-"+prefix) {
		log.Printf("Ignoring Resource %q as it doesn't start with `(renamed-)?%s`", name, prefix)
		return false
	}

	return true
}
