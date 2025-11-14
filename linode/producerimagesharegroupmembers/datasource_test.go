//go:build integration || producerimagesharegroupmembers

package producerimagesharegroupmembers_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmembers/tmpl"
)

// This test requires two separate Linode API tokens, one for the producer
// and one for the consumer.
//
// These can be set using the LINODE_PRODUCER_TOKEN and LINODE_CONSUMER_TOKEN
// environment variables.
//
// If either is not set,the test will be skipped.
func TestAccDataSourceImageShareGroupMembers_basic(t *testing.T) {
	t.Parallel()

	const dsByLabel = "data.linode_producer_image_share_group_members.by_label"
	const dsByStatus = "data.linode_producer_image_share_group_members.by_status"
	const dsByTokenUUID = "data.linode_producer_image_share_group_members.by_token_uuid"

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
	memberLabel := acctest.RandomWithPrefix("tf-test")

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
				Config: tmpl.DataBasic(t, shareGroupLabel, tokenLabel, memberLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsByLabel, "members.#", "1"),
					resource.TestCheckResourceAttrSet(dsByLabel, "members.0.sharegroup_id"),
					resource.TestCheckResourceAttrSet(dsByLabel, "members.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByLabel, "members.0.status"),
					resource.TestCheckResourceAttr(dsByLabel, "members.0.label", memberLabel),

					resource.TestCheckResourceAttr(dsByTokenUUID, "members.#", "1"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "members.0.sharegroup_id"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "members.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByTokenUUID, "members.0.status"),
					resource.TestCheckResourceAttr(dsByTokenUUID, "members.0.label", memberLabel),

					resource.TestCheckResourceAttr(dsByStatus, "members.#", "1"),
					resource.TestCheckResourceAttrSet(dsByStatus, "members.0.sharegroup_id"),
					resource.TestCheckResourceAttrSet(dsByStatus, "members.0.token_uuid"),
					resource.TestCheckResourceAttrSet(dsByStatus, "members.0.status"),
					resource.TestCheckResourceAttr(dsByStatus, "members.0.label", memberLabel),
				),
			},
		},
	})
}
