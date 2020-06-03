package linode

import (
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/acctest"
	acctesthelpers "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// publicKeyMaterial for use while testing
var publicKeyMaterial string

func init() {
	var err error
	publicKeyMaterial, _, err = acctesthelpers.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		log.Fatalf("Failed to generate random SSH key pair for testing: %s", err)
	}
}

func TestMain(m *testing.M) {
	acctest.UseBinaryDriver("linode", Provider)
	resource.TestMain(m)
}
