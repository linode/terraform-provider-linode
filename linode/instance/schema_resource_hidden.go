package instance

import (
	"maps"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// hiddenHostIDEnvVar gates the unexported `host_id` attribute.
//
// `host_id` maps to an internal-only Linode API field used to pin an instance
// to a specific host machine. It is required to provision certain specialized
// plans (e.g. RDMA-capable GPU instances) whose plans report as unavailable in
// every region and can only be created on a qualifying host.
//
// Because the underlying API field is not public, the attribute is deliberately
// kept out of the released provider: unless this environment variable is set,
// `host_id` is absent from the resource schema entirely and will not appear in
// `terraform providers schema`, documentation, or configuration validation.
const hiddenHostIDEnvVar = "LINODE_ENABLE_HIDDEN_HOST_ID"

// hostIDSchema is the schema for the gated `host_id` attribute.
var hostIDSchema = &schema.Schema{
	Type: schema.TypeInt,
	Description: "The ID of the host machine to provision this Linode on. " +
		"This field is internal-only and is not supported for general use.",
	Optional: true,
	ForceNew: true,
}

// hiddenHostIDEnabled reports whether the gated `host_id` attribute should be
// registered on the resource schema.
func hiddenHostIDEnabled() bool {
	_, ok := os.LookupEnv(hiddenHostIDEnvVar)
	return ok
}

// resourceSchemaWithHiddenAttributes returns the instance resource schema,
// additionally registering internal-only attributes when they are enabled.
//
// The base schema is copied rather than mutated so that the package-level
// resourceSchema is never modified.
func resourceSchemaWithHiddenAttributes() map[string]*schema.Schema {
	if !hiddenHostIDEnabled() {
		return resourceSchema
	}

	result := make(map[string]*schema.Schema, len(resourceSchema)+1)
	maps.Copy(result, resourceSchema)
	result["host_id"] = hostIDSchema

	return result
}
