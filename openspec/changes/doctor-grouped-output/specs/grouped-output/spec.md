## ADDED Requirements

### Requirement: Doctor output is organized into grouped sections

`complyctl doctor` SHALL render check results grouped under labeled section headers in the following order: Providers, Policies, Cache, Complypacks, Workspace, Verification. Each section header SHALL be printed on its own line, with check results indented below it.

#### Scenario: All checks pass with grouped output

- **GIVEN** a workspace with healthy providers, cached policies, and valid config
- **WHEN** `complyctl doctor` runs
- **THEN** results SHALL be grouped under section headers (Providers, Policies, Cache, Complypacks, Workspace, Verification)
- **AND** each result SHALL be indented under its section header

### Requirement: Variable checks are nested under their parent provider

Variable check results SHALL be rendered as children of their matching provider result, identified by evaluator-id. Children SHALL display a short label (`variables:`) followed by their message.

#### Scenario: Variables nest under providers

- **GIVEN** providers "ampel" and "opa" are healthy
- **AND** variable checks produce results for evaluators "ampel" and "opa"
- **WHEN** `complyctl doctor` renders output
- **THEN** each variable result SHALL appear indented under its matching provider
- **AND** the variable result SHALL display as `variables: <message>` (not the full qualified name)

#### Scenario: Unmapped providers with no required variables skip variable output

- **GIVEN** provider "opa" is healthy
- **AND** no policies in the workspace reference evaluator-id "opa"
- **AND** the opa provider declares no required global or target variables
- **WHEN** `complyctl doctor` renders output
- **THEN** the "opa" provider result SHALL have no variable child
- **AND** the provider SHALL render as a single line without nested sub-results

#### Scenario: Verbose variable detail nests under variable summary

- **GIVEN** `--verbose` flag is set
- **AND** provider "ampel" has target variables "url" and "specs"
- **WHEN** `complyctl doctor` renders output
- **THEN** per-key detail results SHALL appear indented under the variable summary result (not as siblings)
- **AND** each detail result SHALL display as `<status-emoji> <message>` (message-only, no label)

#### Scenario: Variable resolve warnings appear as top-level in Providers section

- **GIVEN** a variable resolve failure for a policy reference
- **WHEN** `complyctl doctor` renders output
- **THEN** the resolve warning SHALL appear as a top-level result in the Providers section (not nested)
- **AND** the resolve warning SHALL have `Group: GroupProviders` (not `GroupVariables`)

### Requirement: Active-period checks are nested under their parent policy

Active-period check results SHALL be rendered as children of their matching policy version result, identified by policy-id. Verbose detail results SHALL be rendered as children of the active-period result.

#### Scenario: Active-period nests under policy

- **GIVEN** a policy "cis" with a cached version and an evaluation timeline
- **WHEN** `complyctl doctor` renders output
- **THEN** the active-period result SHALL appear indented under the `policy/cis` version result
- **AND** the active-period result SHALL display as `active-period: <message>`

#### Scenario: Verbose detail nests under active-period

- **GIVEN** `--verbose` flag is set
- **AND** a policy has evaluation notes and enforcement timeline
- **WHEN** `complyctl doctor` renders output
- **THEN** enforcement and notes results SHALL appear indented under the active-period result

### Requirement: Per-policy results on registry failure

When a registry is unreachable, `CheckPolicyVersions` SHALL emit a warning result for each affected policy (`policy/<eid>: version check skipped (registry unreachable: <registry>)`) instead of a single `registry/<registry>` result. No policy SHALL be silently skipped.

#### Scenario: Unreachable registry produces per-policy warnings

- **GIVEN** two policies from registry "down.io"
- **WHEN** "down.io" is unreachable
- **THEN** each policy SHALL receive its own warning result with the message indicating the registry is unreachable
- **AND** no `registry/` named result SHALL be produced

#### Scenario: Active-period nests under registry-failed policy

- **GIVEN** a policy whose version check was skipped due to unreachable registry
- **AND** the policy's active-period is resolvable from local cache
- **WHEN** `complyctl doctor` renders output
- **THEN** the active-period result SHALL nest under the policy warning result

### Requirement: Complypacks are listed individually by evaluator-id

