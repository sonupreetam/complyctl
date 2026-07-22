## Context

`complyctl doctor` outputs a flat list of check results. As the number of checks grows, the output becomes hard to scan. Related checks (variables with providers, active-period with policies) are visually separated. Complypacks show only a count. Unreachable registries silently skip affected policies.

Issue #692 tracks the redesign. Feedback from @marcusburghardt established the section ordering and the need for named complypacks.

## Goals / Non-Goals

### Goals
- Group check results into labeled sections following scan prerequisite order
- Nest variable results under their parent provider
- Nest active-period and verbose detail under their parent policy
- List complypacks individually by evaluator-id
- Emit per-policy warnings when a registry is unreachable (no silent skipping)
- Count top-level and children (depth 0–1) in the summary line; exclude grandchildren (depth 2)
- Nest verbose detail under its parent summary (variables detail under variables, active-period detail under active-period)
- Preserve `--verbose` behavior within the nested layout

### Non-Goals
- Machine-readable output format (#611)
- Changing check logic or thresholds
- Adding new diagnostic checks

## Decisions

### D1: Explicit `Children` field instead of string-based name inference

**Decision**: Add `Children []CheckResult` to the `CheckResult` struct. `Run()` assembles the tree by matching results by evaluator/policy-id.

**Rationale**: The initial implementation inferred parent-child relationships by parsing `/`-delimited name segments. This required special-case helpers (`partitionVariables`, `buildChildMap`, `filterParents`, `isChildPolicyResult`) and orphan handling. The explicit field eliminates all string parsing in the renderer, making it a simple tree walk.

**Alternatives considered**: Derive groups from name prefixes at render time (fragile, requires orphan handling); separate render-model struct (over-engineered for the scope).

### D2: Tree assembly in `Run()` as the single wiring point

**Decision**: `Run()` calls individual `Check*` functions (which return flat results), then attaches children to parents via `attachByEvaluatorID` and `attachByPolicyID` helpers. Unmatched children are promoted to top-level in their group.

**Rationale**: Keeps `Check*` functions independently testable (they don't need to know about the tree structure). The wiring logic is centralized and testable in isolation.

### D3: Per-policy results on registry failure

**Decision**: `CheckPolicyVersions` emits a `policy/<eid>: version check skipped (registry unreachable: <reg>)` warning for each affected policy instead of one `registry/<reg>` result that silently skips subsequent policies.

**Rationale**: The current behavior silently drops policies from unreachable registries, and `CheckPolicyActivePeriod` (which uses local cache, not registry) can still produce results for those policies — creating orphan children with no parent. Per-policy results ensure every configured policy has a result, and active-period always has a parent to nest under.

**Implementation note**: `resolvePinnedFallback` already receives the `eid` parameter — the fix is changing the final fallback return from `registry/<ref.Registry>` to `policy/<eid>` and removing the `unreachable` map + `continue` skip in the caller.

### D4: Children display short labels

**Decision**: Nested children use short descriptive `Name` values (`variables`, `active-period`, `evaluation notes`) since the parent already provides context. The renderer displays `Name: Message` at every level uniformly.

**Rationale**: Full qualified paths (`variables/ampel`, `policy/cis/active-period`) are redundant when nested under the parent. Short labels are more readable and the renderer needs no special handling for depth.

### D5: Summary counts top-level and children, excludes grandchildren

**Decision**: The summary line (`N checks: X passed, Y failed, Z warnings`) counts top-level results and their direct children (depth 0 and 1). Grandchildren (depth 2 — verbose detail) are excluded from the count.

**Rationale**: The original design counted only top-level results, which hid child failures (e.g., a variable `❌` not reflected in the summary). Counting children fixes this. However, grandchildren are verbose detail expansions — they break down an already-counted parent into per-key lines (like `FormatScanSummary` counts assessments, not individual steps). Counting them would inflate the summary and make the count `--verbose`-dependent: the same workspace in the same state would report different check counts depending on the flag. The summary should be stable regardless of display depth.

### D6: Section ordering follows scan prerequisite sequence

**Decision**: Sections are ordered: Providers, Policies, Cache, Complypacks, Workspace, Verification.

**Rationale**: Per @marcusburghardt: this follows the sequential requirements for performing a successful scan. Providers must exist before policies can call them, policies must be cached before complypacks are relevant, etc.

### D7: Skip variable results for unmapped providers

**Decision**: `CheckVariables` skips providers that have zero required variables (global and target) AND no policies mapping to them. No `CheckResult` is produced for these providers.

**Rationale**: Nesting variables under providers made the noise visible. A provider like `opa` installed but with no OPA policies configured would display `variables: 0/0 global vars, no target mapping for this evaluator` -- pure noise with no actionable information. The skip keeps the Providers section focused on providers that actually have something to validate.

**Guard**: The skip only triggers when all three conditions are true (no required globals, no required target vars, no policy mapping). Providers with any required variables or any policy reference still get a result, even if the result is passing.

### D8: Verbose detail produced as Children within Check* functions

**Decision**: `CheckVariables` and `CheckPolicyActivePeriod` produce verbose detail results as `Children` of their parent summary `CheckResult`, not as separate flat results. `attachByEvaluatorID` / `attachByPolicyID` only handle the summary-to-parent wiring.

**Rationale**: The initial implementation returned detail results as flat entries with names like `variables/<eid>/detail`. Since `extractID` returns the second segment, both `variables/ampel` and `variables/ampel/detail` extract to `ampel` — making them siblings as direct children of the provider instead of a three-level tree (provider → variables → detail). Producing detail as `Children` within the originating `Check*` function is correct because the relationship is inherent (detail belongs to its summary), not a cross-function relationship that `Run()` needs to wire.

**Result**: The renderer's existing three-level tree walk (top-level → children → grandchildren) displays detail at the deepest indent with message-only formatting, matching the spec intent.

### D9: Detail Status reflects actual check state

**Decision**: Verbose detail `CheckResult`s use `StatusPass` or `StatusFail` based on whether the individual variable is present, not a hardcoded `StatusPass`. The inline emoji (`✅`/`❌`) is removed from the message — the renderer provides the status emoji.

**Rationale**: The prior implementation hardcoded `Status: StatusPass` on all detail results and embedded pass/fail as emoji within the message string. This produced contradictory output: `✅ detail: target[k8s]: url ❌` — the leading emoji says pass, the trailing emoji says fail. Since the renderer prepends `statusEmoji(gc.Status)` at every level, the `Status` field must carry the truth.

**Guard**: Only affects verbose detail results (`--verbose`). Summary-level variable results already have correct `Status`.

### D10: Display-ready Name — eliminate renderer string parsing

**Problem**: D1 states the explicit `Children` field "eliminates all string parsing in the renderer." But the renderer still parses `Name` via `displayName` (strips prefix) and `childLabel` (extracts type label). This contradicts D1 and deviates from the codebase pattern where display structs (`summaryEntry`, `ExecutionPlanRow`) carry explicit display-ready fields.

`Name` serves dual duty:
- **Identity**: `extractID("provider/ampel")` → `"ampel"` for tree assembly
- **Display**: `displayName("provider/ampel")` → `"ampel"` for rendering

**Options considered**:

| Option | Change | Renderer | Identity | Invasiveness |
|--------|--------|----------|----------|-------------|
| A. Add `Label` field | New `Label string` on `CheckResult` | Uses `Label` (falls back to `Name`) | `extractID` on `Name` unchanged | Low — additive field, tests unchanged |
| B. Add `EvaluatorID` field | New `EvaluatorID string`; `Name` becomes display-ready | Uses `Name` directly | `attachBy*` matches on `EvaluatorID` | Medium — assembly helpers change, tests update Name assertions |
| C. Rewrite Name post-assembly | Mutate `Name` to short form after `Run()` assembly | Uses `Name` directly | `extractID` runs before mutation | Low code — but mutating data after assembly is fragile |
| D. Keep as-is | Document deviation from D1 | `displayName`/`childLabel` | `extractID` on `Name` | None — but D1 claim remains inaccurate |

**Decision**: **Option A — Add `Label string` field to `CheckResult`**.

**Rationale**:
- Follows the codebase pattern of explicit display fields (`summaryEntry.emoji`, `ExecutionPlanRow.TargetID`)
- `Name` retains its identity role — `extractID`, test assertions, and tree assembly are unaffected
- Renderer becomes a pure tree walk: display `Label` (or `Name` when `Label` is empty), no string parsing
- `displayName`, `childLabel` helpers and the `strings` import are removed from the renderer
- Zero-value (`""`) is backward-compatible — results without `Label` simply display `Name`

**Label values by level**:

| Level | Name (identity) | Label (display) |
|-------|-----------------|-----------------|
| Top-level provider | `provider/ampel` | `ampel` |
| Top-level policy | `policy/test-ampel-bp` | `test-ampel-bp` |
| Top-level complypack | `complypacks/ampel` | `ampel` |
| Top-level singleton | `cache`, `config`, `verification` | (empty — `Name` is already display-ready) |
| Child variable | `variables/ampel` | `variables` |
| Child active-period | `policy/cis/active-period` | `active-period` |
| Grandchild detail | (set by D8 — only `Message` rendered) | (empty — grandchildren render message-only) |

## Risks / Trade-offs

- **[Data model change]** Adding `Children` and `Group` to `CheckResult` is a structural change. All `Check*` functions and tests reference the struct. Mitigated by: the new fields are additive (zero-value is backward-compatible).
- **[Registry warning contract change]** `CheckPolicyVersions` no longer produces `registry/` named results. Tests that assert on `registry/` names need updating. Mitigated by: limited to `internal/doctor/` test file.
- **[Label field maintenance]** Every `CheckResult` with a qualified `Name` must set `Label` for correct display. Forgetting `Label` is safe (falls back to `Name`) but produces verbose output. Mitigated by: code review and the pattern being consistent across all `Check*` functions.
- **[Summary count change]** The summary counts top-level results and their direct children (depth 0–1). Grandchildren (verbose detail, depth 2) are excluded so the count is stable regardless of `--verbose`. The total may differ from the pre-redesign count because grouping and nesting change which results exist at counted depths.
