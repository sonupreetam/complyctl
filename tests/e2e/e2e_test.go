//go:build e2e

package e2e

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	FrameworkID = "cusp_fedora"
	controlID   = "cusp_fedora_1-1"
	ruleID      = "file_groupowner_grub2_cfg"
	parameterID = "var_rekey_limit_size"
	configPath  = "../testdata/config.yml"
)

func TestComplyctlHelp(t *testing.T) {
	// Run the "complyctl --help" command
	cmd := exec.Command("complyctl", "--help")
	output, err := cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl --help: %v\nOutput: %s", err, string(output))
	}

	// Convert the output to a string and check if expected text is present
	outputStr := string(output)

	// Assert that "Usage" or the expected help message is part of the output
	// Use a table-driven test for clear, maintainable assertions
	expectedSubstrings := []string{
		"Usage:",
		"Aliases:",
		"Available Commands:",
		"Flags:",
		"complyctl [command]",
	}

	for _, expected := range expectedSubstrings {
		t.Run("Contains "+expected, func(t *testing.T) {
			assert.True(t, strings.Contains(outputStr, expected), "Help output should contain '%s'", expected)
		})
	}
}

func TestComplyctlList(t *testing.T) {
	// Run the "complyctl list" command
	cmd := exec.Command("complyctl", "list", "--plain")
	output, err := cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl list: %v\nOutput: %s", err, string(output))
	}

	// Convert the output to a string and check if expected content is returned
	outputStr := string(output)

	// Check if the output contains expected text
	assert.True(t, len(outputStr) > 0, "Output from 'complyctl list' should not be empty")
	assert.True(t, strings.Contains(outputStr, FrameworkID))
}

func TestComplyctlInfo(t *testing.T) {
	testCases := []struct {
		name             string   // A descriptive name for the sub-test.
		args             []string // The command-line arguments to pass to "complyctl info".
		expectedFragment string   // A string we expect to find in the command's output.
	}{
		{
			name:             "show info for FrameworkID",
			args:             []string{"info", FrameworkID, "--plain"},
			expectedFragment: FrameworkID,
		},
		{
			name:             "show info for controlID",
			args:             []string{"info", FrameworkID, "--control", controlID, "--plain"},
			expectedFragment: controlID,
		},
		{
			name:             "show info for rule",
			args:             []string{"info", FrameworkID, "--rule", ruleID, "--plain"},
			expectedFragment: ruleID,
		},
		{
			name:             "show info for parameter",
			args:             []string{"info", FrameworkID, "--parameter", parameterID, "--plain"},
			expectedFragment: parameterID,
		},
	}

	// Iterate over the test cases and run each as a sub-test.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Execute the command
			cmd := exec.Command("complyctl", tc.args...)
			output, err := cmd.CombinedOutput()

			// It stops the current sub-test on failure.
			require.NoError(t, err, "command failed unexpectedly. Output: \n%s", string(output))

			// Use assert.Contains for the actual test assertion.
			assert.Contains(t, string(output), tc.expectedFragment, "output should contain the expected fragment")
		})
	}
}

