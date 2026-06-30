## 1. Add evaluator-id uniqueness validation

- [ ] 1.1 Add `validateUniqueEvaluatorIDs(state *cache.State, complypacks []complytime.PolicyEntry) error` to `cmd/complyctl/cli/get.go`
- [ ] 1.2 Build a map of evaluator-id → []repository from `state.Complypacks`
- [ ] 1.3 For entries with len > 1, format error listing all conflicting repositories
- [ ] 1.4 Call `validateUniqueEvaluatorIDs` after `syncAllComplypacks` returns in `syncComplypacks`

## 2. Unit tests

- [ ] 2.1 [P] `TestValidateUniqueEvaluatorIDs_NoDuplicates` — passes cleanly
- [ ] 2.2 [P] `TestValidateUniqueEvaluatorIDs_DuplicateDetected` — returns error with both repos
- [ ] 2.3 [P] `TestValidateUniqueEvaluatorIDs_SingleEntry` — passes
- [ ] 2.4 [P] `TestValidateUniqueEvaluatorIDs_EmptyState` — passes

## 3. Verification

- [ ] 3.1 Run `go test -race ./cmd/complyctl/cli/`
- [ ] 3.2 Run `go vet ./...`
- [ ] 3.3 E2E: configure two complypacks with same evaluator-id, run `complyctl get`, verify error
