## Context

`complyctl doctor` runs pre-flight diagnostic checks. The existing `CheckComplypacks` verifies evaluator-id presence in config but not that cache files actually exist on disk.

## Goals

- Detect when `state.json` claims a complypack is cached but the files are missing
- Provide actionable remediation advice
- Non-blocking check (warning, not failure)

## Non-Goals

- Verifying content integrity (checksums/signatures)
- Adding a `--force` flag to `complyctl get` (separate issue)
- Auto-healing (re-fetching missing complypacks)

## Decisions

### D1: Check placement

Add `CheckComplypackCacheIntegrity` as a new non-blocking check in the doctor check list. It runs after `CheckComplypacks` (which validates config).

### D2: What to verify

For each complypack entry in `state.json`:
1. Resolve evaluator-id from the entry
2. Check that `{cacheDir}/complypacks/{evaluator-id}/` exists
3. Check that at least one version subdirectory contains `content.tar.gz`

### D3: Non-blocking

This is a WARNING-level check, not a blocking check. Missing cache doesn't prevent all operations — only scan/generate for that specific evaluator would fail.

### D4: Remediation message

Format: `Complypack cache missing for evaluator "{id}" — run "complyctl get" to restore`

## Risks / Trade-offs

- The check requires loading `state.json` which doctor already does for other checks — no additional I/O cost.
- If #646 (version eviction) lands first, there's guaranteed to be at most one version dir, simplifying the "find content.tar.gz" logic.