func TestComplyctlPlan(t *testing.T) {
	// testCases defines all scenarios we want to test for the 'plan' command.
	testCases := []struct {
		name           string   // A descriptive name for the sub-test.
		args           []string // Arguments to pass to the 'complyctl plan' command.
		expectedOutput string   // A substring we expect to find in the command's output.
		expectedFile   string   // The path to a file that should be created.
		expectError    bool     // Whether we expect the command to fail.
	}{
		{
			name:           "DefaultWorkspace",
			args:           []string{FrameworkID},
			expectedOutput: "INFO Assessment plan written to complytime/assessment-plan.json",
			expectedFile:   "complytime/assessment-plan.json",
		},
		{
			name:           "CustomWorkspace",
			args:           []string{FrameworkID, "--workspace", "custom_dir"},
			expectedOutput: "INFO Assessment plan written to custom_dir/assessment-plan.json",
			expectedFile:   "custom_dir/assessment-plan.json",
		},
		{
			name:           "DryRun",
			args:           []string{FrameworkID, "--dry-run"},
			expectedOutput: "controlId:", // Check if the output contains expected text
			expectedFile:   "",           // No file should be created.
		},
		{
			name:           "DryRunOutFile",
			args:           []string{FrameworkID, "--dry-run", "--out", "config.yml"},
			expectedOutput: "", // No output, so we don't check it.
			expectedFile:   "config.yml",
		},
		{
			name:           "ScopeConfig",
			args:           []string{FrameworkID, "--scope-config", configPath, "--workspace", "config_dir"},
			expectedOutput: "INFO Assessment plan written to config_dir/assessment-plan.json",
			expectedFile:   "config_dir/assessment-plan.json",
		},
	}

	// Loop through all defined test cases and run them as sub-tests.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Construct the full command arguments.
			args := append([]string{"plan"}, tc.args...)
			cmd := exec.Command("complyctl", args...)

			// Execute the command.
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Assert that the command's success or failure matches our expectation.
			if tc.expectError {
				require.Error(t, err, "Expected an error but got none. Output:\n%s", outputStr)
			} else {
				require.NoError(t, err, "Command failed unexpectedly. Output:\n%s", outputStr)
			}

			// If we expect specific output, check for it.
			if tc.expectedOutput != "" {
				assert.Contains(t, outputStr, tc.expectedOutput, "Output did not contain expected text")
			}

			// If we expect a file to be created, verify it exists.
			if tc.expectedFile != "" {
				assert.FileExists(t, tc.expectedFile)
			}
		})
	}
}

func TestComplyctlGenerate(t *testing.T) {

	cmd := exec.Command("complyctl", "generate")
	output, err := cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl generate: %v\nOutput: %s", err, string(output))
	}

	// Check if the output contains expected text
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "Policy generation process completed"))
}

func TestComplyctlScan(t *testing.T) {
	// Run the "complyctl scan" command
	// In order to improve the performace, without md will be covered in the CustomizePlanWorkflow
	cmd := exec.Command("complyctl", "scan", "--with-md")
	output, err := cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl scan --with-md: %v\nOutput: %s", err, string(output))
	}

	// Check if the output contains expected text
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "The assessment results in JSON were successfully written to complytime/assessment-results.json"))
	assert.True(t, strings.Contains(outputStr, "The assessment results in markdown were successfully written to complytime/assessment-results.md"))

	// Check if the assessment-results files exist
	filePath_json := "complytime/assessment-results.json"
	filePath_md := "complytime/assessment-results.md"
	assert.FileExists(t, filePath_json, "Expected file exists:", filePath_json)
	assert.FileExists(t, filePath_md, "Expected file exists:", filePath_md)
}

func TestComplyctlCustomizePlanWorkflow(t *testing.T) {
	// The the workflow of customize the assessment plan via the config.yml

	// 1. Load config.yml to customize the generated assessment plan
	cmd := exec.Command("complyctl", "plan", FrameworkID, "--scope-config", configPath)
	output, err := cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl plan: %v\nOutput: %s", err, string(output))
	}
	// Check if the output contains expected text
	outputStr := string(output)
	assert.True(t, strings.Contains(outputStr, "Assessment plan written to complytime/assessment-plan.json"))

	// 2. Generate PVP policy from the assessment plan
	cmd = exec.Command("complyctl", "generate")
	output, err = cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl generate: %v\nOutput: %s", err, string(output))
	}
	// Check if the output contains expected text
	outputStr = string(output)
	assert.True(t, strings.Contains(outputStr, "Policy generation process completed"))

	// 3. Scan environment with the customized assessment plan
	cmd = exec.Command("complyctl", "scan")
	output, err = cmd.CombinedOutput()

	// Ensure there is no error when running the command
	if err != nil {
		t.Fatalf("Error running complyctl scan: %v\nOutput: %s", err, string(output))
	}
	// Check if the output contains expected text
	outputStr = string(output)
	assert.True(t, strings.Contains(outputStr, "The assessment results in JSON were successfully written to complytime/assessment-results.json"))
	assert.True(t, strings.Contains(outputStr, "No assessment result in markdown will be generated"))
}
