## Why

`complyctl doctor` renders all check results as a flat list of status lines. As the number of checks grows (providers, policies, variables, complypacks, verification), the output becomes harder to scan. Related checks are visually separated: variables appear far from their providers, active-period appears far from its parent policy, and complypacks show only a count without naming them. When a registry is unreachable, one `registry/` warning silently skips all affected policies with no per-policy indication.

This is tracked in [issue #692](https://github.com/complytime/complyctl/issues/692). Feedback from @marcusburghardt confirmed the need for grouped, sequentially-ordered output with named complypacks.

## What Changes

- Check results gain a `Children []CheckResult` field for explicit parent-child nesting and a `Group` field for section assignment.
- `Run()` assembles a result tree: attaches variable results to their provider and active-period results to their policy by evaluator/policy-id. Unmatched children are promoted to top-level (never silently dropped).
- `CheckPolicyVersions` emits per-policy warning results when a registry is unreachable, instead of one `registry/` result that silently skips subsequent policies.
- `printDiagnostics` walks the tree with increasing indent, rendering `Name: Message` at every level. Summary counts only top-level results.
- Output sections are reordered to follow scan prerequisite sequence: Providers, Policies, Cache, Complypacks, Workspace, Verification.

## Capabilities

### New Capabilities
- `grouped-output`: Doctor output is organized into labeled sections with nested sub-results, following the logical sequence of scan prerequisites.

### Modified Capabilities
- `doctor`: Output format changes from flat list to grouped, nested tree. `CheckResult` struct gains `Group` and `Children` fields. `CheckPolicyVersions` emits per-policy warnings on registry failure.

### Removed Capabilities
- None. This is a presentation-layer redesign with one behavioral change (per-policy registry warnings).

## Impact

- `internal/doctor/doctor.go` — `CheckResult` gains `Group` and `Children` fields; `CheckGroup` type and constants added; `Run()` gains tree assembly helpers; `CheckPolicyVersions` emits per-policy results for unreachable registries
- `cmd/complyctl/cli/doctor.go` — `printDiagnostics` rewritten as recursive tree walk; string-parsing helpers removed
- `internal/doctor/doctor_test.go` — tests for tree assembly, per-policy registry warnings, and updated complypack assertions

## Constitution Alignment

### I. Autonomous Collaboration

**Assessment**: PASS

Grouped output with nested sub-results is self-describing. Users see which providers have which variable checks, which policies have which timeline checks, and which policies were affected by unreachable registries — all without consulting documentation.

### II. Composability First

**Assessment**: PASS

The `Children` field is a general-purpose mechanism. Any `Check*` function can produce children without the renderer needing domain-specific knowledge. The tree assembly in `Run()` is the single wiring point, keeping concerns separated.

### III. Observable Quality

**Assessment**: PASS

Per-policy registry warnings replace silent skipping. Summary counts are accurate (top-level only). The nested layout makes it visually clear which sub-checks belong to which parent.

### IV. Testability

**Assessment**: PASS

Tree assembly helpers (`attachByEvaluatorID`, `attachByPolicyID`) are independently testable. Existing `Check*` tests are unaffected since they test individual check functions, not assembly. Per-policy registry warning behavior is testable via existing mock infrastructure.

## Documentation Impact

- **CHANGELOG.md**: Entry under `### Changed` for `doctor` command output redesign.
- **AGENTS.md**: Update Recent Changes section with `doctor-grouped-output` entry.
- **Website issue**: Exempt -- presentation-layer change to existing diagnostic command, not a new user-facing feature.
