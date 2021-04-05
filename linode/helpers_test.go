package linode

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"golang.org/x/crypto/ssh"
)

const optInTestsEnvVar = "ACC_OPT_IN_TESTS"

var optInTests map[string]struct{}

func init() {
	optInTests = make(map[string]struct{})
	optInTestsValue, ok := os.LookupEnv(optInTestsEnvVar)
	if !ok {
		return
	}

	for _, testName := range strings.Split(optInTestsValue, ",") {
		optInTests[testName] = struct{}{}
	}
}

func optInTest(t *testing.T) {
	t.Helper()

	if _, ok := optInTests[t.Name()]; !ok {
		t.Skipf("skipping opt-in test; specify test in environment variable %q to run", optInTestsEnvVar)
	}
}

func getSSHClient(t *testing.T, user, addr string) (client *ssh.Client) {
	t.Helper()

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyMaterial))
	if err != nil {
		t.Fatalf("failed to parse private key: %s", err)
	}
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Minute,
	}

	attempts := 3

	for attempts != 0 {
		client, err = ssh.Dial("tcp", addr+":22", config)
		if err == nil {
			break
		}

		t.Logf("ssh dial failed: %s", err)
		attempts--

		time.Sleep(5 * time.Second)
	}

	if client == nil {
		t.Fatal("failed to get ssh client")
	}
	return
}

func testAccCheckResourceAttrNotEqual(resName string, path, notValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		if value, ok := rs.Primary.Attributes[path]; !ok {
			return fmt.Errorf("attribute %s does not exist", path)
		} else if value == notValue {
			return fmt.Errorf("attribute was equal")
		}

		return nil
	}
}

func testAccCheckResourceNonEmptyList(resourceName, attrName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		instCount, err := strconv.Atoi(rs.Primary.Attributes[fmt.Sprintf("%s.#", attrName)])
		if err != nil {
			return fmt.Errorf("failed to parse: %s", err)
		}

		if instCount < 1 {
			return fmt.Errorf("expected at least 1 element in %s", attrName)
		}

		return nil
	}
}
