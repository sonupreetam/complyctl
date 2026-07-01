## 1. Add cache integrity check to doctor

- [ ] 1.1 Add `CheckComplypackCacheIntegrity(cacheDir string, state *cache.State) []DiagnosticResult` to `internal/doctor/doctor.go`
- [ ] 1.2 For each complypack in `state.Complypacks`: resolve evaluator-id, check directory existence, check `content.tar.gz` presence
- [ ] 1.3 Return warning-level results for missing directories or content
- [ ] 1.4 Wire the new check into the doctor check list (non-blocking)

## 2. Unit tests

- [ ] 2.1 [P] `TestCheckComplypackCacheIntegrity_AllPresent` — pass case
- [ ] 2.2 [P] `TestCheckComplypackCacheIntegrity_MissingDir` — warning emitted
- [ ] 2.3 [P] `TestCheckComplypackCacheIntegrity_MissingContent` — warning emitted
- [ ] 2.4 [P] `TestCheckComplypackCacheIntegrity_EmptyState` — skipped cleanly

## 3. Verification

- [ ] 3.1 Run `go test -race ./internal/doctor/`
- [ ] 3.2 Run `go vet ./...`
- [ ] 3.3 E2E: manually delete a complypack dir, run `complyctl doctor`, verify warning appears
