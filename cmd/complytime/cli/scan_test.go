// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAssessmentPlan(t *testing.T) {

	testFile := filepath.Join("./testdata", "assessment-plan.json")

	assessmentPlan, err := loadAssessmentPlan(testFile)
	require.NoError(t, err)
	require.Equal(t, assessmentPlan.Metadata.Title, "REPLACE_ME")
}
