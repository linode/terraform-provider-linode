//go:build integration || regions

package regions_test

import (
	"context"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/regions/tmpl"
)

func TestAccDataSourceRegions_basic_smoke(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_regions.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.country"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv6"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.0.capabilities.#", 0),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceRegions_filterByCountry(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_regions.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	regions, err := client.ListRegions(context.TODO(), nil)
	randIndex := rand.Intn(len(regions))
	region := regions[randIndex]

	country := region.Country
	status := region.Status
	capabilities := region.Capabilities

	randomCapability := capabilities[rand.Intn(len(capabilities))]

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing region.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterCountry(t, country, status, randomCapability),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "regions.0.country", country),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv6"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.0.capabilities.#", 0),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceRegions_filterByStatus(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_regions.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	regions, err := client.ListRegions(context.TODO(), nil)
	randIndex := rand.Intn(len(regions))
	region := regions[randIndex]

	country := region.Country
	status := region.Status
	capabilities := region.Capabilities

	randomCapability := capabilities[rand.Intn(len(capabilities))]

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing region.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterStatus(t, country, status, randomCapability),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.country"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.label"),
					resource.TestCheckResourceAttr(resourceName, "regions.0.status", status),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv6"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.0.capabilities.#", 0),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.#", 0),
				),
			},
		},
	})
}

func TestAccDataSourceRegions_filterByCapabilities(t *testing.T) {
	t.Parallel()

	resourceName := "data.linode_regions.foobar"

	client, err := acceptance.GetTestClient()
	if err != nil {
		t.Fail()
		t.Log("Failed to get testing client.")
	}

	regions, err := client.ListRegions(context.TODO(), nil)
	randIndex := rand.Intn(len(regions))
	region := regions[randIndex]

	country := region.Country
	status := region.Status
	capabilities := region.Capabilities

	randomCapability := capabilities[rand.Intn(len(capabilities))]

	if err != nil {
		t.Fail()
		t.Log("Failed to get testing region.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataFilterCapabilities(t, country, status, randomCapability),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.country"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.label"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "regions.0.resolvers.0.ipv6"),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.0.capabilities.#", 0),
					acceptance.CheckResourceAttrGreaterThan(resourceName, "regions.#", 0),
					acceptance.LoopThroughStringList(resourceName, "regions", func(resourceName, path string, state *terraform.State) error {
						return acceptance.CheckResourceAttrListContains(resourceName, path+".capabilities", randomCapability)(state)
					}),
				),
			},
		},
	})
}
