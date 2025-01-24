// SPDX-License-Identifier: Apache-2.0

package xccdf

import (
	"fmt"
	"os"

	"github.com/ComplianceAsCode/compliance-operator/pkg/xccdf"
	"github.com/antchfx/xmlquery"
)

func getDsProfileID(profileId string) string {
	return fmt.Sprintf("xccdf_org.ssgproject.content_profile_%s", profileId)
}

func getDsElement(dsDom *xmlquery.Node, dsElement string) (*xmlquery.Node, error) {
	// Returns nil if the element is not found
	return xmlquery.Query(dsDom, dsElement)
}

func loadDataStream(dsPath string) (*xmlquery.Node, error) {
	file, err := os.Open(dsPath)
	if err != nil {
		return nil, fmt.Errorf("error opening datastream file: %s", err)
	}
	defer file.Close()

	dsDom, err := xmlquery.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("error parsing datastream file: %s", err)
	}

	return dsDom, nil
}

func getDsProfile(dsDom *xmlquery.Node, dsProfileID string) (*xmlquery.Node, error) {
	return getDsElement(dsDom, fmt.Sprintf("//xccdf-1.2:Profile[@id='%s']", dsProfileID))
}

func getDsProfileTitle(dsProfile *xmlquery.Node) (*xmlquery.Node, error) {
	profileTitle, err := getDsElement(dsProfile, "xccdf-1.2:title")
	if err != nil {
		return nil, fmt.Errorf("error finding title element in profile: %s", err)
	}
	return profileTitle, nil
}

func getDsProfileDescription(dsProfile *xmlquery.Node) (*xmlquery.Node, error) {
	profileDescription, err := getDsElement(dsProfile, "xccdf-1.2:description")
	if err != nil || profileDescription == nil {
		return nil, fmt.Errorf("error finding description element in profile: %s", err)
	}
	return profileDescription, nil
}

func populateProfileInfo(dsProfile *xmlquery.Node, parsedProfile *xccdf.ProfileElement) (*xccdf.ProfileElement, error) {
	profileTitle, err := getDsProfileTitle(dsProfile)
	if err != nil {
		return parsedProfile, fmt.Errorf("error populating profile title: %s", err)
	}
	if parsedProfile.Title == nil {
		parsedProfile.Title = &xccdf.TitleOrDescriptionElement{}
	}
	if profileTitle == nil {
		// log that profile title was not found.
		// It is a valid case therefore but better to log it.
		parsedProfile.Title.Override = false
		parsedProfile.Title.Value = ""
	} else {
		parsedProfile.Title.Override = true
		parsedProfile.Title.Value = profileTitle.InnerText()
	}

	profileDescription, err := getDsProfileDescription(dsProfile)
	if err != nil {
		return parsedProfile, fmt.Errorf("error populating profile description: %s", err)
	}
	if parsedProfile.Description == nil {
		parsedProfile.Description = &xccdf.TitleOrDescriptionElement{}
	}
	if profileDescription == nil {
		// log that profile description was not found.
		parsedProfile.Description.Override = false
		parsedProfile.Description.Value = ""
	} else {
		parsedProfile.Description.Override = true
		parsedProfile.Description.Value = profileDescription.InnerText()
	}

	return parsedProfile, nil
}

func initProfile(dsProfile *xmlquery.Node, dsProfileId string) (*xccdf.ProfileElement, error) {
	parsedProfile := new(xccdf.ProfileElement)
	parsedProfile.ID = dsProfileId

	parsedProfile, err := populateProfileInfo(dsProfile, parsedProfile)
	if err != nil {
		return parsedProfile, fmt.Errorf("error populating profile title and description: %s", err)
	}

	return parsedProfile, nil
}

func GetDsProfileTitle(profileId string, dsPath string) (string, error) {
	profile, err := GetDsProfile(profileId, dsPath)
	if err != nil {
		return "", fmt.Errorf("error processing profile %s in datastream: %s", profileId, err)
	}

	return profile.Title.Value, nil
}

func GetDsProfile(profileId string, dsPath string) (*xccdf.ProfileElement, error) {
	dsDom, err := loadDataStream(dsPath)
	if err != nil {
		return nil, fmt.Errorf("error loading datastream: %s", err)
	}

	dsProfileID := getDsProfileID(profileId)
	dsProfile, err := getDsProfile(dsDom, dsProfileID)
	if err != nil {
		return nil, fmt.Errorf("error processing profile %s in datastream: %s", profileId, err)
	}

	if dsProfile == nil {
		return nil, fmt.Errorf("profile not found: %s", dsProfileID)
	}

	parsedProfile, err := initProfile(dsProfile, dsProfileID)
	if err != nil {
		return nil, fmt.Errorf("error parsing profile %s in datastream: %s", profileId, err)
	}

	return parsedProfile, nil
}

// Getting rule information
// Copied from https://github.com/ComplianceAsCode/compliance-operator/blob/fed54b4b761374578016d79d97bcb7636bf9d920/pkg/utils/parse_arf_result.go#L170

func NewRuleHashTable(dsDom *xmlquery.Node) NodeByIdHashTable {
	return newHashTableFromRootAndQuery(dsDom, "//ds:component/xccdf-1.2:Benchmark", "//xccdf-1.2:Rule")
}

func newHashTableFromRootAndQuery(dsDom *xmlquery.Node, root, query string) NodeByIdHashTable {
	benchmarkDom := dsDom.SelectElement(root)
	rules := benchmarkDom.SelectElements(query)
	return newByIdHashTable(rules)
}

type NodeByIdHashTable map[string]*xmlquery.Node

//type nodeByIdHashVariablesTable map[string][]string

func newByIdHashTable(nodes []*xmlquery.Node) NodeByIdHashTable {
	table := make(NodeByIdHashTable)
	for i := range nodes {
		ruleDefinition := nodes[i]
		ruleId := ruleDefinition.SelectAttr("id")

		table[ruleId] = ruleDefinition
	}

	return table
}
