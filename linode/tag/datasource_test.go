//go:build integration || tag

package tag_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/tag/tmpl"
)

var testRegion string

func init() {
	resource.AddTestSweepers("linode_tag", &resource.Sweeper{
		Name: "linode_tag",
		F:    sweep,
	})

	region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	tags, err := client.ListTags(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("Error getting tags: %s", err)
	}

	for _, tag := range tags {
		if !acceptance.ShouldSweep(prefix, tag.Label) {
			continue
		}
		if err := client.DeleteTag(context.Background(), tag.Label); err != nil {
			return fmt.Errorf("Error destroying tag %q during sweep: %s", tag.Label, err)
		}
	}

	return nil
}

func TestAccDataSourceTag_basic(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_tag.test"
	tagLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatalf("Error getting client: %s", err)
					}
					t.Cleanup(func() {
						_ = client.DeleteTag(context.Background(), tagLabel)
					})
					if _, err := client.CreateTag(context.Background(), linodego.TagCreateOptions{Label: tagLabel}); err != nil {
						t.Fatalf("Error creating tag: %s", err)
					}
				},
				Config: tmpl.DataSource(t, tagLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "label", tagLabel),
					resource.TestCheckResourceAttr(dsName, "id", tagLabel),
				),
			},
		},
	})
}

func TestAccDataSourceTag_reservedIP(t *testing.T) {
	t.Parallel()

	dsName := "data.linode_tag.test"
	tagLabel := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatalf("Error getting client: %s", err)
					}
					ip, err := client.AllocateReserveIP(context.Background(), linodego.AllocateReserveIPOptions{
						Type:     "ipv4",
						Public:   true,
						Reserved: true,
						Region:   testRegion,
					})
					if err != nil {
						t.Fatalf("Error reserving IP: %s", err)
					}
					t.Cleanup(func() {
						_ = client.DeleteReservedIPAddress(context.Background(), ip.Address)
						_ = client.DeleteTag(context.Background(), tagLabel)
					})
					if _, err := client.CreateTag(context.Background(), linodego.TagCreateOptions{
						Label:                 tagLabel,
						ReservedIPv4Addresses: []string{ip.Address},
					}); err != nil {
						t.Fatalf("Error creating tag: %s", err)
					}
				},
				Config: tmpl.DataSource(t, tagLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "label", tagLabel),
					resource.TestCheckResourceAttr(dsName, "id", tagLabel),
				),
			},
		},
	})
}
