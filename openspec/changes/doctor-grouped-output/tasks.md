<!-- Tasks use [P] to mark parallel-eligible items within a phase. -->
<!-- Sequential tasks (no marker) run first, then [P] tasks can run concurrently. -->

## 1. Data model: Group field and section ordering (done — uncommitted)

- [x] 1.1 Add `CheckGroup` type, group constants, and `GroupOrder` slice to `internal/doctor/doctor.go`
- [x] 1.2 Add `Group CheckGroup` field to `CheckResult` struct
- [x] 1.3 Assign `Group` to every `CheckResult` in all `Check*` functions
- [x] 1.4 Reorder `Run()` calls to match scan prerequisite sequence

## 2. Named complypacks (done — uncommitted)

- [x] 2.1 Change `CheckComplypacks` to list each complypack individually by evaluator-id
- [x] 2.2 Update `TestCheckComplypacks_AllPresent` to assert individual named results
- [x] 2.3 Update `TestCheckComplypacks_Missing` to assert evaluator-id in `Name`

## 3. Grouped output rendering with string-based nesting (done — uncommitted)

- [x] 3.1 Rewrite `printDiagnostics` with grouped section headers and indent
- [x] 3.2 Add `partitionVariables`, `buildChildMap`, `filterParents`, `extractEvaluatorID` helpers for string-based nesting
- [x] 3.3 Add `statusEmoji`, `countStatus`, `printResult` helpers

## 4. Data model: Add Children field to CheckResult

- [x] 4.1 Add `Children []CheckResult` field to `CheckResult` struct
- [x] 4.2 Remove `GroupVariables` from `GroupOrder` (variables are nested, not a section)
- [x] 4.3 Change `CheckVariables` resolve warnings and nil-config error from `Group: GroupVariables` to `Group: GroupProviders` so they render in the Providers section when `GroupVariables` is removed

## 5. Per-policy registry failure warnings

- [x] 5.1 Change `resolvePinnedFallback` network-failure fallback (line 462) to use `policy/<eid>` name instead of `registry/<reg>` — `eid` is already a parameter, just use it in the final return
- [x] 5.2 Remove the `unreachable` map and `continue` skip logic from `CheckPolicyVersions` — each policy from an unreachable registry gets its own warning via `resolvePinnedFallback`
- [x] 5.3 [P] Update 3 specific tests to expect per-policy results instead of `registry/` names:
  - `TestCheckPolicyVersions_RegistryUnreachable` (line 347: `registry/unreachable.io` → `policy/<eid>`)
  - `TestCheckPolicyVersions_MixedRegistries_Unreachable_And_404` (line 521: `registry/down.io` → `policy/<eid>`)
  - `TestCheckPolicyVersions_PinnedNetworkFailure_BothFail` (line 554-557: expect 2 results instead of 1, both `policy/<eid>` named)

## 6. Tree assembly in Run()

- [x] 6.1 Add `attachByEvaluatorID(providers, vars []CheckResult) []CheckResult` helper — matches variable results to provider results by evaluator-id, attaches as `Children`, returns unmatched results
- [x] 6.2 Add `attachByPolicyID(policies, active []CheckResult) []CheckResult` helper — matches active-period results to policy results by policy-id, attaches as `Children`, returns unmatched results
- [x] 6.3 Update `Run()` to call `attachByEvaluatorID` and `attachByPolicyID` after collecting flat results, assembling the tree before returning
- [x] 6.4 [P] Add unit tests for `attachByEvaluatorID`:
  - `TestAttachByEvaluatorID_MatchesCorrectly`
  - `TestAttachByEvaluatorID_UnmatchedPromoted`
  - `TestAttachByEvaluatorID_NoProviders`
- [x] 6.5 [P] Add unit tests for `attachByPolicyID`:
  - `TestAttachByPolicyID_MatchesCorrectly`
  - `TestAttachByPolicyID_UnmatchedPromoted`
  - `TestAttachByPolicyID_NoPolicies`

## 7. Simplify printDiagnostics to recursive tree walk

- [x] 7.1 Rewrite `printDiagnostics` as a tree walk: group by `Group`, iterate `GroupOrder`, render top-level `Name: Message`, recurse into `Children` with deeper indent displaying short `Name: Message`
- [x] 7.2 Count top-level and child results in the summary line (D5 — excludes grandchildren)
- [x] 7.3 Delete string-parsing helpers: `partitionVariables`, `buildChildMap`, `filterParents`, `isChildPolicyResult`, `extractEvaluatorID`
- [x] 7.4 Keep `statusEmoji`, `countStatus` helpers

## 8. Verification (phases 1-7)

- [x] 8.1 Run `make lint` and confirm 0 issues
- [x] 8.2 Run `go test -race -count=1 ./internal/doctor/ ./cmd/complyctl/cli/` and confirm all pass
- [x] 8.3 Run `make test-unit` and confirm no regressions beyond pre-existing failures
- [x] 8.4 Update issue #692 body with design decisions and Marcus feedback
- [x] 8.5 Post design rationale comment on issue #692

## 9. Verbose detail nesting (D8, D9)

- [x] 9.1 `CheckVariables`: collect verbose detail `CheckResult`s into a local slice and assign as `Children` of the summary `CheckResult` instead of appending to the flat `results` slice
- [x] 9.2 `CheckVariables`: set detail `Status` to `StatusFail` when variable is missing, `StatusPass` when present; remove inline emoji (`✅`/`❌`) from detail message
- [x] 9.3 [P] `CheckPolicyActivePeriod`: collect verbose detail `CheckResult`s into a local slice and assign as `Children` of the active-period `CheckResult` instead of appending to the flat `results` slice
- [x] 9.4 [P] Update or add unit tests for verbose detail nesting in `CheckVariables` and `CheckPolicyActivePeriod`

## 10. Display-ready Label field (D10)

- [x] 10.1 Add `Label string` field to `CheckResult` struct in `internal/doctor/doctor.go`
- [x] 10.2 Set `Label` on every `CheckResult` in `CheckProviders` (label: evaluator-id)
- [x] 10.3 Set `Label` on every `CheckResult` in `CheckPolicyVersions` (label: policy effective-id)
- [x] 10.4 [P] Set `Label` on every `CheckResult` in `CheckComplypacks` (label: evaluator-id)
- [x] 10.5 [P] Set `Label` on variable summary results in `CheckVariables` (label: `"variables"`)
- [x] 10.6 [P] Set `Label` on active-period results in `CheckPolicyActivePeriod` (label: `"active-period"`)
- [x] 10.7 [P] Set `Label` on resolve-warning results in `CheckVariables` (label: display-ready resolve name)
- [x] 10.8 Update `printDiagnostics`: use `Label` (with `Name` fallback) instead of `displayName`/`childLabel`; remove `displayName`, `childLabel` helpers and `strings` import; remove `countStatus` call for grandchildren (D5 — summary excludes depth 2)

## 11. Final verification

- [x] 11.1 Run `make lint` and confirm 0 issues
- [x] 11.2 Run `go test -race -count=1 ./internal/doctor/ ./cmd/complyctl/cli/` and confirm all pass
- [x] 11.3 Run `make test-unit` and confirm no regressions
- [x] 11.4 Build and run `complyctl doctor --verbose` against mock registry workspace to visually confirm three-level tree rendering and Label-based display
