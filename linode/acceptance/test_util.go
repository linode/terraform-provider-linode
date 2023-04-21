package acceptance

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode"
	"github.com/linode/terraform-provider-linode/linode/helper"
	"github.com/linode/terraform-provider-linode/version"
	"golang.org/x/crypto/ssh"
)

const (
	optInTestsEnvVar         = "ACC_OPT_IN_TESTS"
	SkipInstanceReadyPollKey = "skip_instance_ready_poll"
)

type AttrValidateFunc func(val string) error

type ListAttrValidateFunc func(resourceName, path string, state *terraform.State) error

var (
	optInTests               map[string]struct{}
	privateKeyMaterial       string
	PublicKeyMaterial        string
	TestAccProviders         map[string]*schema.Provider
	TestAccProvider          *schema.Provider
	TestAccFrameworkProvider *linode.FrameworkProvider
	ConfigTemplates          *template.Template
	TestImageLatest          string
	TestImagePrevious        string
)

func initOptInTests() {
	optInTests = make(map[string]struct{})

	optInTestsValue, ok := os.LookupEnv(optInTestsEnvVar)
	if !ok {
		return
	}

	for _, testName := range strings.Split(optInTestsValue, ",") {
		optInTests[testName] = struct{}{}
	}
}

// initTestImages grabs the latest Linode Alpine images for acceptance test configurations
func initTestImages() {
	client, err := GetClientForSweepers()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	imageFilter := &linodego.Filter{}

	imageFilter.AddField(linodego.Eq, "vendor", "Alpine")

	filterJSON, err := imageFilter.MarshalJSON()
	if err != nil {
		log.Fatalf("failed to create image filter json: %s", err)
	}

	images, err := client.ListImages(context.Background(), &linodego.ListOptions{Filter: string(filterJSON)})
	if err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(images, func(i, j int) bool {
		return images[i].Created.After(*images[j].Created)
	})

	TestImageLatest = images[0].ID
	TestImagePrevious = images[1].ID
}

func init() {
	var err error
	PublicKeyMaterial, privateKeyMaterial, err = acctest.RandSSHKeyPair("linode@ssh-acceptance-test")
	if err != nil {
		log.Fatalf("Failed to generate random SSH key pair for testing: %s", err)
	}

	initOptInTests()

	TestAccProvider = linode.Provider()
	TestAccFrameworkProvider = linode.CreateFrameworkProvider(version.ProviderVersion).(*linode.FrameworkProvider)
	TestAccProviders = map[string]*schema.Provider{
		"linode": TestAccProvider,
	}

	var templateFiles []string

	err = filepath.Walk(
		"../",
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if filepath.Ext(path) == ".gotf" {
				templateFiles = append(templateFiles, path)
			}

			return nil
		},
	)

	if err != nil {
		log.Fatalf("failed to load template files: %v", err)
	}

	ConfigTemplates = template.New("tf-test")
	if _, err := ConfigTemplates.ParseFiles(templateFiles...); err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	initTestImages()
}

func TestProvider(t *testing.T) {
	t.Parallel()

	if err := linode.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func PreCheck(t *testing.T) {
	t.Helper()

	if v := os.Getenv("LINODE_TOKEN"); v == "" {
		t.Fatal("LINODE_TOKEN must be set for acceptance tests")
	}
}

func OptInTest(t *testing.T) {
	t.Helper()

	if _, ok := optInTests[t.Name()]; !ok {
		t.Skipf("skipping opt-in test; specify test in environment variable %q to run", optInTestsEnvVar)
	}
}

func GetSSHClient(t *testing.T, user, addr string) (client *ssh.Client) {
	t.Helper()

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyMaterial))
	if err != nil {
		t.Fatalf("failed to parse private key: %s", err)
	}

	// #nosec G106 -- Test data, not used in production
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Minute,
	}

	attempts := 3

	for attempts != 0 {
		client, err = ssh.Dial("tcp", addr+":22", config)
		if err == nil {
			break
		}

		t.Logf("ssh dial failed: %s", err)
		attempts--

		time.Sleep(5 * time.Second)
	}

	if client == nil {
		t.Fatal("failed to get ssh client")
	}
	return
}

