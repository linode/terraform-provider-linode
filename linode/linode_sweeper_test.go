package linode

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/linode/linodego"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func getClientForSweepers() (*linodego.Client, error) {
	token := os.Getenv("LINODE_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("LINODE_TOKEN must be set for acceptance tests")
	}

	config := &Config{AccessToken: token, APIVersion: "v4beta"}
	client := config.Client()
	return &client, nil
}

func sweeperListOptions(prefix, field string) *linodego.ListOptions {
	filterFmt := "{ %q : {\"+contains\": %q }}"

	filter := fmt.Sprintf(filterFmt, field, prefix)
	listOpts := linodego.NewListOptions(0, filter)
	return listOpts
}

func shouldSweepAcceptanceTestResource(prefix, name string) bool {
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
