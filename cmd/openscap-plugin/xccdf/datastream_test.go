// SPDX-License-Identifier: Apache-2.0

package xccdf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/antchfx/xmlquery"
)

// TestGetDsProfileID tests the getDsProfileID function.
func TestGetDsProfileID(t *testing.T) {
	tests := []struct {
		profileId string
		expected  string
	}{
		{"test", "xccdf_org.ssgproject.content_profile_test"},
		{"profile1", "xccdf_org.ssgproject.content_profile_profile1"},
		{"", "xccdf_org.ssgproject.content_profile_"},
	}

	for _, tt := range tests {
		t.Run(tt.profileId, func(t *testing.T) {
			result := getDsProfileID(tt.profileId)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestGetDsElement tests the getDsElement function with the Profile "description" element.
func TestGetDsElement(t *testing.T) {
	dsFile := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap", "ssg-rhel-ds.xml")
	dsData, err := os.ReadFile(dsFile)
	if err != nil {
		t.Fatalf("failed to read Datastream file: %v", err)
	}

	doc, err := xmlquery.Parse(strings.NewReader(string(dsData)))
	if err != nil {
		t.Fatalf("failed to parse Datastream: %v", err)
	}

	tests := []struct {
		dsElement string
		expected  string
		wantErr   bool
	}{
		{"//xccdf-1.2:Profile[@id='xccdf_org.ssgproject.content_profile_test_profile']", "This profile is only used for Unit Tests", false},
		{"//xccdf-1.2:Profile[@id='xccdf_org.ssgproject.content_profile_absent']", "", false},
		{"//invalid:Profile[@id='xccdf_org.ssgproject.content_profile_test_profile']", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.dsElement, func(t *testing.T) {
			result, err := getDsElement(doc, tt.dsElement)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDsElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != nil && err == nil {
				profileDescription := result.SelectElement("xccdf-1.2:description")
				if profileDescription.InnerText() != tt.expected {
					t.Errorf("got %s, want %s", profileDescription.InnerText(), tt.expected)
				}
			}
		})
	}
}

// TestLoadDataStream tests the loadDataStream function.
// Errors are expected for nonexistent and invalid XML files that cannot be parsed.
func TestLoadDataStream(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap")

	tests := []struct {
		dsPath  string
		wantErr bool
	}{
		{filepath.Join(testDataDir, "ssg-rhel-ds.xml"), false},
		{filepath.Join(testDataDir, "nonexistent.xml"), true},
		{filepath.Join(testDataDir, "invalid.xml"), true},
	}

	for _, tt := range tests {
		t.Run(tt.dsPath, func(t *testing.T) {
			_, err := loadDataStream(tt.dsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadDataStream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestGetDsProfileTitle tests the GetDsProfileTitle function.
// It also uses the getDsElement function with some additional logic specific for profile titles.
func TestGetDsProfileTitle(t *testing.T) {
	testDataDir := filepath.Join("..", "..", "..", "internal", "complytime", "testdata", "openscap")

	tests := []struct {
		profileId string
		dsPath    string
		expected  string
		wantErr   bool
	}{
		{"test_profile", filepath.Join(testDataDir, "ssg-rhel-ds.xml"), "Test Profile", false},
		{"test_profile_no_title", filepath.Join(testDataDir, "ssg-rhel-ds.xml"), "", false},
		{"nonexistent_profile", filepath.Join(testDataDir, "ssg-rhel-ds.xml"), "", true},
		{"invalid_profile", filepath.Join(testDataDir, "nonexistent.xml"), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.profileId, func(t *testing.T) {
			result, err := GetDsProfileTitle(tt.profileId, tt.dsPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDsProfileTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}
