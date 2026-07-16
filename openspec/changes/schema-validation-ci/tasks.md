## 1. Validation Script

- [ ] 1.1 Create `tests/schema-validation/validate.sh` that accepts directory arguments,
  finds `evaluation-log-*.yaml` files, and validates each with
  `cue vet -c -d '#EvaluationLog' <module> <file>`
- [ ] 1.2 Pin Gemara CUE module to `github.com/gemaraproj/gemara@v0.23.0` (tracking
  `go-gemara v0.7.0` in `go.mod`) with override via `GEMARA_CUE_MODULE` env var
- [ ] 1.3 Report per-file PASS/FAIL with native `cue vet` error output on failure
- [ ] 1.4 Exit 0 (all pass), 1 (any fail), or 2 (no files found)

## 2. Cross-Repo Test Script Modification

- [ ] 2.1 Add conditional `SCAN_OUTPUT_DIR` export block at end of
  `tests/cross-repo/cross_repo_integration_test.sh` that copies
  `evaluation-log-*.yaml` files when the env var is set

## 3. CI Workflow Modification

- [ ] 3.1 Add `workflow_dispatch` trigger to `ci_cross_repo_integration.yml`
- [ ] 3.2 Add `SCAN_OUTPUT_DIR` env var to the integration test step
- [ ] 3.3 Add step to install CUE: `go install cuelang.org/go/cmd/cue@v0.17.1`
- [ ] 3.4 Add step to run `./tests/schema-validation/validate.sh` on the exported
  scan output directory

## 4. Devcontainer Tooling

- [ ] 4.1 Add `go install cuelang.org/go/cmd/cue@v0.17.1` to
  `.devcontainer/scripts/post-create.sh` alongside existing snappy/ampel/conftest
  installs

## 5. Verification

- [ ] 5.1 Run the cross-repo integration test locally with `SCAN_OUTPUT_DIR` set and
  confirm evaluation log files are exported
- [ ] 5.2 Run `validate.sh` on the exported files and confirm PASS for all
- [ ] 5.3 Confirm CI workflow passes on the PR branch
- [ ] 5.4 Confirm `make test-unit` and `make test-integration` still pass (no regressions)
