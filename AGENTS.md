# Terraform Provider for Linode - AI Agent Instructions

## Architecture Overview

This is a Terraform provider that uses **both** SDKv2 and Plugin Framework patterns (muxed together). Resources/data sources live under `linode/<resource-name>/` with each package being self-contained.

- **SDKv2 resources** (legacy): `linode/instance/`, `linode/domain/`, `linode/lke/` - use `resource.go`, `datasource.go`
- **Plugin Framework resources** (preferred for new work): `linode/vpc/`, `linode/volume/`, `linode/vpcsubnet/` - use `framework_resource.go`, `framework_datasource.go`, `framework_models.go`
- **Provider registration**: SDKv2 in `linode/provider.go`, Framework in `linode/framework_provider.go`
- **Shared utilities**: `linode/helper/` - conversion functions, base resource/datasource, plan modifiers

## Framework Resource Structure (New Resources)

Each Framework resource package follows this pattern:

```
linode/<resource-name>/
├── framework_resource.go          # CRUD operations
├── framework_datasource.go        # Data source Read
├── framework_models.go            # Terraform state models with Flatten/CopyFrom methods
├── framework_schema_resource.go   # Resource schema definition
├── framework_schema_datasource.go # Data source schema definition
├── resource_test.go               # Integration tests
├── datasource_test.go             # Data source integration tests
├── tmpl/                          # Test templates
│   ├── template.go                # Go functions returning HCL configs
│   └── *.gotf                     # HCL template files
```

### Key Model Patterns

Models must implement:
- `FlattenXxx(ctx, apiObject, preserveKnown)` - Converts Linode API response to Terraform state
- `CopyFrom(ctx, other, preserveKnown)` - Copies values between model instances for updates

Use `helper.KeepOrUpdateValue()`, `helper.KeepOrUpdateString()`, `helper.KeepOrUpdateInt64()` to handle `preserveKnown` flag which prevents overwriting known plan values with computed values.

## Test Commands

```bash
# Run all tests for a package
make TEST_SUITE="vpcsubnet" test-int

# Run a specific test
make PKG_NAME="volume" TEST_CASE="TestAccResourceVolume_basic" test-int

# Run unit tests only
make test-unit

# Run unit tests for a specific package
make PKG_NAME="instance" test-unit
```

**Important**: Set `LINODE_TOKEN` environment variable or use `.env` file. Tests create real resources (costs money).

## Test Template Pattern

Tests use `.gotf` template files with Go text/template syntax:

```go
// In tmpl/template.go
func Basic(t testing.TB, label, region string) string {
    return acceptance.ExecuteTemplate(t, "resource_basic", TemplateData{
        Label: label, Region: region,
    })
}
```

```hcl
// In tmpl/basic.gotf
{{ define "resource_basic" }}
resource "linode_example" "foobar" {
    label = "{{.Label}}"
    region = "{{.Region}}"
}
{{ end }}
```

## Build Tags

- `//go:build integration` or `//go:build <resource-name>` - Integration tests (require API token)
- `//go:build unit` - Unit tests (no API calls)

## Helper Functions Reference

| Function | Purpose |
|----------|---------|
| `helper.NewBaseResource()` | Creates base Framework resource with common config |
| `helper.KeepOrUpdateValue()` | Conditionally preserves known values during refresh |
| `helper.FrameworkSafeInt64ToInt()` | Safe int64→int conversion with diagnostics |
| `helper.MapSlice()` | Transforms slices with a mapping function |
| `acceptance.GetRandomRegionWithCaps()` | Gets random region with required capabilities |
| `acceptance.ExecuteTemplate()` | Renders HCL test templates |

## Common Workflows

**Adding a new Framework resource:**
1. Create package under `linode/<resource-name>/`
2. Define schema in `framework_schema_resource.go`
3. Define models in `framework_models.go` with `Flatten*` and `CopyFrom` methods
4. Implement CRUD in `framework_resource.go`
5. Register in `linode/framework_provider.go` Resources() method
6. Add tests with `tmpl/` directory
7. Add docs in `docs/resources/<resource>.md`

**Debugging tests:**
- `TF_LOG_PROVIDER=DEBUG` - Provider logging
- `TF_LOG_PROVIDER_LINODE_REQUESTS=DEBUG` - API request logging
- `TF_SCHEMA_PANIC_ON_ERROR=1` - Force panic on schema errors

## Linode API Client

Uses `github.com/linode/linodego` client. Access via:
- SDKv2: `meta.(*helper.ProviderMeta).Client`
- Framework: `r.Meta.Client`

## Code Style

- Use `golangci-lint fmt` for code formatting (run `make format`), or `gofmt -w` if unavailable
- Use `tflog.Debug(ctx, ...)` for logging in resources
- Prefer Framework over SDKv2 for new resources
- Unit test files use `*_unit_test.go` naming with `//go:build unit` tag

## Go Idioms

**Sets using maps (Go 1.23+):**
Use `helper.StringSet` and `helper.ExistsInSet` for set operations, then extract keys with `slices.Collect(maps.Keys())`:

```go
import (
    "maps"
    "slices"
    "github.com/linode/terraform-provider-linode/v3/linode/helper"
)

// Create a set
regionSet := make(helper.StringSet)
for _, endpoint := range endpoints {
    regionSet[endpoint.Region] = helper.ExistsInSet
}

// Extract keys as a slice (Go 1.23+)
regions := slices.Collect(maps.Keys(regionSet))
```

This is preferred over the manual loop pattern:
```go
// Avoid this verbose pattern
regions := make([]string, 0, len(regionSet))
for region := range regionSet {
    regions = append(regions, region)
}
```