Each cached complypack SHALL be listed as an individual result (`complypacks/<evaluator-id>: cached`) instead of a single summary count.

#### Scenario: Multiple complypacks listed individually

- **GIVEN** complypacks cached for evaluators "ampel" and "opa"
- **WHEN** `complyctl doctor` renders output
- **THEN** two results SHALL appear in the Complypacks section: `complypacks/ampel: cached` and `complypacks/opa: cached`

### Requirement: Verbose detail Status reflects actual check state

Verbose detail `CheckResult`s SHALL use `StatusPass` when the variable is present and `StatusFail` when missing. The inline emoji (`✅`/`❌`) SHALL NOT appear in the detail message — the renderer provides the status emoji based on the `Status` field.

#### Scenario: Missing variable detail shows fail status

- **GIVEN** `--verbose` flag is set
- **AND** target "test-k8s-deployment" is missing required variable "url"
- **WHEN** `complyctl doctor` renders the detail line
- **THEN** the detail result SHALL have `Status: StatusFail`
- **AND** the rendered line SHALL show the fail emoji from the renderer (not an inline emoji in the message)

#### Scenario: Present variable detail shows pass status

- **GIVEN** `--verbose` flag is set
- **AND** target "complyctl-repo" has required variable "url" configured
- **WHEN** `complyctl doctor` renders the detail line
- **THEN** the detail result SHALL have `Status: StatusPass`
- **AND** the rendered line SHALL show the pass emoji from the renderer

### Requirement: Summary counts top-level and children, excludes grandchildren

The summary line (`N checks: X passed, Y failed, Z warnings`) SHALL count top-level results (depth 0) and their direct children (depth 1). Grandchildren (depth 2 — verbose detail) SHALL NOT be counted. This ensures the summary is stable regardless of `--verbose` while still reflecting child-level failures.

#### Scenario: Summary includes children

- **GIVEN** 3 provider results (top-level) with 3 variable children and 2 policy results with 2 active-period children
- **WHEN** `complyctl doctor` renders the summary
- **THEN** the summary SHALL report 10 checks (top-level + children)

#### Scenario: Summary excludes grandchildren in verbose mode

- **GIVEN** `--verbose` flag is set
- **AND** 3 provider results with 3 variable children, each having 2 verbose detail grandchildren
- **WHEN** `complyctl doctor` renders the summary
- **THEN** the summary SHALL report the same count as without `--verbose` (grandchildren are not counted)
- **AND** the 6 grandchild detail lines SHALL be rendered but excluded from the total

### Requirement: Renderer uses explicit Label field for display

The renderer SHALL display `Label` when set, falling back to `Name` when `Label` is empty. The renderer SHALL NOT parse `Name` to derive display text. The `displayName`, `childLabel` helpers and the `strings` import SHALL be removed from the renderer.

#### Scenario: Top-level provider displays Label

- **GIVEN** a provider result with `Name: "provider/ampel"` and `Label: "ampel"`
- **WHEN** `complyctl doctor` renders the result
- **THEN** the displayed name SHALL be `ampel` (from `Label`)
- **AND** no string splitting of `Name` SHALL occur

#### Scenario: Child variable displays Label

- **GIVEN** a variable child result with `Name: "variables/ampel"` and `Label: "variables"`
- **WHEN** `complyctl doctor` renders the child
- **THEN** the displayed label SHALL be `variables` (from `Label`)

#### Scenario: Singleton result falls back to Name

- **GIVEN** a cache result with `Name: "cache"` and `Label: ""`
- **WHEN** `complyctl doctor` renders the result
- **THEN** the displayed name SHALL be `cache` (fallback to `Name`)

## MODIFIED Requirements

- `CheckResult` struct gains `Group CheckGroup`, `Children []CheckResult`, and `Label string` fields.
- `CheckPolicyVersions` changes from emitting `registry/<reg>` results to per-policy `policy/<eid>` warnings for unreachable registries.
- `CheckVariables` produces verbose detail results as `Children` of the summary `CheckResult` (not as separate flat results). Detail `Status` reflects actual variable presence.
- `CheckPolicyActivePeriod` produces verbose detail results as `Children` of the active-period `CheckResult` (not as separate flat results).

## REMOVED Requirements

None.
