package diffsuppressfuncs

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// AutoAllocRange is a DiffSuppressFunc for fields that either accept a CIDR, "auto",
// or a / followed by a prefix length.
//
// This prevents values of "auto" from causing repeating diffs on subsequent applies.
func AutoAllocRange(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue != "" && newValue == "auto" {
		return true
	}

	addr, prefix, err := helper.ParseRangeOptionalAddress(oldValue)
	if err != nil {
		log.Printf("Failed to parse old range: %s", err)
		return false
	}

	newAddr, newPrefix, err := helper.ParseRangeOptionalAddress(newValue)
	if err != nil {
		log.Printf("Failed to parse new range: %s", err)
		return false
	}

	// One of the addresses only has a prefix specified,
	// so we should only diff on the prefix.
	if (addr == nil) || (newAddr == nil) {
		return prefix == newPrefix
	}

	return prefix == newPrefix && addr.Compare(*newAddr) == 0
}
