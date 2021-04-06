package linode

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

const providerKeySkipInstanceReadyPoll = "skip_instance_ready_poll"

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"linode": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LINODE_TOKEN"); v == "" {
		t.Fatal("LINODE_TOKEN must be set for acceptance tests")
	}
}

func accTestWithProvider(config string, options map[string]interface{}) string {
	sb := strings.Builder{}
	sb.WriteString("provider \"linode\" {\n")
	for key, value := range options {
		sb.WriteString(fmt.Sprintf("\t%s = %#v\n", key, value))
	}
	sb.WriteString("}\n")
	return sb.String() + config
}
