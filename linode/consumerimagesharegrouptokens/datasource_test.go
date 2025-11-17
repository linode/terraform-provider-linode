//go:build integration || consumerimagesharegrouptokens

package consumerimagesharegrouptokens_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegrouptokens/tmpl"
)

// This test requires two separate Linode API tokens, one for the producer
// and one for the consumer.
//
// These can be set using the LINODE_PRODUCER_TOKEN and LINODE_CONSUMER_TOKEN
// environment variables.
//
// If either is not set,the test will be skipped.
func TestAccDataSourceImageShareGroupTokens_basic(t *testing.T) {
	t.Parallel()

	const dsByLabel = "data.linode_consumer_image_share_group_tokens.by_label"
	const dsByStatus = "data.linode_consumer_image_share_group_tokens.by_status"
	const dsByTokenUUID = "data.linode_consumer_image_share_group_tokens.by_token_uuid"
	const dsByValidForShareGroupUUID = "data.linode_consumer_image_share_group_tokens.by_valid_for_sharegroup_uuid"

	producerToken := os.Getenv("LINODE_PRODUCER_TOKEN")
	consumerToken := os.Getenv("LINODE_CONSUMER_TOKEN")

	if producerToken == "" || consumerToken == "" {
		t.Skip("Skipping test: both LINODE_PRODUCER_TOKEN and LINODE_CONSUMER_TOKEN must be set")
	}

	producerClient, err := acceptance.GetTestClientAlternateToken("LINODE_PRODUCER_TOKEN")
	if err != nil {
		t.Fatalf("Failed to create producer client: %s", err)
	}

	consumerClient, err := acceptance.GetTestClientAlternateToken("LINODE_CONSUMER_TOKEN")
	if err != nil {
		t.Fatalf("Failed to create consumer client: %s", err)
	}

	producerProvider := acceptance.NewFrameworkProviderWithClient(producerClient)
	consumerProvider := acceptance.NewFrameworkProviderWithClient(consumerClient)

	shareGroupLabel := acctest.RandomWithPrefix("tf-test")
	tokenLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode-producer": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](producerProvider, nil)
			},
			"linode-consumer": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](consumerProvider, nil)
			},
		},
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, shareGroupLabel, tokenLabel),
				Check: resource.ComposeTestCheckFunc(
					acceptance.CheckResourceAttrGreaterThan(dsByStatus, "tokens.#", 0),
					resource.TestCheckResourceAttrSet(dsByStatus, "tokens.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByStatus, "tokens.0.status"),
					resource.TestCheckResourceAttrSet(dsByStatus, "tokens.0.valid_for_sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByStatus, "tokens.0.sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByStatus, "tokens.0.sharegroup_label"),
					resource.TestCheckResourceAttr(dsByStatus, "tokens.0.label", tokenLabel),

					resource.TestCheckResourceAttr(dsByLabel, "tokens.#", "1"),
					resource.TestCheckResourceAttrSet(dsByLabel, "tokens.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByLabel, "tokens.0.status"),
					resource.TestCheckResourceAttrSet(dsByLabel, "tokens.0.valid_for_sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByLabel, "tokens.0.sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByLabel, "tokens.0.sharegroup_label"),
					resource.TestCheckResourceAttr(dsByLabel, "tokens.0.label", tokenLabel),

					resource.TestCheckResourceAttr(dsByTokenUUID, "tokens.#", "1"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "tokens.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "tokens.0.status"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "tokens.0.valid_for_sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByTokenUUID, "tokens.0.sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByTokenUUID, "tokens.0.sharegroup_label"),
					resource.TestCheckResourceAttr(dsByTokenUUID, "tokens.0.label", tokenLabel),

					resource.TestCheckResourceAttr(dsByValidForShareGroupUUID, "tokens.#", "1"),
					resource.TestCheckResourceAttrSet(dsByValidForShareGroupUUID, "tokens.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByValidForShareGroupUUID, "tokens.0.status"),
					resource.TestCheckResourceAttrSet(dsByValidForShareGroupUUID, "tokens.0.valid_for_sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByValidForShareGroupUUID, "tokens.0.sharegroup_uuid"),
					resource.TestCheckNoResourceAttr(dsByValidForShareGroupUUID, "tokens.0.sharegroup_label"),
					resource.TestCheckResourceAttr(dsByValidForShareGroupUUID, "tokens.0.label", tokenLabel),
				),
			},
		},
	})
}
