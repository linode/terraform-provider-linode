package linode

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// publicKeyMaterial for use while testing
var (
	privateKeyMaterial string
	publicKeyMaterial  string
)

func init() {
	var err error
	publicKeyMaterial, privateKeyMaterial, err = acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		log.Fatalf("Failed to generate random SSH key pair for testing: %s", err)
	}
}
