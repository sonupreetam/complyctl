// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// MockFetcher provides an in-memory mock registry for testing
type MockFetcher struct {
	mu       sync.RWMutex
	policies map[string]map[string]*MockPolicy
}

// MockPolicy represents a mock policy version
type MockPolicy struct {
	Digest   string
	Version  string
	Manifest []byte
}

// NewMockFetcher creates a mock fetcher with optional seed data
func NewMockFetcher() *MockFetcher {
	return &MockFetcher{
		policies: make(map[string]map[string]*MockPolicy),
	}
}

// AddPolicy registers a policy version for the mock
func (m *MockFetcher) AddPolicy(modulePath, version, digest string, manifest []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.policies[modulePath] == nil {
		m.policies[modulePath] = make(map[string]*MockPolicy)
	}
	m.policies[modulePath][version] = &MockPolicy{
		Digest:   digest,
		Version:  version,
		Manifest: manifest,
	}
	m.policies[modulePath]["latest"] = &MockPolicy{
		Digest:   digest,
		Version:  version,
		Manifest: manifest,
	}
}

// DefinitionVersion returns digest and version for modulePath.
// ref is empty or "latest" to use the "latest" entry; otherwise a specific version key.
func (m *MockFetcher) DefinitionVersion(_ context.Context, modulePath, ref string) (string, string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	versions, ok := m.policies[modulePath]
	if !ok {
		return "", "", fmt.Errorf("policy %s not found", modulePath)
	}
	key := "latest"
	if ref != "" && ref != "latest" {
		key = ref
	}
	p := versions[key]
	if p == nil && (ref == "" || ref == "latest") {
		p = versions["latest"]
	}
	if p == nil {
		return "", "", fmt.Errorf("policy %s not found for ref %q", modulePath, ref)
	}
	return p.Digest, p.Version, nil
}

// SeedTestPolicy adds default test policy for integration tests
func (m *MockFetcher) SeedTestPolicy(modulePath string) {
	manifest := map[string]interface{}{
		"policy_id": modulePath,
		"version":   "v1.0.0",
		"layers":    []string{"controls", "guidelines", "assessments"},
	}
	data, _ := json.Marshal(manifest)
	m.AddPolicy(modulePath, "v1.0.0", "sha256:abc123", data)
}
