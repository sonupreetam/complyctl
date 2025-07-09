// SPDX-License-Identifier: Apache-2.0

package complytime

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/models"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/oscal-compass/oscal-sdk-go/validation"
)

// WritePlan writes an AssessmentPlan to a given path location with consistency.
func WritePlan(plan *oscalTypes.AssessmentPlan, frameworkId string, planLocation string) error {
	// Ensure UserWorkspace exists before writing the plan
	userWorkspace := filepath.Dir(planLocation)
	if err := os.MkdirAll(userWorkspace, 0700); err != nil {
		return err
	}

	// Add the framework property needed for ComplyTime
	if plan.Metadata.Props == nil {
		plan.Metadata.Props = &[]oscalTypes.Property{}
	}
	frameworkProperty := oscalTypes.Property{
		Name:  extensions.FrameworkProp,
		Value: frameworkId,
		Ns:    extensions.TrestleNameSpace,
	}
	*plan.Metadata.Props = append(*plan.Metadata.Props, frameworkProperty)

	// To ensure we can easily read the plan once written, include under
	// OSCAL Model type to include the top-level "assessment-plan" key.
	oscalModels := oscalTypes.OscalModels{
		AssessmentPlan: plan,
	}
	assessmentPlanData, err := json.MarshalIndent(oscalModels, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(planLocation, assessmentPlanData, 0600)
}

// ReadPlan reads an assessment plans from a given file path.
func ReadPlan(assessmentPlanPath string, validator validation.Validator) (*oscalTypes.AssessmentPlan, error) {
	file, err := os.Open(assessmentPlanPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	plan, err := models.NewAssessmentPlan(file, validator)
	if err != nil {
		return nil, fmt.Errorf("failed to load assessment plan from %s: %w", assessmentPlanPath, err)
	}
	return plan, nil
}

var ErrNoActivities = errors.New("no local activities detected")

// Settings return a new compliance Settings instance based on the
// given assessment plan path.
func Settings(plan *oscalTypes.AssessmentPlan) (settings.Settings, error) {
	if plan.LocalDefinitions != nil && plan.LocalDefinitions.Activities != nil {
		return settings.NewAssessmentActivitiesSettings(*plan.LocalDefinitions.Activities), nil
	}
	return settings.Settings{}, ErrNoActivities
}

// GetControlTitle retrieves the title for a control from the catalog
func GetControlTitle(controlID string, controlSource string, appDir ApplicationDirectory, validator validation.Validator) (string, error) {
	profile, err := LoadProfile(appDir, controlSource, validator)
	if err != nil {
		return "", fmt.Errorf("failed to load profile from source '%s': %w", controlSource, err)
	}

	if profile.Imports == nil {
		return "", fmt.Errorf("profile '%s' has no imports", controlSource)
	}

	for _, imp := range profile.Imports {
		catalog, err := LoadCatalogSource(appDir, imp.Href, validator)
		if err != nil {
			continue
		}
		if catalog.Groups == nil {
			continue
		}
		for _, group := range *catalog.Groups {
			if group.Controls == nil {
				continue
			}
			for _, control := range *group.Controls {
				if control.ID == controlID && control.Title != "" {
					return control.Title, nil
				}
			}
		}
	}
	return "", fmt.Errorf("title for control '%s' not found in catalog", controlID)
}
