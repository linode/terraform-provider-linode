//go:build integration || producerimagesharegroupmember

package producerimagesharegroupmember_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroupmember/tmpl"
)

// This test requires two separate Linode API tokens, one for the producer
// and one for the consumer.
//
// These can be set using the LINODE_PRODUCER_TOKEN and LINODE_CONSUMER_TOKEN
// environment variables.
//
// If either is not set,the test will be skipped.
func TestAccDataSourceImageShareGroupMember_basic(t *testing.T) {
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

	resourceName := "data.linode_producer_image_share_group_member.foobar"
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
				Config: tmpl.DataBasic(t, shareGroupLabel, tokenLabel, memberLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "sharegroup_id"),
					resource.TestCheckResourceAttrSet(resourceName, "token_uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttr(resourceName, "label", memberLabel),
				),
			},
		},
	})
}