func CheckResourceAttrContains(resName string, path, desiredValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		value, ok := rs.Primary.Attributes[path]
		if !ok {
			return fmt.Errorf("attribute %s does not exist", path)
		}

		if !strings.Contains(value, desiredValue) {
			return fmt.Errorf("value '%s' was not found", desiredValue)
		}

		return nil
	}
}

func ValidateResourceAttr(resName, path string, comparisonFunc AttrValidateFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		value, ok := rs.Primary.Attributes[path]
		if !ok {
			return fmt.Errorf("attribute %s does not exist", path)
		}

		err := comparisonFunc(value)
		if err != nil {
			return fmt.Errorf("comparison failed: %s", err)
		}

		return nil
	}
}

func CheckResourceAttrGreaterThan(resName, path string, target int) resource.TestCheckFunc {
	return ValidateResourceAttr(resName, path, func(val string) error {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return err
		}

		if !(valInt > target) {
			return fmt.Errorf("%d <= %d", valInt, target)
		}

		return nil
	})
}

func CheckResourceAttrNotEqual(resName string, path, notValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		if value, ok := rs.Primary.Attributes[path]; !ok {
			return fmt.Errorf("attribute %s does not exist", path)
		} else if value == notValue {
			return fmt.Errorf("attribute was equal")
		}

		return nil
	}
}

func CheckResourceAttrListContains(resName, path, desiredValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		length, err := strconv.Atoi(rs.Primary.Attributes[path+".#"])
		if err != nil {
			return fmt.Errorf("attribute %s does not exist", path)
		}

		for i := 0; i < length; i++ {
			if rs.Primary.Attributes[path+"."+strconv.Itoa(i)] == desiredValue {
				return nil
			}
		}

		return fmt.Errorf("Desired value not found in resource attribute")
	}
}

func LoopThroughStringList(resName, path string, listValidateFunc ListAttrValidateFunc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resName]
		if !ok {
			return fmt.Errorf("Not found: %s", resName)
		}

		length, err := strconv.Atoi(rs.Primary.Attributes[path+".#"])
		if err != nil {
			return fmt.Errorf("attribute %s does not exist", path)
		}

		for i := 0; i < length; i++ {
			err := listValidateFunc(resName, path+"."+strconv.Itoa(i), s)
			if err != nil {
				return fmt.Errorf("Value not found:%s", err)
			}
		}

		return nil
	}
}

func CheckLKEClusterDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_lke_cluster" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to parse LKE Cluster ID: %s", err)
		}

		if id == 0 {
			return fmt.Errorf("should not have LKE Cluster ID of 0")
		}

		if _, err = client.GetLKECluster(context.Background(), id); err == nil {
			return fmt.Errorf("should not find Linode ID %d existing after delete", id)
		} else if apiErr, ok := err.(*linodego.Error); !ok {
			return fmt.Errorf("expected API Error but got %#v", err)
		} else if apiErr.Code != 404 {
			return fmt.Errorf("expected an error 404 but got %#v", apiErr)
		}
	}

	return nil
}

func CheckVolumeDestroy(s *terraform.State) error {
	client := TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_volume" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)
		}

		_, err = client.GetVolume(context.Background(), id)

		if err == nil {
			return fmt.Errorf("Linode Volume with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("Error requesting Linode Volume with id %d", id)
		}
	}

	return nil
}

func CheckVolumeExists(name string, volume *linodego.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetVolume(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Volume %s: %s", rs.Primary.Attributes["label"], err)
		}

		*volume = *found

		return nil
	}
}

func CheckFirewallExists(provider *schema.Provider, name string, firewall *linodego.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := provider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		found, err := client.GetFirewall(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Error retrieving state of Firewall %s: %s", rs.Primary.Attributes["label"], err)
		}

		*firewall = *found

		return nil
	}
}

