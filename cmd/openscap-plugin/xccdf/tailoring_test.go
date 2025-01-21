// SPDX-License-Identifier: Apache-2.0

package xccdf

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/ComplianceAsCode/compliance-operator/pkg/xccdf"
	"github.com/complytime/complytime/cmd/openscap-plugin/config"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
)

// TestGetTailoringID tests the getTailoringID function.
func TestGetTailoringID(t *testing.T) {
	expected := "xccdf_complytime.openscapplugin_tailoring_complytime-tailoring-profile"
	result := getTailoringID()

	if result != expected {
		t.Errorf("getTailoringID() = %v; want %v", result, expected)
	}
}

// TestGetTailoringProfileID tests the getTailoringProfileID function.
func TestGetTailoringProfileID(t *testing.T) {
	profileId := "test-profile"
	expected := "xccdf_complytime.openscapplugin_profile_test-profile_complytime-tailoring-profile"
	result := getTailoringProfileID(profileId)

	if result != expected {
		t.Errorf("getTailoringProfileID(%v) = %v; want %v", profileId, result, expected)
	}
}

// TestGetTailoringProfileTitle tests the getTailoringProfileTitle function.
func TestGetTailoringProfileTitle(t *testing.T) {
	dsPath := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap", "ssg-rhel-ds.xml")

	tests := []struct {
		profileId string
		dsPath    string
		expected  string
	}{
		{"test_profile", dsPath, "ComplyTime Tailoring Profile - Test Profile"},
		{"no-ds-profile", dsPath, "ComplyTime Tailoring Profile - no-ds-profile"},
		{"test_profile_no_title", dsPath, "ComplyTime Tailoring Profile - test_profile_no_title"},
	}

	for _, tt := range tests {
		result := getTailoringProfileTitle(tt.profileId, tt.dsPath)
		if result != tt.expected {
			t.Errorf("getTailoringProfileTitle(%v, %v) = %v; want %v", tt.profileId, tt.dsPath, result, tt.expected)
		}
	}
}

// TestGetTailoringVersion tests the getTailoringVersion function.
func TestGetTailoringVersion(t *testing.T) {
	result := getTailoringVersion()

	if result.Value != "1" {
		t.Errorf("getTailoringVersion().Value = %v; want %v", result.Value, "1")
	}

	_, err := time.Parse(time.RFC3339, result.Time)
	if err != nil {
		t.Errorf("getTailoringVersion().Time = %v; not in RFC3339 format", result.Time)
	}
}

// TestGetTailoringBenchmarkHref tests the getTailoringBenchmarkHref function.
func TestGetTailoringBenchmarkHref(t *testing.T) {
	dsPath := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap", "ssg-rhel-ds.xml")
	expected := xccdf.BenchmarkElement{
		Href: dsPath,
	}
	result := getTailoringBenchmarkHref(dsPath)

	if result != expected {
		t.Errorf("getTailoringBenchmarkHref(%v) = %v; want %v", dsPath, result, expected)
	}
}

// TestGetTailoringProfile tests the getTailoringProfile function.
func TestGetTailoringProfile(t *testing.T) {
	profileId := "test-profile"
	dsPath := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap", "ssg-rhel-ds.xml")
	tailoringPolicy := policy.Policy{
		{
			Rule: extensions.Rule{
				ID:          "rule1",
				Description: "rule1_description",
				Parameter: &extensions.Parameter{
					ID:          "param1",
					Description: "param1_description",
					Value:       "value1",
				},
			},
		},
	}

	expected := xccdf.ProfileElement{
		ID: getTailoringProfileID(profileId),
		Title: &xccdf.TitleOrDescriptionElement{
			Value: getTailoringProfileTitle(profileId, dsPath),
		},
		Selections: []xccdf.SelectElement{
			{IDRef: "rule1", Selected: true},
		},
		Values: []xccdf.SetValueElement{
			{IDRef: "param1", Value: "value1"},
		},
	}

	result := getTailoringProfile(profileId, tailoringPolicy, dsPath)

	if result.ID != expected.ID {
		t.Errorf("getTailoringProfile().ID = %v; want %v", result.ID, expected.ID)
	}
	if result.Title.Value != expected.Title.Value {
		t.Errorf("getTailoringProfile().Title.Value = %v; want %v", result.Title.Value, expected.Title.Value)
	}
	if len(result.Selections) != len(expected.Selections) {
		t.Errorf("getTailoringProfile().Selections = %v; want %v", result.Selections, expected.Selections)
	}
	if len(result.Values) != len(expected.Values) {
		t.Errorf("getTailoringProfile().Values = %v; want %v", result.Values, expected.Values)
	}
}

// TestGetPolicySelections tests the getPolicySelections function.
func TestGetPolicySelections(t *testing.T) {
	tailoringPolicy := policy.Policy{
		{
			Rule: extensions.Rule{
				ID:          "rule1",
				Description: "rule1_description",
			},
		},
		{
			Rule: extensions.Rule{
				ID:          "rule2",
				Description: "rule2_description",
			},
		},
	}

	expected := []xccdf.SelectElement{
		{IDRef: "rule1", Selected: true},
		{IDRef: "rule2", Selected: true},
	}

	result := getPolicySelections(tailoringPolicy)

	if len(result) != len(expected) {
		t.Errorf("getPolicySelections() length = %v; want %v", len(result), len(expected))
	}

	for i, selection := range result {
		if selection.IDRef != expected[i].IDRef || selection.Selected != expected[i].Selected {
			t.Errorf("getPolicySelections()[%d] = %v; want %v", i, selection, expected[i])
		}
	}
}

