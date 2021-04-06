package linode

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
)

var (
	testCertifcate = `-----BEGIN CERTIFICATE-----
MIIF3DCCA8QCCQC0dUFu1HvjazANBgkqhkiG9w0BAQsFADCBrzELMAkGA1UEBhMC
VVMxCzAJBgNVBAgMAlBBMRUwEwYDVQQHDAxQaGlsYWRlbHBoaWExDzANBgNVBAoM
Bkxpbm9kZTELMAkGA1UECwwCRFgxKDAmBgNVBAMMH2xpbm9kZS1vYmotYnVja2V0
LWNlcnQtdGVzdC54eXoxNDAyBgkqhkiG9w0BCQEWJWFkbWluQGxpbm9kZS1vYmot
YnVja2V0LWNlcnQtdGVzdC54eXowHhcNMjAxMDA1MTg0MDUyWhcNMjExMDA1MTg0
MDUyWjCBrzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlBBMRUwEwYDVQQHDAxQaGls
YWRlbHBoaWExDzANBgNVBAoMBkxpbm9kZTELMAkGA1UECwwCRFgxKDAmBgNVBAMM
H2xpbm9kZS1vYmotYnVja2V0LWNlcnQtdGVzdC54eXoxNDAyBgkqhkiG9w0BCQEW
JWFkbWluQGxpbm9kZS1vYmotYnVja2V0LWNlcnQtdGVzdC54eXowggIiMA0GCSqG
SIb3DQEBAQUAA4ICDwAwggIKAoICAQCy4LqfRYXE314e6YkpR1BbKPH8ohO4lcMt
+YzMUNlOC1KUktGjX8pWk4wAXYar7Mxccmbbh68pgE8iSio8V97CdQb8O64OQmre
/y33z7Yts37/6mH5mBnfeiilVHOenQmh+4400tvF1jljU8MZSg6sLM4ZEBhfcT0V
3yqxAwwzV8vk0t7uLRCMuDI5B4h4ZCsheCkA2roF4RGUG6KwGzf+dLSKzBcjy5ho
h4huzp5jDYer7S86dV6/9Gwzh8CPhVaixbymHGoMbJM8lUtc/hFI+J8WVh/qLTKQ
CcqvoZ96QU0LX2ib+ElvCMGl/UrznpHZUrGkLPfnnoxK/vKBNycJsENtWno9KgtN
fsdmYy/blxNRW/qpi+l92f3zbjjpRqJ/oyA+hsSMn19O/v3O4wz+YS55xnVeEPIf
fOq6VJ9BfVdXPPRp33sllM8EVWuS4ry3oJKI1CFTlhV7eU1RpJmbc5X8GhytiD2M
gIrVlYzJTftSHw7J3v0orRD6SxI9enXI4o4pS1MMxRNb+ZQDvwx3ZujxjFXe3+qI
kme3ih+Vl9W9rDeKAd95ciII9CxBqOvsso8zqDAEV25fn3tutk/7hQNMqv0APAah
Lo/eY1NK9i9YVJknVSzWBkE2MUyvpfFhiw6TPYh88qH+wN3CznWaCtXiAjH3kbOk
6y2OmI8+4QIDAQABMA0GCSqGSIb3DQEBCwUAA4ICAQCP2UawP8GDWxyMOsHDPqKp
PtedCxPpEPsQm8KMnt5KJ55NFqTcpARz1miHXT1aBedu9IoqxvTP4g8BQ4QFjP2s
ddNu2WKqnwyzkCtnB2zOrOKlvUtRAZ4x2iyhKNqls6D7I4tw22HMbTzW2TVeuGVa
oiRtawFcUsjSAcarRw6swLTln+BK54dWa9E5hiulBoHLosMWCEyUDrUnaiB+2+7C
bsExYZTXRlii7YPSr46zPmte2iKa1+b0g5DXkzSazWp+R/dlGYp84uLWk71e4b/9
So1pIitPasCJHgO/ii9nIcmDXarkaGT5CEUP8WPp6mLY5W9NxgF2czdz6AMJa3P9
2jNd4J1VFl8k+LDZ4GnwHGhyL3h3lFUmmoQV/0YVoXmA59SxE2JPvc2d1V6xh2gz
yg2M+xcKliSXxshhAopsSSoEp5g3II2mCvzeSxwsXa4Ob5c5TJNdXslm1pugRCbB
tjFNh70wZmCq+jY8C+vGsDwkf/5UeAd+c+14s3bwsBfWqZBGokVxyf/UWHtsWlVn
p3USWBwLxEWyQIioMmj4O6wROZeyePDlFDVky4hzTCrTS6EFIqkGBs5RneCHhTN0
gNHFG8Ixql6mybJAwopvWGEL+7E4pbNdbhmgVvf2YEQuMZBCM7fGdBsRNkTs6jIA
/8soO6buQgQoCq3GFbodZA==
-----END CERTIFICATE-----
`

	testPrivateKey = `-----BEGIN PRIVATE KEY-----
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
7aJXvW2Q3oCt1sF576QNr9rLxhHl8A==
-----END PRIVATE KEY-----
`
)