func CheckEventAbsent(name string, entityType linodego.EntityType, action linodego.EventAction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*helper.ProviderMeta).Client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing %v to int", rs.Primary.ID)
		}

		event, err := helper.GetLatestEvent(context.Background(), &client, id, entityType, action)
		if err != nil {
			return err
		}

		if event != nil {
			return fmt.Errorf("event exists: %d", event.ID)
		}

		return nil
	}
}

func ExecuteTemplate(t *testing.T, templateName string, data interface{}) string {
	t.Helper()

	var b bytes.Buffer

	err := ConfigTemplates.ExecuteTemplate(&b, templateName, data)
	if err != nil {
		t.Fatalf("failed to execute template %s: %v", templateName, err)
	}

	return b.String()
}

func CreateTempFile(t *testing.T, name, content string) *os.File {
	file, err := ioutil.TempFile(os.TempDir(), name)
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}

	t.Cleanup(func() {
		if err := os.Remove(file.Name()); err != nil {
			t.Fatalf("failed to remove test file: %s", err)
		}
	})

	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write to temp file: %s", err)
	}

	return file
}

func CreateTestProvider() (*schema.Provider, map[string]*schema.Provider) {
	provider := linode.Provider()
	providerMap := map[string]*schema.Provider{
		"linode": provider,
	}
	return provider, providerMap
}

type ProviderMetaModifier func(ctx context.Context, config *helper.ProviderMeta) error

func ModifyProviderMeta(provider *schema.Provider, modifier ProviderMetaModifier) {
	oldConfigure := provider.ConfigureContextFunc

	provider.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config, err := oldConfigure(ctx, data)
		if err != nil {
			return nil, err
		}

		if err := modifier(ctx, config.(*helper.ProviderMeta)); err != nil {
			return nil, diag.FromErr(err)
		}

		return config, nil
	}
}

// GetRegionsWithCaps returns a list of regions that support the given capabilities.
func GetRegionsWithCaps(capabilities []string) ([]string, error) {
	result := make([]string, 0)

	client, err := GetClientForSweepers()
	if err != nil {
		return nil, err
	}

	regions, err := client.ListRegions(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	regionHasCaps := func(r linodego.Region) bool {
		capsMap := make(map[string]bool)

		for _, c := range r.Capabilities {
			capsMap[strings.ToUpper(c)] = true
		}

		for _, c := range capabilities {
			if _, ok := capsMap[strings.ToUpper(c)]; !ok {
				return false
			}
		}

		return true
	}

	for _, region := range regions {
		if region.Status != "ok" || !regionHasCaps(region) {
			continue
		}

		result = append(result, region.ID)
	}

	return result, nil
}

// GetRandomRegionWithCaps gets a random region given a list of region capabilities.
func GetRandomRegionWithCaps(capabilities []string) (string, error) {
	rand.Seed(time.Now().UnixNano())

	regions, err := GetRegionsWithCaps(capabilities)
	if err != nil {
		return "", err
	}

	if len(regions) < 1 {
		return "", fmt.Errorf("no region found with the provided caps")
	}

	// #nosec G404 -- Test data, doesn't need to be cryptography
	return regions[rand.Intn(len(regions))], nil
}

// GetRandomOBJCluster gets a random Object Storage cluster.
func GetRandomOBJCluster() (string, error) {
	rand.Seed(time.Now().UnixNano())

	client, err := GetClientForSweepers()
	if err != nil {
		return "", err
	}

	clusters, err := client.ListObjectStorageClusters(context.Background(), nil)
	if err != nil {
		return "", err
	}

	if len(clusters) < 1 {
		return "", fmt.Errorf("no clusters found")
	}

	// #nosec G404 -- Test data, doesn't need to be cryptography
	return clusters[rand.Intn(len(clusters))].ID, nil
}
