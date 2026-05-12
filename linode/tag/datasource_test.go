//go:build integration || tag

package tag_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/tag/tmpl"
)

func init() {
	resource.AddTestSweepers("linode_tag", &resource.Sweeper{
		Name: "linode_tag",
		F:    sweep,
	})
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
	var reservedIPAddress string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					region, err := acceptance.GetRandomRegionWithCaps(nil, "core")
					if err != nil {
						t.Fatalf("Error finding region: %s", err)
					}
					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatalf("Error getting client: %s", err)
					}
					ip, err := client.AllocateReserveIP(context.Background(), linodego.AllocateReserveIPOptions{
						Type:     "ipv4",
						Public:   true,
						Reserved: true,
						Region:   region,
					})
					if err != nil {
						t.Fatalf("Error reserving IP: %s", err)
					}
					reservedIPAddress = ip.Address
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

					// Poll until the reserved IP association is visible in the API
					// before letting Terraform proceed (eventual consistency).
					// Skip if the API environment does not support reserved_ipv4_addresses
					// in POST /tags (feature may not be rolled out yet).
					deadline := time.Now().Add(30 * time.Second)
					visible := false
					for time.Now().Before(deadline) {
						objects, err := client.ListTaggedObjects(context.Background(), tagLabel, nil)
						if err == nil && len(objects) > 0 {
							visible = true
							break
						}
						time.Sleep(2 * time.Second)
					}
					if !visible {
						t.Skip("reserved_ipv4_addresses tag association not visible via API; skipping (feature may not be available in this environment)")
					}
				},
				Config: tmpl.DataSource(t, tagLabel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "label", tagLabel),
					resource.TestCheckResourceAttr(dsName, "id", tagLabel),
					resource.TestCheckResourceAttr(dsName, "objects.#", "1"),
					resource.TestCheckResourceAttr(dsName, "objects.0.type", "reserved_ipv4_address"),
					resource.TestCheckResourceAttrWith(dsName, "objects.0.id", func(val string) error {
						if val != reservedIPAddress {
							return fmt.Errorf("expected objects.0.id to be %q, got %q", reservedIPAddress, val)
						}
						return nil
					}),
				),
			},
		},
	})
}
