//go:build integration || domain

package domain_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/linode/linodego/v2"
	"github.com/linode/terraform-provider-linode/v4/linode"
	"github.com/linode/terraform-provider-linode/v4/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v4/linode/domain/tmpl"
	"github.com/linode/terraform-provider-linode/v4/linode/helper"
	"github.com/stretchr/testify/assert"
)

func init() {
	resource.AddTestSweepers("linode_domain", &resource.Sweeper{
		Name: "linode_domain",
		F:    sweep,
	})
}

func sweep(prefix string) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}

	listOpts := acceptance.SweeperListOptions(prefix, "domain")
	domains, err := client.ListDomains(context.Background(), listOpts)
	if err != nil {
		return fmt.Errorf("Error getting domains: %s", err)
	}
	for _, domain := range domains {
		if !acceptance.ShouldSweep(prefix, domain.Domain) {
			continue
		}
		err := client.DeleteDomain(context.Background(), domain.ID)
		if err != nil {
			return fmt.Errorf("Error destroying %s during sweep: %s", domain.Domain, err)
		}
	}

	return nil
}

func TestSmokeTests_domain_resource(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{"TestAccResourceDomain_basic_smoke", TestAccResourceDomain_basic_smoke},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestAccResourceDomain_basic_smoke(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttrSet(resName, "type"),
					resource.TestCheckResourceAttrSet(resName, "soa_email"),
					resource.TestCheckResourceAttrSet(resName, "description"),
					resource.TestCheckResourceAttrSet(resName, "retry_sec"),
					resource.TestCheckResourceAttrSet(resName, "expire_sec"),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckNoResourceAttr(resName, "master_ips"),
					resource.TestCheckNoResourceAttr(resName, "axfr_ips"),
					resource.TestCheckResourceAttr(resName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
				),
			},

			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_update(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
				),
			},
			{
				Config: tmpl.Updates(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", fmt.Sprintf("renamed-%s", domainName)),
					resource.TestCheckResourceAttr(resName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resName, "tags.0", "tf_test"),
					resource.TestCheckResourceAttr(resName, "tags.1", "tf_test_2"),
				),
			},
		},
	})
}

func TestAccResourceDomain_roundedDomainSecs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.RoundedSec(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "refresh_sec", "3600"),
					resource.TestCheckResourceAttr(resName, "retry_sec", "7200"),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "300"),
					resource.TestCheckResourceAttr(resName, "expire_sec", "2419200"),
				),
			},
			{
				Config:            tmpl.RoundedSec(t, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_zeroSecs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.ZeroSec(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "refresh_sec", "0"),
					resource.TestCheckResourceAttr(resName, "retry_sec", "0"),
					resource.TestCheckResourceAttr(resName, "ttl_sec", "0"),
					resource.TestCheckResourceAttr(resName, "expire_sec", "0"),
				),
			},
			{
				Config:            tmpl.ZeroSec(t, domainName),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDomain_updateIPs(t *testing.T) {
	t.Parallel()

	domainName := acctest.RandomWithPrefix("tf-test") + ".example"
	resName := "linode_domain.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.IPS(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "domain", domainName),
					resource.TestCheckResourceAttr(resName, "master_ips.0", "12.34.56.78"),
					resource.TestCheckResourceAttr(resName, "axfr_ips.0", "87.65.43.21"),
				),
			},
			{
				Config: tmpl.IPSUpdates(t, domainName),
				Check: resource.ComposeTestCheckFunc(
					checkDomainExists,
					resource.TestCheckResourceAttr(resName, "master_ips.#", "0"),
					resource.TestCheckResourceAttr(resName, "axfr_ips.#", "0"),
				),
			},
		},
	})
}

var domainGetPathRegex = regexp.MustCompile(`/domains/\d+$`)

// notFoundInjectionTransport injects a configurable number of synthetic 404
// responses for GET requests to an individual domain, simulating the eventual
// consistency behavior occasionally observed in the API.
type notFoundInjectionTransport struct {
	next      http.RoundTripper
	remaining atomic.Int32
	injected  atomic.Int32
}

func (t *notFoundInjectionTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.Method != http.MethodGet || !domainGetPathRegex.MatchString(request.URL.Path) {
		return t.next.RoundTrip(request)
	}

	if t.remaining.Add(-1) < 0 {
		return t.next.RoundTrip(request)
	}

	t.injected.Add(1)

	body := `{"errors": [{"reason": "Not found"}]}`

	return &http.Response{
		Status:        "404 Not Found",
		StatusCode:    http.StatusNotFound,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header{"Content-Type": []string{"application/json"}},
		Body:          io.NopCloser(strings.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       request,
	}, nil
}

// providerWithNotFoundInjection returns an SDKv2 provider whose client injects
// the given number of 404 responses into GET requests for individual domains.
func providerWithNotFoundInjection(t *testing.T, notFoundCount int32) (*schema.Provider, *notFoundInjectionTransport) {
	t.Helper()

	transport := &notFoundInjectionTransport{}
	transport.remaining.Store(notFoundCount)

	client := acceptance.GetFrameworkTestClient(t, []helper.HTTPClientModifier{
		func(client *http.Client) error {
			transport.next = client.Transport
			client.Transport = transport
			return nil
		},
	})

	overriddenProvider := acceptance.ModifyProviderMeta(
		linode.Provider(),
		func(ctx context.Context, data *schema.ResourceData, config *helper.ProviderMeta) error {
			config.Client = *client
			return nil
		},
	)

	return overriddenProvider, transport
}

// checkDomainExistsWithTestClient is equivalent to checkDomainExists but uses a
// standalone test client, making it safe to use with tests that override the
// provider's client.
func checkDomainExistsWithTestClient(s *terraform.State) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get test client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		if _, err := client.GetDomain(context.Background(), id); err != nil {
			return fmt.Errorf("Error retrieving state of Domain %s: %s", rs.Primary.Attributes["domain"], err)
		}
	}

	return nil
}

