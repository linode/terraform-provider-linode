//go:build integration || vpcsubnets

package linode_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func TestAccProvider_Overrides(t *testing.T) {
	token := getEnvToken(t)

	// Broken test config
	file := createTestConfig(t, fmt.Sprintf(`
[default]
token = 54321
api_url = https://cool.linode.com
api_version = v4reallycoolapiversion

[cool]
token = %s
api_url = https://api.linode.com
api_version = v4beta
`, token))

	config := &helper.Config{
		AccessToken:           token,
		APIURL:                "https://api.linode.com",
		APIVersion:            "v4beta",
		SkipInstanceReadyPoll: false,
		ConfigPath:            file.Name(),
		ObjAccessKey:          "abcd",
		ObjSecretKey:          "efgh",
	}

	client, err := config.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// This request should be successful
	if _, err := client.ListTypes(context.Background(), nil); err != nil {
		t.Fatal(err)
	}

	// Defer to the config file
	config.AccessToken = ""
	config.APIURL = ""
	config.APIVersion = ""

	client, err = config.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// This request should be unsuccessful (using config default profile)
	if _, err := client.ListTypes(context.Background(), nil); err == nil {
		t.Fatal("expected error")
	}

	config.ConfigProfile = "cool"
	client, err = config.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// This request should be successful (using config alternate profile)
	if _, err := client.ListTypes(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
}

func createTestConfig(t *testing.T, conf string) *os.File {
	file, err := os.CreateTemp("", "linode")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Fprint(file, conf)

	t.Cleanup(func() {
		file.Close()
		os.Remove(file.Name())
	})

	return file
}

func getEnvToken(t *testing.T) string {
	token, ok := os.LookupEnv("LINODE_TOKEN")
	if !ok {
		t.Fatal("LINODE_TOKEN must be specified")
	}

	return token
}
