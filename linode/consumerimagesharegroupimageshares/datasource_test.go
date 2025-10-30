//go:build integration || consumerimagesharegroupimageshares

package consumerimagesharegroupimageshares_test

import (
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/consumerimagesharegroupimageshares/tmpl"
)

// This test requires two separate Linode API tokens, one for the producer
// and one for the consumer.
//
// These can be set using the LINODE_PRODUCER_TOKEN and LINODE_CONSUMER_TOKEN
// environment variables.
//
// If either is not set,the test will be skipped.
func TestAccDataSourceImageShareGroupImageShares_basic(t *testing.T) {
	t.Parallel()

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

	const dsAll = "data.linode_consumer_image_share_group_image_shares.all"
	const dsByLabel = "data.linode_consumer_image_share_group_image_shares.by_label"

	fwLabel := acctest.RandomWithPrefix("tf_test")
	instanceLabel := acctest.RandomWithPrefix("tf_test")

	instanceRegion, err := acceptance.GetRandomRegionWithCaps([]string{}, "core")
	if err != nil {
		log.Fatal(err)
	}

	imageLabel1 := acctest.RandomWithPrefix("tf-test")
	imageLabel2 := acctest.RandomWithPrefix("tf-test")
	shareGroupLabel := acctest.RandomWithPrefix("tf-test")
	tokenLabel := acctest.RandomWithPrefix("tf-test")
	memberLabel := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode-producer": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](producerProvider)
			},
			"linode-consumer": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](consumerProvider)
			},
		},
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, fwLabel, instanceLabel, instanceRegion, imageLabel1, imageLabel2, shareGroupLabel, tokenLabel, memberLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsAll, "image_shares.#", "2"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.0.label", "image_one_label"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.0.description", "image one description"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.1.label", "image_two_label"),
					resource.TestCheckResourceAttr(dsAll, "image_shares.1.description", "image two description"),

					resource.TestCheckResourceAttr(dsByLabel, "image_shares.#", "1"),
					resource.TestCheckResourceAttr(dsByLabel, "image_shares.0.label", "image_two_label"),
					resource.TestCheckResourceAttr(dsByLabel, "image_shares.0.description", "image two description"),
				),
			},
		},
	})
}