func TestAccLinodeNodeBalancerConfig_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	config := testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName)
	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckLinodeNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:       config,
				ResourceName: resName,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
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
				ImportStateIdFunc: testAccStateIDNodeBalancerConfig,
			},
		},
	})
}

func TestAccLinodeNodeBalancerConfig_ssl(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")
	config := testAccCheckLinodeNodeBalancerConfigSSL(nodebalancerName)
	resource.Test(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		CheckDestroy:              testAccCheckLinodeNodeBalancerConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:       config,
				ResourceName: resName,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
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
				ImportStateIdFunc:       testAccStateIDNodeBalancerConfig,
			},
		},
	})
}

func TestAccLinodeNodeBalancerConfig_update(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeNodeBalancerConfigBasic(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
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
				Config: testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "8088"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolHTTP)),
					resource.TestCheckResourceAttr(resName, "check", string(linodego.CheckHTTP)),
					resource.TestCheckResourceAttr(resName, "check_path", "/foo"),
					resource.TestCheckResourceAttr(resName, "check_attempts", "3"),
					resource.TestCheckResourceAttr(resName, "check_timeout", "30"),
					resource.TestCheckResourceAttr(resName, "check_passive", "false"),

					resource.TestCheckResourceAttr(resName, "stickiness", string(linodego.StickinessHTTPCookie)),
				),
			},
		},
	})
}

func TestAccLinodeNodeBalancerConfig_proxyProtocol(t *testing.T) {
	t.Parallel()

	resName := "linode_nodebalancer_config.foofig"
	nodebalancerName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLinodeNodeBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLinodeNodeBalancerConfigProxyProtocol(nodebalancerName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLinodeNodeBalancerConfigExists,
					resource.TestCheckResourceAttr(resName, "port", "80"),
					resource.TestCheckResourceAttr(resName, "protocol", string(linodego.ProtocolTCP)),
					resource.TestCheckResourceAttr(resName, "proxy_protocol", string(linodego.ProxyProtocolV2)),
				),
			},
		},
	})
}

func testAccCheckLinodeNodeBalancerConfigExists(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client

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

func testAccCheckLinodeNodeBalancerConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderMeta).Client
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

func testAccStateIDNodeBalancerConfig(s *terraform.State) (string, error) {
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

func testAccCheckLinodeNodeBalancerConfigBasic(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8080
	protocol = "HttP"
	check = "http"
	check_passive = true
	check_path = "/"
}
`
}

func testAccCheckLinodeNodeBalancerConfigSSL(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + fmt.Sprintf(`
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8080
	protocol = "https"
	check = "http"
	check_passive = true
	check_path = "/"
	ssl_cert = <<EOT
%s
EOT
	ssl_key = <<EOT
%s
EOT
}
`, testCertifcate, testPrivateKey)
}

func testAccCheckLinodeNodeBalancerConfigProxyProtocol(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 80
	protocol = "tcp"
	proxy_protocol = "v2"
}
`
}

func testAccCheckLinodeNodeBalancerConfigUpdates(nodebalancer string) string {
	return testAccCheckLinodeNodeBalancerBasic(nodebalancer) + `
resource "linode_nodebalancer_config" "foofig" {
	nodebalancer_id = "${linode_nodebalancer.foobar.id}"
	port = 8088
	protocol = "http"
	check = "http"
	check_path = "/foo"
	check_attempts = 3
	check_timeout = 30
	check_passive = false
	stickiness = "http_cookie"
	algorithm = "source"
}
`
}
