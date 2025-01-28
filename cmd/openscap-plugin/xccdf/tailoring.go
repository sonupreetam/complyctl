// SPDX-License-Identifier: Apache-2.0

package xccdf

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/ComplianceAsCode/compliance-operator/pkg/xccdf"
	"github.com/oscal-compass/compliance-to-policy-go/v2/policy"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
)

const (
	XCCDFNamespace       string = "complytime.openscapplugin"
	XCCDFTailoringSuffix string = "complytime-tailoring-profile"
)

func getTailoringID() string {
	return fmt.Sprintf("xccdf_%s_tailoring_%s", XCCDFNamespace, XCCDFTailoringSuffix)
}

func getTailoringProfileID(profileId string) string {
	return fmt.Sprintf(
		"xccdf_%s_profile_%s_%s", XCCDFNamespace, profileId, XCCDFTailoringSuffix)
}

func getTailoringProfileTitle(dsProfileTitle string) string {
	return fmt.Sprintf("ComplyTime Tailoring Profile - %s", dsProfileTitle)
}

func getTailoringVersion() xccdf.VersionElement {
	return xccdf.VersionElement{
		Time:  time.Now().Format(time.RFC3339),
		Value: "1",
	}
}

func getTailoringBenchmarkHref(datastreamPath string) xccdf.BenchmarkElement {
	return xccdf.BenchmarkElement{
		Href: datastreamPath,
	}
}

func validateRuleExistence(policyRuleID string, dsRules []DsRules) bool {
	for _, dsRule := range dsRules {
		ruleID := xccdf.GetRuleNameFromID(dsRule.ID)
		if policyRuleID == ruleID {
			return true
		}
	}
	return false
}

func unselectAbsentRules(tailoringSelections, dsProfileSelections []xccdf.SelectElement, oscalPolicy policy.Policy) []xccdf.SelectElement {
	for _, dsRule := range dsProfileSelections {
		dsRuleAlsoInPolicy := false
		ruleID := xccdf.GetRuleNameFromID(dsRule.IDRef)
		for _, rule := range oscalPolicy {
			if ruleID == rule.Rule.ID {
				dsRuleAlsoInPolicy = true
				break
			}
		}
		if !dsRuleAlsoInPolicy {
			tailoringSelections = append(tailoringSelections, xccdf.SelectElement{
				IDRef:    dsRule.IDRef,
				Selected: false,
			})
		}
	}
	return tailoringSelections
}

func selectAdditionalRules(tailoringSelections, dsProfileSelections []xccdf.SelectElement, oscalPolicy policy.Policy) []xccdf.SelectElement {
	for _, rule := range oscalPolicy {
		ruleAlreadyInDsProfile := false
		for _, dsRule := range dsProfileSelections {
			ruleID := xccdf.GetRuleNameFromID(dsRule.IDRef)
			if rule.Rule.ID == ruleID {
				// Not a common case, but a rule be be unselected in Datastream Profile
				if dsRule.Selected {
					ruleAlreadyInDsProfile = true
				}
				break
			}
		}

		if !ruleAlreadyInDsProfile {
			tailoringSelections = append(tailoringSelections, xccdf.SelectElement{
				IDRef:    getDsRuleID(rule.Rule.ID),
				Selected: true,
			})
		}
	}
	return tailoringSelections
}

func getTailoringSelections(oscalPolicy policy.Policy, dsProfile *xccdf.ProfileElement, dsPath string) ([]xccdf.SelectElement, error) {
	if oscalPolicy == nil {
		return nil, fmt.Errorf("OSCAL policy is empty")
	}

	dsRules, err := GetDsRules(dsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get rules from datastream: %w", err)
	}

	// All policy rules should be present in the Datastream
	for _, rule := range oscalPolicy {
		if !validateRuleExistence(rule.Rule.ID, dsRules) {
			return nil, fmt.Errorf("rule not found in datastream: %s", rule.Rule.ID)
		}
	}

	var tailoringSelections []xccdf.SelectElement
	// Rules in dsProfile but not in OSCAL Policy must be unselected in Tailoring file.
	tailoringSelections = unselectAbsentRules(tailoringSelections, dsProfile.Selections, oscalPolicy)
	tailoringSelections = selectAdditionalRules(tailoringSelections, dsProfile.Selections, oscalPolicy)

	return tailoringSelections, nil
}

func getTailoringProfile(profileId string, dsPath string, oscalPolicy policy.Policy) (*xccdf.ProfileElement, error) {
	tailoringProfile := new(xccdf.ProfileElement)
	tailoringProfile.ID = getTailoringProfileID(profileId)

	dsProfile, err := GetDsProfile(profileId, dsPath)
	if err != nil {
		return tailoringProfile, fmt.Errorf("failed to get base profile from datastream: %w", err)
	}

	tailoringProfile.Title = &xccdf.TitleOrDescriptionElement{
		Override: true,
		Value:    getTailoringProfileTitle(dsProfile.Title.Value),
	}

	tailoringProfile.Selections, err = getTailoringSelections(oscalPolicy, dsProfile, dsPath)
	if err != nil {
		return tailoringProfile, fmt.Errorf("failed to get selections for tailoring profile: %w", err)
	}

	//	tailoringProfile.Values = getTailoringValues(oscalPolicy, dsProfile, dsPath)
	return tailoringProfile, nil
}

func PolicyToXML(oscalPolicy policy.Policy, config *config.Config) (string, error) {
	datastreamPath := config.Files.Datastream
	profileId := config.Parameters.Profile

	tailoringProfile, err := getTailoringProfile(profileId, datastreamPath, oscalPolicy)
	if err != nil {
		return "", err
	}

	tailoring := xccdf.TailoringElement{
		XMLNamespaceURI: xccdf.XCCDFURI,
		ID:              getTailoringID(),
		Version:         getTailoringVersion(),
		Benchmark:       getTailoringBenchmarkHref(datastreamPath),
		Profile:         *tailoringProfile,
	}

	output, err := xml.MarshalIndent(tailoring, "", "  ")
	if err != nil {
		return "", err
	}
	return xccdf.XMLHeader + "\n" + string(output), nil
}

/*
func getValuesFromPolicyVariables(tp policy.Policy) []xccdf.SetValueElement {
	var values []xccdf.SetValueElement
	if len(tp) != 0 {
		for _, rule := range tp {
			if rule.Rule.Parameter == nil {
				continue
			}

			values = append(values, xccdf.SetValueElement{
				IDRef: rule.Rule.Parameter.ID,
				Value: rule.Rule.Parameter.Value,
			})
		}
		return values
	}

	return nil
}
*/
