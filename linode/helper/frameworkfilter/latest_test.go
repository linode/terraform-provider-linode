//go:build unit

package frameworkfilter

import (
	"testing"
	"time"

	"github.com/linode/terraform-provider-linode/v2/linode/helper"
)

func TestGetLatestCreated(t *testing.T) {
	type ElemType struct {
		Created *time.Time `json:"-"`
	}

	timeNow := time.Now()
	timeSoon := timeNow.Add(time.Minute)

	elems := []ElemType{
		{
			Created: &timeNow,
		},
		{
			Created: &timeSoon,
		},
	}

	result, d := testFilterConfig.GetLatestCreated(
		helper.TypedSliceToAny(elems),
		"Created",
	)
	if d != nil {
		t.Fatal(d.Detail())
	}

	if !result[0].(ElemType).Created.Equal(*elems[1].Created) {
		t.Fatalf("Expected %s, got %s",
			*elems[1].Created,
			result[0].(ElemType).Created)
	}
}

func TestGetLatestVersion(t *testing.T) {
	type ElemType struct {
		Version string `json:"version"`
	}

	elems := []ElemType{
		{
			Version: "v1.2.3",
		},
		{
			Version: "v1.3.1",
		},
	}

	result, d := testFilterConfig.GetLatestVersion(
		helper.TypedSliceToAny(elems),
		"Version",
	)
	if d != nil {
		t.Fatal(d.Detail())
	}

	if result.(ElemType).Version != elems[1].Version {
		t.Fatalf("Expected %s, got %s",
			elems[1].Version,
			result.(ElemType).Version)
	}
}
