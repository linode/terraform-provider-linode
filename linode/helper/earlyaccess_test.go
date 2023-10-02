//go:build unit

package helper_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

func TestAttemptWarnEarlyAccessSDKv2(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	helper.AttemptWarnEarlyAccessSDKv2(&helper.ProviderMeta{
		Config: &helper.Config{
			APIVersion: "v4",
		},
	})

	if !strings.Contains(buf.String(), "[WARN] This resource is in early access but the provider "+
		"API version is set to \"v4\" (expected \"v4beta\").") {
		t.Fatalf("expected warning")
	}

	buf.Reset()

	helper.AttemptWarnEarlyAccessSDKv2(&helper.ProviderMeta{
		Config: &helper.Config{
			APIVersion: "v4beta",
		},
	})

	if buf.Len() != 0 {
		t.Fatal("expected none, got warning")
	}
}

func TestAttemptWarnEarlyAccessFramework(t *testing.T) {
	d := helper.AttemptWarnEarlyAccessFramework(&helper.FrameworkProviderModel{
		APIVersion: types.StringValue("v4"),
	})

	if d.WarningsCount() < 1 {
		t.Fatal("expected warning, got none")
	}

	warnings := d.Warnings()

	if warnings[0].Summary() != "Non-Beta Target API Version" {
		t.Fatal("Invalid warning summary")
	}

	if warnings[0].Detail() != "This resource is in early access but the provider "+
		"API version is set to \"v4\" (expected \"v4beta\")." {
		t.Fatal("Invalid warning detail")
	}

	d = helper.AttemptWarnEarlyAccessFramework(&helper.FrameworkProviderModel{
		APIVersion: types.StringValue("v4beta"),
	})

	if d.WarningsCount() > 0 {
		t.Fatalf("expected no warnings, got %d", d.WarningsCount())
	}
}
