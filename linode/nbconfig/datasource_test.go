//go:build integration || nbconfig

package nbconfig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig/tmpl"
)

func TestAccDataSourceNodeBalancerConfig_basic(t *testing.T) {
	t.Parallel()

	resName := "data.linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.DataBasic(t, nodebalancerName, testRegion),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/"),

					resource.TestCheckResourceAttrSet(resName, "algorithm"),
					resource.TestCheckResourceAttrSet(resName, "stickiness"),
					resource.TestCheckResourceAttrSet(resName, "check_attempts"),
					resource.TestCheckResourceAttrSet(resName, "check_timeout"),
					resource.TestCheckResourceAttrSet(resName, "check_interval"),
					resource.TestCheckResourceAttrSet(resName, "check_passive"),
					resource.TestCheckResourceAttrSet(resName, "cipher_suite"),
					resource.TestCheckNoResourceAttr(resName, "ssl_common"),
					resource.TestCheckNoResourceAttr(resName, "ssl_ciphersuite"),
					resource.TestCheckResourceAttr(resName, "node_status.0.up", "0"),
					resource.TestCheckResourceAttr(resName, "node_status.0.down", "0"),
					resource.TestCheckNoResourceAttr(resName, "ssl_cert"),
					resource.TestCheckNoResourceAttr(resName, "ssl_key"),
				),
			},
		},
	})
}
