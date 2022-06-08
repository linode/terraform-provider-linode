package token_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/linode/token/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_token", &resource.Sweeper{
		Name: "linode_token",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "label")
	tokens, err := client.ListTokens(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting tokens: %s", err)
	}
	for _, token := range tokens {
		if !acceptance.ShouldSweep(prefix, token.Label) {
			continue
		}
		err := client.DeleteToken(context.Background(), token.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", token.Label, err)
		}
	}

	return nil
}

func TestAccResourceToken_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_token.foobar"
	tokenName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, tokenName),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
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
				Config: tmpl.Updates(t, tokenName),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
					resource.TestCheckResourceAttr(resName, "label", fmt.Sprintf("%s_renamed", tokenName)),
				),
			},
		},
	})
}

func checkTokenExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

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

func checkTokenDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
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
