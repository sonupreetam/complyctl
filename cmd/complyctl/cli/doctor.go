// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/complytime/complyctl/internal/cache"
	"github.com/complytime/complyctl/internal/complytime"
	"github.com/complytime/complyctl/internal/doctor"
	"github.com/complytime/complyctl/internal/policy"
	"github.com/complytime/complyctl/internal/registry"
)

func doctorCmd(common *Common) *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run pre-flight diagnostics on the workspace",
		Long: `Run pre-flight diagnostics on the workspace.

Checks include provider discovery, policy cache integrity, workspace
configuration validation, and complypack availability. When complypacks are
configured, the doctor verifies that each referenced complypack is cached
and reports missing entries.`,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			baseDir, err := common.ResolveWorkspace()
			if err != nil {
				return err
			}
			return runDoctor(baseDir, verbose)
		},
	}
	cmd.Flags().BoolVar(&verbose, "verbose", false, "expand per-provider variable detail")
	return cmd
}

// registryVersionResolver adapts registry.Client to doctor.VersionResolver.
// See R55: specs/001-gemara-native-workflow/spec.md
type registryVersionResolver struct {
	timeout time.Duration
}

func (r *registryVersionResolver) ResolveLatestVersion(registryURL, repository string) (string, error) {
	return r.resolve(registryURL, repository, "")
}

func (r *registryVersionResolver) ResolveVersion(registryURL, repository, version string) (string, error) {
	return r.resolve(registryURL, repository, version)
}

func (r *registryVersionResolver) resolve(registryURL, repository, version string) (string, error) {
	credFunc, err := registry.NewCredentialFunc()
	if err != nil {
		credFunc = nil
	}
	client := registry.NewClient(registryURL, credFunc)
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	lookup := repository
	if version != "" {
		lookup = repository + ":" + version
	}
	_, resolved, err := client.DefinitionVersion(ctx, lookup)
	if err != nil {
		return "", err
	}
	return resolved, nil
}

// See FR-039, R44, R51, R52, R55: specs/001-gemara-native-workflow/spec.md
func runDoctor(baseDir string, verbose bool) error {
	providerDir, err := complytime.ResolveProviderDir()
	if err != nil {
		return fmt.Errorf("failed to resolve provider directory: %w", err)
	}

	cacheDir, err := complytime.ResolveCacheDir()
	if err != nil {
		return fmt.Errorf("failed to resolve cache directory: %w", err)
	}

	dataDir, err := complytime.ResolveDataDir()
	if err != nil {
		return fmt.Errorf("failed to resolve data directory: %w", err)
	}

	ws := complytime.NewWorkspace(baseDir)
	configPath := ws.Path()
	var cfg *complytime.WorkspaceConfig

	if loadErr := ws.Load(); loadErr == nil {
		cfg = ws.Config()
	}

	var resolver doctor.PolicyGraphResolver
	cacheMgr := cache.NewCache(cacheDir)
	loader := policy.NewLoader(cacheMgr)
	resolver = policy.NewResolver(loader)

	versionResolver := &registryVersionResolver{timeout: 5 * time.Second}

	results := doctor.Run(cfg, configPath, providerDir, cacheDir, dataDir, resolver, versionResolver, verbose, logger)
	return printDiagnostics(results)
}

func printDiagnostics(results []doctor.CheckResult) error {
	fmt.Println("Running workspace diagnostics...")
	fmt.Println()

	grouped := make(map[doctor.CheckGroup][]doctor.CheckResult)
	for _, r := range results {
		grouped[r.Group] = append(grouped[r.Group], r)
	}

	var passCount, failCount, warnCount int
	hasBlockingFailure := false
	firstSection := true

	for _, group := range doctor.GroupOrder() {
		checks, ok := grouped[group]
		if !ok {
			continue
		}

		if !firstSection {
			fmt.Println()
		}
		firstSection = false
		fmt.Println(string(group))

		for _, r := range checks {
			emoji := statusEmoji(r.Status)
			countStatus(r.Status, &passCount, &failCount, &warnCount)
			fmt.Printf("  %s %s: %s\n", emoji, resultLabel(r), r.Message)
			if r.Blocking && r.Status == doctor.StatusFail {
				hasBlockingFailure = true
			}

			for _, child := range r.Children {
				childEmoji := statusEmoji(child.Status)
				countStatus(child.Status, &passCount, &failCount, &warnCount)
				fmt.Printf("      %s %s: %s\n", childEmoji, resultLabel(child), child.Message)
				if child.Blocking && child.Status == doctor.StatusFail {
					hasBlockingFailure = true
				}

				// Grandchildren (verbose detail) are rendered but not
				// counted in the summary (D5 — stable count regardless
				// of --verbose).
				for _, gc := range child.Children {
					gcEmoji := statusEmoji(gc.Status)
					fmt.Printf("          %s %s\n", gcEmoji, gc.Message)
				}
			}
		}
	}

	total := passCount + failCount + warnCount
	fmt.Printf("\n%d checks: %d passed, %d failed, %d warnings\n", total, passCount, failCount, warnCount)

	if hasBlockingFailure {
		return fmt.Errorf("one or more blocking checks failed")
	}
	return nil
}

// resultLabel returns the display label for a CheckResult. Uses Label
// when set, falls back to Name when Label is empty. No string parsing
// of Name occurs — display text is set at the source (D10).
func resultLabel(r doctor.CheckResult) string {
	if r.Label != "" {
		return r.Label
	}
	return r.Name
}

func statusEmoji(s doctor.CheckStatus) string {
	switch s {
	case doctor.StatusPass:
		return complytime.StatusPassed
	case doctor.StatusFail:
		return complytime.StatusFailed
	case doctor.StatusWarn:
		return complytime.StatusSkipped
	default:
		return "?"
	}
}

func countStatus(s doctor.CheckStatus, pass, fail, warn *int) {
	switch s {
	case doctor.StatusPass:
		*pass++
	case doctor.StatusFail:
		*fail++
	case doctor.StatusWarn:
		*warn++
	}
}
