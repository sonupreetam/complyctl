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

func getTailoringProfileTitle(profileId string, dsPath string) string {
	dsProfileTitle, err := GetDsProfileTitle(profileId, dsPath)
	if err != nil || dsProfileTitle == "" {
		// log that the profile title was not found or is empty in Datastream
		dsProfileTitle = profileId
	}
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

func getTailoringProfile(
	profileId string, tailoringPolicy policy.Policy, dsPath string) xccdf.ProfileElement {
	return xccdf.ProfileElement{
		ID: getTailoringProfileID(profileId),
		Title: &xccdf.TitleOrDescriptionElement{
			Value: getTailoringProfileTitle(profileId, dsPath),
		},
		// TODO: Should include only the diff from the original profile in Datastream
		Selections: getPolicySelections(tailoringPolicy),
		Values:     getValuesFromPolicyVariables(tailoringPolicy),
	}
}

func getPolicySelections(tailoringPolicy policy.Policy) []xccdf.SelectElement {
	var selections []xccdf.SelectElement
	for _, rule := range tailoringPolicy {
		selections = append(selections, xccdf.SelectElement{
			IDRef:    rule.Rule.ID,
			Selected: true,
		})
	}
	return selections
}

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

func PolicyToXML(tailoringPolicy policy.Policy, config *config.Config) (string, error) {
	datastreamPath := config.Files.Datastream
	profileId := config.Parameters.Profile

	tailoring := xccdf.TailoringElement{
		XMLNamespaceURI: xccdf.XCCDFURI,
		ID:              getTailoringID(),
		Version:         getTailoringVersion(),
		Benchmark:       getTailoringBenchmarkHref(datastreamPath),
		Profile:         getTailoringProfile(profileId, tailoringPolicy, datastreamPath),
	}

	output, err := xml.MarshalIndent(tailoring, "", "  ")
	if err != nil {
		return "", err
	}
	return xccdf.XMLHeader + "\n" + string(output), nil
}
