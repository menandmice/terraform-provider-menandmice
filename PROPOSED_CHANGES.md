# Proposed Changes for terraform-provider-menandmice

> **Status: IMPLEMENTED** - See commit for changes to the following files:
> - `menandmice/resource_range.go` - Added `register_only` and `range_identifier` attributes
> - `menandmice/client_range.go` - Added `RegisterRange()` method
> - `menandmice/data_source_range.go` - Added `range_identifier` attribute
> - `menandmice/data_source_ranges.go` - Added `range_identifier` attribute
> - `docs/resources/range.md` - Updated documentation
> - `examples/resources/menandmice_range/resource.tf` - Added example

## Context

This document summarizes issues identified with the `menandmice_range` resource and proposes solutions based on feedback from a large organization using Terraform for infrastructure automation.

---

## Issue 1: No "Register Only" Mode for Existing Ranges

### Problem Description

In distributed automation environments, infrastructure (including subnets) is created through Terraform pipelines directly in cloud providers. Micetro does not have write access to all subscriptions. The current provider forces Micetro to **create** ranges, which fails when:

- The range already exists in the cloud provider
- Micetro lacks privileges to create resources in that subscription

This leads to:
- Subnets created by Terraform not being tracked in Micetro
- Address space depletion because ranges aren't freed on deletion
- Orphaned subnet entries when parent containers are deleted

### Root Cause

The `CreateRange` function in `menandmice/client_range.go` always performs a `POST /Ranges` call, which instructs Micetro to create the range:

```go
func (c *Mmclient) CreateRange(iprange Range, discovery Discovery) (string, error) {
    // ...
    err := c.Post(postcreate, &re, "Ranges")  // Always tries to CREATE
    // ...
}
```

There is no option to simply **register/claim** an existing range for IPAM tracking purposes.

### Proposed Solution

Add a `register_only` attribute to the `menandmice_range` resource.

#### 1. Schema Change in `menandmice/resource_range.go`

```go
"register_only": {
    Type:        schema.TypeBool,
    Description: "Register an existing range in Micetro without attempting to create it. Use when the range already exists in the cloud provider and you only want to track it in Micetro for IPAM purposes.",
    Optional:    true,
    Default:     false,
    ForceNew:    true,
},
```

#### 2. Client Change in `menandmice/client_range.go`

Add a new method for registering existing ranges (requires corresponding Micetro API support):

```go
func (c *Mmclient) RegisterRange(iprange Range) (string, error) {
    // API call to register/claim existing range without creation
    // The exact implementation depends on Micetro API capabilities
}
```

#### 3. Create Logic Modification in `resourceRangeCreate`

```go
func resourceRangeCreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Mmclient)
    
    registerOnly := d.Get("register_only").(bool)
    
    // ... existing free_range logic ...
    
    iprange := readRangeSchema(d)
    
    var objRef string
    var err error
    
    if registerOnly {
        objRef, err = client.RegisterRange(iprange)
    } else {
        discovery := Discovery{Enabled: false}
        objRef, err = client.CreateRange(iprange, discovery)
    }
    
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(objRef)
    return resourceRangeRead(c, d, m)
}
```

#### Open Questions

- Does the Micetro API have a dedicated endpoint for registering/claiming existing ranges?
- If not, what API modifications would be required on the Micetro side?

---

## Issue 2: Attribute Naming Confusion (`name` vs `title`)

### Problem Description

The `name` attribute returns the CIDR of the range (e.g., `192.168.5.0/24`), not a human-readable title. This is counter-intuitive since `name` typically suggests a human-readable label.

### Root Cause

The naming follows the Micetro API convention where:
- `name` = the range identifier (CIDR or from-to format)
- `title` = human-readable label (stored as a custom property)

Current behavior in `menandmice/resource_range.go`:

```go
func flattenRange(iprange Range, tz *time.Location) (map[string]interface{}, diag.Diagnostics) {
    // ...
    m["name"] = iprange.Name  // This is "192.168.5.0/24" (the CIDR!)
    m["title"] = iprange.RangeProperties.CustomProperties["Title"]  // Human-readable
    // ...
}
```

### Current Attribute Mapping

| Attribute | What it Contains | Source |
|-----------|------------------|--------|
| `name` | CIDR or "from-to" format | API's `Name` field |
| `cidr` | CIDR only (if applicable) | Computed from `Name` |
| `title` | Human-readable name | Custom property |

### Proposed Solution

**Option A: Documentation Clarification (Non-breaking)**

Update documentation to clearly explain the `name` attribute contains the range identifier (CIDR), not a human-readable title.

**Option B: Add New Attribute (Non-breaking)**

Add a `range_identifier` computed attribute that explicitly holds the CIDR/from-to value:

```go
"range_identifier": {
    Type:        schema.TypeString,
    Description: "The range identifier in CIDR notation or from-to format.",
    Computed:    true,
},
```

And update `flattenRange`:

```go
m["range_identifier"] = iprange.Name
```

**Option C: Deprecate and Rename (Breaking)**

1. Deprecate `name` attribute
2. Introduce `range_identifier` as the replacement
3. Remove `name` in a future major version

### Recommendation

Option B is recommended as it maintains backward compatibility while providing clearer semantics.

---

## Issue 3: Orphaned Subnets on Deletion

### Problem Description

When Terraform deletes a subscription:
1. The parent container gets deleted in Micetro
2. Child subnets remain in Micetro
3. Micetro's discovery scan does NOT delete the subnet entries
4. Address space becomes depleted without ever being cleaned

### Root Cause

This is a consequence of Issue 1. Since ranges aren't properly registered via Terraform:
- Terraform doesn't track subnets in its state
- When parent container is deleted, Terraform has no knowledge of child ranges
- Micetro's discovery preserves entries that should be deleted

### Proposed Solution

This issue is resolved by implementing Issue 1's `register_only` mode:

1. Subnets are registered in Micetro on `terraform apply`
2. Terraform tracks them in state
3. When `terraform destroy` runs, subnets are properly deleted from Micetro
4. Address space is correctly freed

---

## Summary

| Issue | Impact | Solution | Breaking Change |
|-------|--------|----------|-----------------|
| No register-only mode | High - address space depletion | Add `register_only` attribute | No |
| `name` confusion | Low - documentation issue | Add `range_identifier` attribute | No |
| Orphaned subnets | High - address space depletion | Resolved by register-only mode | No |

---

## Implementation Priority

1. **High**: `register_only` mode (resolves Issues 1 and 3)
2. **Low**: Attribute naming clarification (resolves Issue 2)

---

## Next Steps

1. Confirm Micetro API capabilities for registering existing ranges
2. Implement `register_only` attribute and corresponding client method
3. Add `range_identifier` attribute for clarity
4. Update documentation
5. Add acceptance tests for new functionality
