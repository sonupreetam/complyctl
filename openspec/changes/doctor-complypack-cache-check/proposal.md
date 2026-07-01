## Why

`ComplypackSync.SyncComplypack()` checks freshness by comparing `state.json` digest against remote. If the local cache directory has been deleted but `state.json` still records a matching digest, sync is permanently skipped. `LookupByEvaluatorID` then returns an empty path, the provider gets no complypack content, and scan may produce incomplete results without error.

There is no user-facing detection or self-healing for this state. `complyctl doctor` currently checks whether evaluator-ids are present in complypack config but does not verify file integrity.

This is tracked in [#649](https://github.com/complytime/complyctl/issues/649).

## What Changes

- `complyctl doctor` gains a new check that verifies each complypack entry in `state.json` has a corresponding cache directory with `content.tar.gz` present.
- When the directory is missing, doctor emits a warning suggesting `complyctl get --force` or clearing state.

## Capabilities

### New Capabilities
- `doctor-complypack-integrity`: Verify that complypack cache directories referenced by `state.json` actually exist on disk with expected content files.

### Modified Capabilities
- `doctor`: Gains an additional diagnostic check in the complypack section.

### Removed Capabilities
- None.

## Impact

- `internal/doctor/doctor.go` — new `CheckComplypackCacheIntegrity` function
- `internal/doctor/doctor_test.go` — tests for the new check
- No CLI flag changes (check is always run as part of `complyctl doctor`)

## Constitution Alignment

### I. Autonomous Collaboration

**Assessment**: PASS

Doctor proactively detects and reports the inconsistency with actionable remediation advice.

### II. Composability First

**Assessment**: PASS

The check is additive — existing checks are unaffected. The new check is independent and scoped to complypack cache integrity only.

### III. Observable Quality

**Assessment**: PASS

Clear diagnostic output identifying which complypack entries have missing cache directories, with suggested fix commands.

### IV. Testability

**Assessment**: PASS

Testable by pre-seeding `state.json` with complypack entries and verifying the check detects missing/present directories.
