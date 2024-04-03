//go:build integration || nbconfig

package nbconfig_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/v2/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v2/linode/helper"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig"
	"github.com/linode/terraform-provider-linode/v2/linode/nbconfig/tmpl"
)

var testRegion string

func init() {
	region, err := acceptance.GetRandomRegionWithCaps([]string{"nodebalancers"})
	if err != nil {
		log.Fatal(err)
	}

	testRegion = region
}

func TestAccResourceNodeBalancerConfig_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		CheckDestroy:              checkNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:       tmpl.Basic(t, nodebalancerName, testRegion),
				ResourceName: resName,
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
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: resourceImportStateID,
			},
		},
	})
}

func TestAccResourceNodeBalancerConfig_ssl(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		CheckDestroy:              checkNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:       tmpl.SSL(t, nodebalancerName, testRegion, tmpl.TestCertifcate, tmpl.TestPrivateKey),
				ResourceName: resName,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTPS)),
					resource.TestCheckResourceAttrSet(resName, "ssl_cert"),
					resource.TestCheckResourceAttrSet(resName, "ssl_key"),
				),
			},
			{
				ResourceName:            resName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ssl_cert", "ssl_key"},
				ImportStateIdFunc:       resourceImportStateID,
			},
		},
	})
}

func TestAccResourceNodeBalancerConfig_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8080"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/"),
					resource.TestCheckResourceAttr(resName, "check_passive", "true"),

					resource.TestCheckResourceAttrSet(resName, "stickiness"),
					resource.TestCheckResourceAttrSet(resName, "check_attempts"),
					resource.TestCheckResourceAttrSet(resName, "check_timeout"),
				),
			},
			{
				Config: tmpl.Updates(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8088"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/foo"),
					resource.TestCheckResourceAttr(resName, "check_attempts", "3"),
					resource.TestCheckResourceAttr(resName, "check_timeout", "30"),
					resource.TestCheckResourceAttr(resName, "check_interval", "31"),
					resource.TestCheckResourceAttr(resName, "check_passive", "false"),

					resource.TestCheckResourceAttr(resName, "stickiness", string(linodego.StickinessHTTPCookie)),
				),
			},
		},
	})
}

func TestAccResourceNodeBalancerConfig_proxyProtocol(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             checkNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ProxyProtocol(t, nodebalancerName, testRegion),
				Check: resource.ComposeTestCheckFunc(
					checkNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "80"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolTCP)),
					resource.TestCheckResourceAttr(resName, "proxy_protocol", string(linodego.ProxyProtocolV2)),
				),
			},
		},
	})
}

func TestLinodeNodeBalancerConfig_UpgradeV0(t *testing.T) {
	t.Parallel()

	oldState := map[string]interface{}{
		"node_status": map[string]interface{}{
			"down": "13",
			"up":   "37",
		},
	}

	desiredState := map[string]interface{}{
		"node_status": []map[string]interface{}{
			{
				"down": 13,
				"up":   37,
			},
		},
	}

	newState, err := nbconfig.ResourceNodeBalancerConfigV0Upgrade(context.Background(), oldState, nil)
	if err != nil {
		t.Fatalf("error migrating state: %v", err)
	}

	if !reflect.DeepEqual(desiredState, newState) {
		t.Fatalf("expected %v, got %v", desiredState, newState)
	}
}

func TestLinodeNodeBalancerConfig_UpgradeV0Empty(t *testing.T) {
	t.Parallel()

	oldState := map[string]interface{}{
		"node_status": map[string]interface{}{
			"down": "",
			"up":   "",
		},
	}

	desiredState := map[string]interface{}{
		"node_status": []map[string]interface{}{
			{
				"down": 0,
				"up":   0,
			},
		},
	}

	newState, err := nbconfig.ResourceNodeBalancerConfigV0Upgrade(context.Background(), oldState, nil)
	if err != nil {
		t.Fatalf("error migrating state: %v", err)
	}

	if !reflect.DeepEqual(desiredState, newState) {
		t.Fatalf("expected %v, got %v", desiredState, newState)
	}
}

func checkNodeBalancerConfigExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}

		_, err = client.GetNodeBalancerConfig(context.Background(), nodebalancerID, id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of NodeBalancer Config %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkNodeBalancerConfigDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}

		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetNodeBalancerConfig(context.Background(), nodebalancerID, id)

		if err == nil {
			return fmt.Errorf("NodeBalancer Config with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting NodeBalancer Config with id %d", id)
		}
	}

	return nil
}

func resourceImportStateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_nodebalancer_config" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		nodebalancerID, err := strconv.Atoi(rs.Primary.Attributes["nodebalancer_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing nodebalancer_id %v to int", rs.Primary.Attributes["nodebalancer_id"])
		}
		return fmt.Sprintf("%d,%d", nodebalancerID, id), nil
	}

	return "", fmt.Errorf("Error finding linode_nodebalancer_config")
}
