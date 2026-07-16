## ADDED Requirements

### Requirement: CUE schema validation in cross-repo CI workflow

The cross-repo integration workflow SHALL validate each `evaluation-log-*.yaml` file
produced by `complyctl scan` against the Gemara CUE schema using `cue vet`. Validation
SHALL run on the GitHub Actions runner (bare-metal), not inside a container. The workflow
SHALL fail if any file violates a schema constraint.

#### Scenario: All evaluation logs pass schema validation

- **WHEN** the cross-repo integration test completes successfully and produces
  evaluation log files
- **THEN** `cue vet -c -d '#EvaluationLog'` passes for each file and the workflow
  succeeds

#### Scenario: A schema violation is detected

- **WHEN** an evaluation log file violates a CUE schema constraint (e.g., empty `steps`
  array)
- **THEN** the workflow fails and the CI log reports the exact constraint path and
  violation message from `cue vet`

#### Scenario: No evaluation log files found

- **WHEN** the integration test fails before producing scan output
- **THEN** the validation step reports a configuration error (exit code 2) and the
  workflow fails

### Requirement: Validation script

The complyctl repository SHALL provide a reusable validation script at
`tests/schema-validation/validate.sh` that accepts one or more directory paths as
arguments, discovers `evaluation-log-*.yaml` files in those directories, and validates
each against the Gemara CUE schema module.

#### Scenario: Script validates files in given directories

- **WHEN** `./tests/schema-validation/validate.sh <dir>` is invoked with a directory
  containing evaluation log files
- **THEN** the script validates each file and prints per-file PASS/FAIL status

#### Scenario: Script reports constraint errors

- **WHEN** a file fails validation
- **THEN** the script prints the `cue vet` error output showing the exact constraint
  and JSON path that was violated

#### Scenario: Script exits with distinct codes

- **WHEN** all files pass → exit 0
- **WHEN** one or more files fail → exit 1
- **WHEN** no evaluation log files are found → exit 2

### Requirement: Pinned CUE and Gemara schema versions

The CI workflow SHALL install a pinned version of CUE (`cuelang.org/go/cmd/cue@v0.17.1`).
The validation script SHALL use a pinned Gemara schema module version
(`github.com/gemaraproj/gemara@v0.23.0`) that tracks the `go-gemara` version in
complyctl's `go.mod`.

#### Scenario: Version pinning ensures reproducibility

- **WHEN** the workflow runs on two different dates without code changes
- **THEN** the validation uses the same CUE binary version and the same Gemara schema
  version, producing identical results

### Requirement: Scan output export for validation

The cross-repo integration test script SHALL export `evaluation-log-*.yaml` files to
a directory specified by the `SCAN_OUTPUT_DIR` environment variable when that variable
is set. When `SCAN_OUTPUT_DIR` is not set, no export occurs (local development behavior
is unchanged).

#### Scenario: CI sets SCAN_OUTPUT_DIR

- **WHEN** `SCAN_OUTPUT_DIR` is set and the integration test completes
- **THEN** all `evaluation-log-*.yaml` files from `.complytime/scan/` are copied to
  the specified directory

#### Scenario: Local development without SCAN_OUTPUT_DIR

- **WHEN** `SCAN_OUTPUT_DIR` is not set
- **THEN** the test script runs normally with no export side-effect

### Requirement: CUE available in devcontainer

The devcontainer `post-create.sh` script SHALL install CUE
(`cuelang.org/go/cmd/cue@v0.17.1`) so developers can run `validate.sh` locally without
manual tool installation.

#### Scenario: Developer validates schema locally

- **WHEN** a developer opens the devcontainer and runs
  `./tests/schema-validation/validate.sh .complytime/scan/`
- **THEN** CUE is available and validation executes without additional setup

### Requirement: Manual workflow dispatch

The cross-repo integration workflow SHALL support `workflow_dispatch` for manual
triggering from the GitHub Actions UI.

#### Scenario: Manual trigger

- **WHEN** a maintainer triggers the workflow manually via the Actions UI
- **THEN** the full pipeline (integration test + schema validation) runs