// checkDestroyWithTestClient is equivalent to checkDestroy but uses a
// standalone test client, making it safe to use with tests that override the
// provider's client.
func checkDestroyWithTestClient(s *terraform.State) error {
	client, err := acceptance.GetTestClient()
	if err != nil {
		return fmt.Errorf("failed to get test client: %w", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		if _, err := client.GetDomain(context.Background(), id); err == nil {
			return fmt.Errorf("Linode Domain with id %d still exists", id)
		} else if !linodego.IsNotFound(err) {
			return fmt.Errorf("Error requesting Linode Domain with id %d: %w", id, err)
		}
	}

	return nil
}

// TestAccResourceDomain_readNotFoundRetry ensures that transient 404 responses
// when reading a domain are retried rather than causing the domain to be
// dropped from the state.
func TestAccResourceDomain_readNotFoundRetry(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	overriddenProvider, transport := providerWithNotFoundInjection(t, 2)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"linode": func() (tfprotov6.ProviderServer, error) {
				return acceptance.ProtoV6CustomProviderFactories["linode"](nil, overriddenProvider)
			},
		},
		CheckDestroy: checkDestroyWithTestClient,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check:  checkDomainExistsWithTestClient,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName, tfjsonpath.New("domain"), knownvalue.StringExact(domainName),
					),
					statecheck.ExpectKnownValue(
						resName, tfjsonpath.New("type"), knownvalue.StringExact("master"),
					),
					statecheck.ExpectKnownValue(
						resName, tfjsonpath.New("status"), knownvalue.NotNull(),
					),
				},
			},
		},
	})

	// The domain read should have transparently retried past every injected 404.
	assert.EqualValues(t, 2, transport.injected.Load())
	assert.LessOrEqual(t, transport.remaining.Load(), int32(0))
}

// TestAccResourceDomain_readNotFoundDeleted ensures that a domain deleted
// outside of Terraform is still removed from the state once the 404 retry
// deadline has been exceeded.
func TestAccResourceDomain_readNotFoundDeleted(t *testing.T) {
	t.Parallel()

	resName := "linode_domain.foobar"
	domainName := acctest.RandomWithPrefix("tf-test") + ".example"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, domainName),
				Check:  checkDomainExists,
			},
			{
				PreConfig: func() {
					client, err := acceptance.GetTestClient()
					if err != nil {
						t.Fatalf("failed to get test client: %s", err)
					}

					domains, err := client.ListDomains(context.Background(), linodego.NewListOptions(
						0, fmt.Sprintf(`{"domain": %q}`, domainName),
					))
					if err != nil {
						t.Fatalf("failed to list domains: %s", err)
					}

					if len(domains) != 1 {
						t.Fatalf("expected exactly one domain with name %s, got %d", domainName, len(domains))
					}

					if err := client.DeleteDomain(context.Background(), domains[0].ID); err != nil {
						t.Fatalf("failed to delete domain out-of-band: %s", err)
					}
				},
				// The domain should be recreated after the read exhausts its
				// retries and removes the resource from the state.
				Config: tmpl.Basic(t, domainName),
				Check:  checkDomainExists,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resName, tfjsonpath.New("domain"), knownvalue.StringExact(domainName),
					),
				},
			},
		},
	})
}

func checkDomainExists(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetDomain(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Domain %s: %s", rs.Primary.Attributes["domain"], err)
		}
	}

	return nil
}

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccSDKv2Provider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_domain" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetDomain(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Domain with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Domain with id %d", id)
		}
	}

	return nil
}

func configBasic(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "%s"
	type = "master"
	status = "active"
	soa_email = "example@%s"
	description = "tf-testing"
	tags = ["tf_test"]
}`, domain, domain)
}

func configRoundedSec(domain string) string {
	return fmt.Sprintf(`
resource "linode_domain" "foobar" {
	domain = "%s"
	type = "master"
	status = "active"
	soa_email = "example@%[1]s"
	description = "tf-testing"
	ttl_sec = 299
	refresh_sec = 600
	retry_sec = 3601
	expire_sec = 2419201
	tags = ["tf_test"]
}`, domain)
}
