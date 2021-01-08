package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

func init() {
	resource.AddTestSweepers("linode_token", &resource.Sweeper{
		Name: "linode_token",
		F:    testSweepLinodeToken,
	})
}

func testSweepLinodeToken(prefix string) error {
	client, err := getClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := sweeperListOptions(prefix, "label")
	tokens, err := client.ListTokens(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting tokens: %s", err)
	}
	for _, token := range tokens {
		if !shouldSweepAcceptanceTestResource(prefix, token.Label) {
			continue
		}
		err := client.DeleteToken(context.Background(), token.ID)

		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", token.Label, err)
		}
	}

	return nil
}

func TestAccLinodeToken_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_token.foobar"
	var tokenName = acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeTokenConfigBasic(tokenName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeTokenExists,
					resource.TestCheckResourceAttr(resName, "label", tokenName),
					resource.TestCheckResourceAttr(resName, "expiry", "2100-01-02T03:04:05Z"),
					resource.TestCheckResourceAttrSet(resName, "scopes"),
					resource.TestCheckResourceAttrSet(resName, "token"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			{
				Config: testAccCheckLinodeTokenConfigUpdates(tokenName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeTokenExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", tokenName)),
				),
			},
		},
	})
}

func testAccCheckLinodeTokenExists(s *terraform.State) error {
	client := testAccProvider.Meta().(linodego.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_token" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetToken(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Token %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func testAccCheckLinodeTokenDestroy(s *terraform.State) error {
	client, ok := testAccProvider.Meta().(linodego.Client)
	if !ok {
		return fmt.Errorf("Error getting Linode client")
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_token" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetToken(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Token with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Token with id %d", id)
		}
	}

	return nil
}

func testAccCheckLinodeTokenConfigBasic(token string) string {
	return fmt.Sprintf(`
	resource "linode_token" "foobar" {
		label = "%s"
		scopes = "linodes:read_only"
		expiry = "2100-01-02T03:04:05Z"
	}`, token)
}

func testAccCheckLinodeTokenConfigUpdates(token string) string {
	return fmt.Sprintf(`
	resource "linode_token" "foobar" {
		label = "%s_renamed"
		scopes = "linodes:read_only"
		expiry = "2100-01-02T03:04:05Z"
	}`, token)
}
