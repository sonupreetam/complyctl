## Context

The cross-repo integration test (`ci_cross_repo_integration.yml`) runs the full
`complyctl get` / `generate` / `scan` pipeline with real provider binaries and produces
`evaluation-log-*.yaml` files as output. These files are the only runtime artifact that
can be validated against the Gemara CUE schema — they cannot be validated statically or
at compile time because they depend on live provider data flowing through gRPC.

The Gemara CUE schema (`gemaraproj/gemara`) defines constraints that Go's type system
does not enforce: minimum list lengths (`steps: [#AssessmentStep, ...#AssessmentStep]`),
enum values, required fields. The Go SDK (`go-gemara`) uses `[]*AssessmentStep` which
accepts `nil` — valid Go, invalid CUE. No validation API exists in `go-gemara` today.

## Goals / Non-Goals

**Goals:**

- Validate every EvaluationLog produced by the cross-repo test against the Gemara CUE
  schema before merge.
- Report the exact CUE constraint that failed, with the JSON path to the violating field.
- Keep CI overhead minimal (target: <10s added to the existing workflow).
- Ensure developers can reproduce the same validation locally in the devcontainer.

**Non-Goals:**

- Validating EvaluationLog inside the devcontainer in CI (unnecessary overhead).
- Adding schema validation to `complytime-providers` (separate concern, separate repo).
- Adding a `Validate()` API to `go-gemara` (upstream concern, tracked separately).
- Validating non-EvaluationLog Gemara artifacts (future extension).

## Decisions

### D1: Validate in the existing cross-repo workflow (not a separate workflow)

Schema validation is appended as steps after the integration test in
`ci_cross_repo_integration.yml`.

**Rationale**: The EvaluationLog files only exist after `complyctl scan` completes. The
cross-repo workflow already runs the full pipeline. A separate workflow would need to
duplicate the entire build+run pipeline just to produce the same files. Appending to the
existing workflow adds ~5s, not ~55s.

**Alternative rejected**: Standalone `ci_schema_validation.yml` workflow. Rejected because
it would duplicate the full pipeline (checkout both repos, build, install tools, run
get/generate/scan) adding ~55s of redundant CI time per PR.

### D2: Run cue vet bare-metal on the runner (not inside a devcontainer)

CUE is installed via `go install` on the Actions runner. The validation script runs
directly on the host.

**Rationale**: `cue vet` produces identical results regardless of execution environment.
The devcontainer adds ~90s of overhead (Docker image build without layer cache + CUE
compilation inside container) with no functional benefit to validation. The runner already
has Go installed (from the build steps) and the evaluation log files on disk.

**Alternative rejected**: Build devcontainer image and `docker run` with `cue vet` inside.
Rejected after investigation showed +90s CI overhead (47s Docker build + 30s CUE compile
+ 14s actual validation) vs ~5s bare-metal (CUE already compiled by `go install` in the
Go cache from the build steps).

### D3: Pin CUE binary and Gemara schema module versions

- CUE: `cuelang.org/go/cmd/cue@v0.17.1` (latest stable at time of implementation)
- Gemara schema: `github.com/gemaraproj/gemara@v0.23.0` (tracks `go-gemara v0.7.0`
  used by complyctl per `go.mod`)

**Rationale**: Pinning ensures reproducible CI results. The Gemara schema version must
track the `go-gemara` version in `go.mod` because complyctl produces output using those
Go types. A schema version newer than the Go types could introduce false failures for
fields not yet implemented.

**Update policy**: Bump the Gemara schema pin when `go-gemara` is updated in `go.mod`.
Bump CUE when a new stable release is available.

### D4: Install CUE in post-create.sh (not Containerfile)

CUE is added to `.devcontainer/scripts/post-create.sh` alongside existing `go install`
commands for snappy, ampel, and conftest.

**Rationale**: Follows the established pattern — `post-create.sh` installs project-specific
development tools via `go install`. The Containerfile installs system-level packages
(Go, git, make, curl, jq). CUE is a project-specific tool, not a system dependency.
Installing in `post-create.sh` ensures developers can reproduce CI schema validation
locally without additional manual setup.

**Alternative rejected**: Add to Containerfile via `RUN go install`. Rejected because it
breaks the existing separation of concerns (Containerfile = system packages,
post-create = project tools) and would require maintaining the version in an additional
location.

### D5: Validation script is a standalone reusable tool

`tests/schema-validation/validate.sh` accepts one or more directory paths, discovers
`evaluation-log-*.yaml` files, and validates each with `cue vet`. It exits 0 on all-pass,
1 on any failure, 2 on configuration error (no files found).

**Rationale**: A standalone script is invocable from CI, from `make`, and from a developer
shell. It decouples the validation logic from the CI workflow YAML, making it testable
and reusable. The same script works whether invoked by the cross-repo workflow or by a
developer running it manually in the devcontainer.

## Risks / Trade-offs

**[CUE Central Registry dependency]** → `cue vet` fetches the Gemara schema module from
the CUE Central Registry at runtime. If the registry is unavailable, validation fails.
Mitigation: the registry is highly available; pinned version ensures no surprise breakage
from upstream schema changes.

**[Version drift]** → If `go-gemara` is bumped in `go.mod` but the Gemara schema pin in
`validate.sh` is not updated, false failures may occur. Mitigation: documented update
policy; version comment in `validate.sh` references the `go.mod` dependency.

**[Provider bugs surface in complyctl CI]** → A schema violation caused by a provider bug
(on `providers@main`) will fail complyctl PRs. Mitigation: this is the desired behavior —
it signals that the integrated output is non-conformant. The fix belongs in providers,
but complyctl CI is the integration boundary where both sides meet.

## Open Questions

- Should `complytime-providers` add its own `cue vet` validation to catch bugs at source
  before they reach `providers@main`? (Recommended: yes, as a follow-up issue.)
- Should `go-gemara` provide a programmatic `Validate()` API? (Recommended: yes, as an
  upstream feature request, to enable unit-test-level validation without shelling out.)