// TestGetValuesFromPolicyVariables tests the getValuesFromPolicyVariables function.
func TestGetValuesFromPolicyVariables(t *testing.T) {
	tests := []struct {
		name            string
		tailoringPolicy policy.Policy
		expected        []xccdf.SetValueElement
	}{
		{
			name: "Single rule with parameter",
			tailoringPolicy: policy.Policy{
				{
					Rule: extensions.Rule{
						ID:          "rule1",
						Description: "rule1_description",
						Parameter: &extensions.Parameter{
							ID:          "param1",
							Description: "param1_description",
							Value:       "value1",
						},
					},
				},
			},
			expected: []xccdf.SetValueElement{
				{IDRef: "param1", Value: "value1"},
			},
		},
		{
			name: "Multiple rules with parameters",
			tailoringPolicy: policy.Policy{
				{
					Rule: extensions.Rule{
						ID:          "rule1",
						Description: "rule1_description",
						Parameter: &extensions.Parameter{
							ID:          "param1",
							Description: "param1_description",
							Value:       "value1",
						},
					},
				},
				{
					Rule: extensions.Rule{
						ID:          "rule2",
						Description: "rule2_description",
						Parameter: &extensions.Parameter{
							ID:          "param2",
							Description: "param2_description",
							Value:       "value2",
						},
					},
				},
			},
			expected: []xccdf.SetValueElement{
				{IDRef: "param1", Value: "value1"},
				{IDRef: "param2", Value: "value2"},
			},
		},
		{
			name: "Rule without parameter",
			tailoringPolicy: policy.Policy{
				{
					Rule: extensions.Rule{
						ID:          "rule1",
						Description: "rule1_description",
					},
				},
			},
			expected: nil,
		},
		{
			name:            "Empty policy",
			tailoringPolicy: policy.Policy{},
			expected:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getValuesFromPolicyVariables(tt.tailoringPolicy)
			if len(result) != len(tt.expected) {
				t.Errorf("getValuesFromPolicyVariables() length = %v; want %v", len(result), len(tt.expected))
			}
			for i, value := range result {
				if value.IDRef != tt.expected[i].IDRef || value.Value != tt.expected[i].Value {
					t.Errorf("getValuesFromPolicyVariables()[%d] = %v; want %v", i, value, tt.expected[i])
				}
			}
		})
	}
}

// TestPolicyToXML tests the PolicyToXML function.
func TestPolicyToXML(t *testing.T) {
	tempDir := t.TempDir()
	dsPath := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap", "ssg-rhel-ds.xml")
	profileId := "test_profile"

	tailoringPolicy := policy.Policy{
		{
			Rule: extensions.Rule{
				ID:          "rule1",
				Description: "rule1_description",
				Parameter: &extensions.Parameter{
					ID:          "param1",
					Description: "param1_description",
					Value:       "value1",
				},
			},
		},
	}

	cfg := &config.Config{
		Files: struct {
			PluginDir  string "yaml:\"plugindir\""
			Workspace  string "yaml:\"workspace\""
			Datastream string "yaml:\"datastream\""
			Results    string "yaml:\"results\""
			ARF        string "yaml:\"arf\""
			Policy     string "yaml:\"policy\""
		}{
			PluginDir:  "plugins",
			Workspace:  filepath.Join(tempDir, "workspace"),
			Datastream: dsPath,
			Results:    "results.xml",
			ARF:        "arf.xml",
			Policy:     "absent_policy.yaml",
		},
		Parameters: struct {
			Profile string `yaml:"profile"`
		}{
			Profile: profileId,
		},
	}

	expectedXML := `<?xml version="1.0" encoding="UTF-8"?>
<xccdf-1.2:Tailoring xmlns:xccdf-1.2="http://checklists.nist.gov/xccdf/1.2" id="xccdf_complytime.openscapplugin_tailoring_complytime-tailoring-profile">
  <xccdf-1.2:benchmark href="` + dsPath + `"></xccdf-1.2:benchmark>
  <xccdf-1.2:version time="` + getTailoringVersion().Time + `">1</xccdf-1.2:version>
  <xccdf-1.2:Profile id="xccdf_complytime.openscapplugin_profile_test_profile_complytime-tailoring-profile">
    <xccdf-1.2:title override="false">ComplyTime Tailoring Profile - Test Profile</xccdf-1.2:title>
    <xccdf-1.2:select idref="rule1" selected="true"></xccdf-1.2:select>
    <xccdf-1.2:set-value idref="param1">value1</xccdf-1.2:set-value>
  </xccdf-1.2:Profile>
</xccdf-1.2:Tailoring>`

	result, err := PolicyToXML(tailoringPolicy, cfg)
	if err != nil {
		t.Fatalf("PolicyToXML() error = %v", err)
	}

	if result != expectedXML {
		t.Errorf("PolicyToXML() = %v; want %v", result, expectedXML)
	}
}
