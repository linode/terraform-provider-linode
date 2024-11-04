//go:build integration || instance

package instance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/instance/tmpl"
)

func TestAccDataSourceInstances_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(64)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, instanceName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.id"),
					resource.TestCheckResourceAttr(resName, "instances.0.type", "g6-nanode-1"),
					resource.TestCheckResourceAttr(resName, "instances.0.tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "instances.0.image", acceptance.TestImageLatest),
					resource.TestCheckResourceAttr(resName, "instances.0.region", testRegion),
					resource.TestCheckResourceAttr(resName, "instances.0.group", "tf_test"),
					resource.TestCheckResourceAttr(resName, "instances.0.swap_size", "256"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.disk_encryption"),
					resource.TestCheckResourceAttr(resName, "instances.0.ipv4.#", "2"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.ipv6"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.host_uuid"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.has_user_data"),
					resource.TestCheckResourceAttr(resName, "instances.0.disk.#", "2"),
					resource.TestCheckResourceAttr(resName, "instances.0.config.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.config.0.id"),
					resource.TestCheckResourceAttr(resName, "instances.0.placement_group.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceInstances_withBlockStorageEncryption(t *testing.T) {
	t.Parallel()

	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	// Resolve a region with support for Block Storage Encryption
	targetRegion, err := acceptance.GetRandomRegionWithCaps(
		[]string{"Linodes", "Block Storage Encryption"},
		"core",
	)
	if err != nil {
		t.Fatal(err)
	}

	rootPass := acctest.RandString(64)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithBlockStorageEncryption(t, instanceName, targetRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.id"),
					resource.TestCheckTypeSetElemAttr(resName, "instances.0.capabilities.*", "Block Storage Encryption"),
				),
			},
		},
	})
}

func TestAccDataSourceInstances_withPG(t *testing.T) {
	t.Parallel()

	resName := "data.linode_instances.foobar"
	instanceName := acctest.RandomWithPrefix("tf_test")

	pgIDs := []string{"foobar"}

	// Resolve a region with support for PGs
	targetRegion, err := acceptance.GetRandomRegionWithCaps(
		[]string{"Linodes", "Placement Group"},
		"core",
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,

		Steps: []resource.TestStep{
			{
				Config: tmpl.DataWithPG(t, instanceName, targetRegion, "foobar", pgIDs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "1"),
					resource.TestCheckResourceAttr(resName, "instances.0.placement_group.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.placement_group.0.id"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.placement_group.0.placement_group_type"),
					resource.TestCheckResourceAttrSet(resName, "instances.0.placement_group.0.placement_group_policy"),
				),
			},
		},
	})
}

func TestAccDataSourceInstances_multipleInstances(t *testing.T) {
	resName := "data.linode_instances.foobar"
	resNameDesc := "data.linode_instances.desc"
	resNameAsc := "data.linode_instances.asc"

	instanceName := acctest.RandomWithPrefix("tf_test")
	tagName := acctest.RandomWithPrefix("tf_test")
	rootPass := acctest.RandString(64)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             acceptance.CheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataMultiple(t, instanceName, tagName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "3"),
				),
			},
			{
				Config: tmpl.DataMultipleOrder(t, instanceName, tagName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					// Ensure order is correctly appended to filter
					resource.TestCheckResourceAttr(resNameDesc, "instances.#", "3"),
					resource.TestCheckResourceAttr(resNameAsc, "instances.#", "3"),
				),
			},
			{
				Config: tmpl.DataMultipleRegex(t, instanceName, tagName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "3"),
				),
			},
			{
				Config: tmpl.DataClientFilter(t, instanceName, tagName, testRegion, rootPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "instances.#", "1"),
					resource.TestCheckResourceAttr(resName, "instances.0.status", "running"),
				),
			},
		},
	})
}
