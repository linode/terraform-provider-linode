package linode

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceLinodeNodeBalancer_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer.foobar"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinodeNodeBalancerBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerExists,
					resource.TestCheckResourceAttr(resName, "label", nodebalancerName),
					resource.TestCheckResourceAttr(resName, "client_conn_throttle", "20"),
					resource.TestCheckResourceAttr(resName, "region", "us-east"),
					resource.TestCheckResourceAttrSet(resName, "hostname"),
					resource.TestCheckResourceAttrSet(resName, "ipv4"),
					resource.TestCheckResourceAttrSet(resName, "ipv6"),
					resource.TestCheckResourceAttrSet(resName, "created"),
					resource.TestCheckResourceAttrSet(resName, "updated"),
					resource.TestCheckResourceAttr(resName, "transfer.#", "1"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.in"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.out"),
					resource.TestCheckResourceAttrSet(resName, "transfer.0.total"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
				),
			},
		},
	})
}

func testDataSourceLinodeNodeBalancerBasic(nodeBalancerName string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodeBalancerName) + `
data "linode_nodebalancer" "foobar" {
	id = "${linode_nodebalancer.foobar.id}"
}
`
}
