## Why

The cross-repo integration test validates structural correctness of complyctl's scan
output: file existence, correct requirement IDs, pipeline exit codes. It does not
validate that the EvaluationLog YAML conforms to the Gemara CUE schema. This gap
allowed `complytime/complytime-providers#63` to reach main — the OPA provider returned
`steps: []`, violating the schema's minimum-one-step constraint
(`steps: [#AssessmentStep, ...#AssessmentStep]`). Go's type system permits empty slices
where CUE requires at least one element, so the bug was invisible to existing tests.

Adding `cue vet` on the real EvaluationLog files produced by the cross-repo integration
test closes this gap at the point where the files are already generated — no duplicate
pipeline, no additional infrastructure.

## What Changes

- **CI workflow modification**: `.github/workflows/ci_cross_repo_integration.yml` gains
  a CUE install step and a schema validation step that runs `cue vet` on each
  `evaluation-log-*.yaml` file produced by the integration test.
- **New validation script**: `tests/schema-validation/validate.sh` accepts directories
  as arguments, finds EvaluationLog files, and validates each against the Gemara CUE
  schema module. Reports per-file pass/fail with native `cue vet` error output.
- **Test script modification**: `tests/cross-repo/cross_repo_integration_test.sh` exports
  evaluation log files to `SCAN_OUTPUT_DIR` when the env var is set (CI only).
- **Devcontainer tooling**: `.devcontainer/scripts/post-create.sh` installs CUE so
  developers can reproduce schema validation locally.
- **Workflow trigger**: `workflow_dispatch` added for manual triggering.

## Capabilities

### New Capabilities

- `schema-validation`: Automated CUE schema validation of EvaluationLog output in CI,
  catching constraint violations (cardinality, required fields, enum values) that Go's
  type system does not enforce.

### Modified Capabilities

- `cross-repo-integration-test`: Extended with schema validation step and scan output
  export mechanism.

## Impact

- **ci_cross_repo_integration.yml**: Two new steps (~5s total overhead on bare-metal).
- **cross_repo_integration_test.sh**: Conditional export block at end of script.
- **post-create.sh**: One additional `go install` line (consistent with existing pattern).
- **No new workflows**: Validation is appended to the existing cross-repo workflow.
- **No devcontainer usage in CI**: CUE runs directly on the runner. Devcontainer installs
  CUE for local developer use only.

## Constitution Alignment

| Principle | Assessment |
|:----------|:-----------|
| I. Single Source of Truth | Validation script is the single implementation; CI and local dev both invoke it. |
| II. Simplicity & Isolation | No Docker, no separate workflow. Two new steps in an existing job. |
| III. Incremental Improvement | Adds schema validation without restructuring the existing test infrastructure. |
| IV. Readability First | `validate.sh` is self-documenting with clear pass/fail reporting. |
| V. Do Not Reinvent the Wheel | Uses `cue vet` (the CUE project's official validation tool) directly. |
| VI. Composability | `validate.sh` is reusable: CI passes a directory, developers pass their own paths. |
| VII. Convention Over Configuration | Follows the existing `post-create.sh` pattern for tool installation. |
