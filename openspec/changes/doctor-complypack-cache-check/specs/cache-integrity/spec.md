## ADDED Requirements

### FR-001: Detect missing complypack cache directory

**Given** `state.json` records a complypack entry with evaluator-id `E` and digest `D`
**And** the directory `~/.complytime/complypacks/E/` does not exist
**When** `complyctl doctor` runs
**Then** a warning MUST be emitted indicating the cache is missing for evaluator `E`
**And** the warning MUST suggest running `complyctl get` to restore

### FR-002: Detect missing content.tar.gz

**Given** `state.json` records a complypack entry with evaluator-id `E`
**And** the evaluator directory exists but contains no `content.tar.gz`
**When** `complyctl doctor` runs
**Then** a warning MUST be emitted indicating the complypack content is incomplete

### FR-003: Pass when cache is consistent

**Given** `state.json` records a complypack entry with evaluator-id `E`
**And** `~/.complytime/complypacks/E/{version}/content.tar.gz` exists
**When** `complyctl doctor` runs
**Then** the complypack cache integrity check MUST pass

### FR-004: Handle empty state (no complypacks configured)

**Given** `state.json` has no complypack entries (or does not exist)
**When** `complyctl doctor` runs
**Then** the complypack cache integrity check MUST be skipped without error
