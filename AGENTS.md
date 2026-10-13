# Terraform Provider for Linode - AI Agent Instructions

## Architecture Overview

This is a Terraform provider that uses **both** SDKv2 and Plugin Framework patterns (muxed together). Resources and data sources generally live under `linode/<resource-name>/` and use shared code from `linode/helper/` and, where needed, other resource packages.

- **SDKv2 resources** (legacy): for example, `linode/instance/` uses `resource.go` and `datasource.go`
- **Mixed packages**: `linode/domain/` and `linode/lke/` use an SDKv2 `resource.go` with a Framework data source
- **Plugin Framework resources** (preferred for new work): `linode/vpc/`, `linode/volume/`, `linode/vpcsubnet/` - use `framework_resource.go`, `framework_datasource.go`, `framework_models.go`
- **Provider registration**: SDKv2 in `linode/provider.go`, Framework in `linode/framework_provider.go`
- **Shared utilities**: `linode/helper/` - conversion functions, base resource/datasource, plan modifiers

## Framework Resource Structure (New Resources)

A full Framework resource package commonly follows this pattern, although smaller packages may omit files. New Plugin Framework tests must use `framework_resource_test.go` and `framework_datasource_test.go`; `resource_test.go` and `datasource_test.go` are legacy SDKv2 filenames.

```
linode/<resource-name>/
├── framework_resource.go          # CRUD operations
├── framework_datasource.go        # Data source Read
├── framework_models.go            # Terraform state models with Flatten/CopyFrom methods
├── framework_schema_resource.go   # Resource schema definition
├── framework_schema_datasource.go # Data source schema definition
├── framework_resource_test.go     # Resource integration tests
├── framework_datasource_test.go   # Data source integration tests
├── tmpl/                          # Test templates
│   ├── template.go                # Go functions returning HCL configs
│   └── *.gotf                     # HCL template files
```

### Key Model Patterns

Resource models commonly provide:
- A `FlattenXxx` or `ParseXxx` method that converts a Linode API response to Terraform state
- A `CopyFrom` method for preserving or copying values during updates when the resource workflow needs one

Signatures vary by package; some methods accept a context or diagnostics parameter, return diagnostics, or omit context when it is unnecessary. Follow the established pattern in the package being changed.

Use `helper.KeepOrUpdateValue()`, `helper.KeepOrUpdateString()`, `helper.KeepOrUpdateInt64()` to handle `preserveKnown` flag which prevents overwriting known plan values with computed values.

## Test Commands

```bash
# Run an integration test suite selected by build tag
make TEST_SUITE="vpcsubnet" test-int

# Run a specific test
make PKG_NAME="volume" TEST_CASE="TestAccResourceVolume_basic" test-int

# Run unit tests only
make test-unit

# Run unit tests for a specific package
make PKG_NAME="instance" test-unit
```

**Important**: Export the `LINODE_TOKEN` environment variable before running acceptance tests. The repository does not automatically load a `.env` file. Acceptance tests create real resources and may incur costs.

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

- `//go:build integration || <resource-name>` - Typical integration-test tag expression (requires an API token); some specialized suites add other tags or conditions
- `//go:build unit` - Unit tests (no API calls)

## Integration Test Style

When writing integration tests, **prefer `ConfigStateChecks` with `statecheck.*` functions** over the legacy `Check` field with `resource.TestCheckResourceAttr` / `resource.ComposeTestCheckFunc`. The `statecheck` API is the modern Terraform testing approach and should be used for all new tests.

**Preferred — `ConfigStateChecks` with `statecheck`:**
```go
Steps: []resource.TestStep{
    {
        Config: tmpl.Basic(t),
        ConfigStateChecks: []statecheck.StateCheck{
            statecheck.ExpectKnownValue(resourceName,
                tfjsonpath.New("label"), knownvalue.StringExact("my-label")),
            statecheck.ExpectKnownValue(resourceName,
                tfjsonpath.New("status"), knownvalue.StringExact("ready")),
            statecheck.ExpectKnownValue(resourceName,
                tfjsonpath.New("backups_enabled"), knownvalue.NotNull()),
        },
    },
}
```

**Legacy — `Check` with `resource.TestCheckResourceAttr` (avoid in new tests):**
```go
Steps: []resource.TestStep{
    {
        Config: tmpl.Basic(t),
        Check: resource.ComposeTestCheckFunc(
            resource.TestCheckResourceAttr(resourceName, "label", "my-label"),
            resource.TestCheckResourceAttrSet(resourceName, "status"),
        ),
    },
}
```

Common `statecheck` / `knownvalue` helpers:
| Function | Purpose |
|----------|---------|
| `statecheck.ExpectKnownValue(name, path, check)` | Assert an attribute matches a known value |
| `knownvalue.StringExact(v)` | Exact string match |
| `knownvalue.Bool(v)` | Exact bool match |
| `knownvalue.NotNull()` | Attribute is set (non-null) |
| `knownvalue.Int64Exact(v)` | Exact int64 match |
| `tfjsonpath.New("attr")` | Root-level attribute path |
| `tfjsonpath.New("list").AtSliceIndex(0).AtMapKey("key")` | Nested path |

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
6. Add `framework_resource_test.go` and/or `framework_datasource_test.go` with a `tmpl/` directory
7. Add docs in `docs/resources/<resource>.md`

**Debugging tests:**
- `TF_LOG_PROVIDER=DEBUG` - Provider logging
- `TF_LOG_PROVIDER_LINODE_REQUESTS=DEBUG` - API request logging

## Linode API Client

Uses `github.com/linode/linodego/v2` client. Access via:
- SDKv2: `meta.(*helper.ProviderMeta).Client`
- Framework: `r.Meta.Client`

## Code Style

- Use `golangci-lint fmt` for code formatting (run `make format`), or `gofmt -w` if unavailable
- Use `tflog.Debug(ctx, ...)` for logging in resources
- Prefer Framework over SDKv2 for new resources
- Unit test files use `*_unit_test.go` naming with `//go:build unit` tag
- Keep this `AGENTS.md` synchronized with the repository: update it when a change contradicts its guidance, and correct any factual inaccuracies you discover while working.

## Go Idioms

**Sets (Go 1.23+):**
Use `github.com/hashicorp/go-set/v3` for set operations:

```go
import (
    "github.com/hashicorp/go-set/v3"
)

// Create a set
regionSet := set.New[string](len(endpoints))
for _, endpoint := range endpoints {
    regionSet.Insert(endpoint.Region)
}

// Extract elements as a slice
regions := regionSet.Slice()
```
