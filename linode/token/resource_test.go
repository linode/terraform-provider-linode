//go:build integration || token

package token_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/token/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_token", &resource.Sweeper{
		Name: "linode_token",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
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
		PreCheck:                 func() { acceptance.PreCheck(t) },
		CheckDestroy:             checkTokenDestroy,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, tokenName),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
					resource.TestCheckResourceAttr(resName, "label", tokenName),
					resource.TestCheckResourceAttr(resName, "expiry", "2100-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(resName, "scopes", "linodes:read_only"),
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

func TestAccResourceToken_recreative_update(t *testing.T) {
	t.Parallel()

	resName := "linode_token.foobar"
	tokenName := acctest.RandomWithPrefix("tf_test")

	var currentToken string
	tokenRecreatedCheck := func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "linode_token" {
				continue
			}

			newToken, ok := rs.Primary.Attributes["token"]
			if !ok {
				return fmt.Errorf("Can't find the token in the state.")
			}
			if newToken == currentToken {
				return fmt.Errorf("The token suppose to be but was not recreated.")
			}

		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		CheckDestroy:             checkTokenDestroy,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, tokenName),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
					tokenRecreatedCheck,
					resource.TestCheckResourceAttr(resName, "label", tokenName),
					resource.TestCheckResourceAttr(resName, "expiry", "2100-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(resName, "scopes", "linodes:read_only"),
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
				Config: tmpl.RecreateNewExpiryDate(t, tokenName, "2099-05-04T03:02:01+00:00"),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
					tokenRecreatedCheck,
					resource.TestCheckResourceAttr(resName, "expiry", "2099-05-04T03:02:01+00:00"),
				),
			},
			{
				Config: tmpl.RecreateNewScopes(t, tokenName, "linodes:read_only lke:read_only"),
				Check: resource.ComposeTestCheckFunc(
					checkTokenExists,
					tokenRecreatedCheck,
					resource.TestCheckResourceAttr(resName, "scopes", "linodes:read_only lke:read_only"),
				),
			},
		},
	})
}

func checkTokenExists(s *terraform.State) error {
	client := acceptance.TestAccFrameworkProvider.Meta.Client

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
	client := acceptance.TestAccFrameworkProvider.Meta.Client
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
