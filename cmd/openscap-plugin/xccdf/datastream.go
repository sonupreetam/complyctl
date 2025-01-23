// SPDX-License-Identifier: Apache-2.0

package xccdf

import (
	"fmt"
	"os"

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

func GetDsProfileTitle(profileId string, dsPath string) (string, error) {
	dsDom, err := loadDataStream(dsPath)
	if err != nil {
		return "", fmt.Errorf("error loading datastream: %s", err)
	}

	dsProfileID := getDsProfileID(profileId)
	profile, err := getDsElement(dsDom, fmt.Sprintf("//xccdf-1.2:Profile[@id='%s']", dsProfileID))
	if err != nil {
		return "", fmt.Errorf("error processing profile %s in datastream: %s", dsProfileID, err)
	}

	if profile == nil {
		return "", fmt.Errorf("profile not found: %s", dsProfileID)
	}

	profileTitle, err := xmlquery.Query(profile, "xccdf-1.2:title")
	if err != nil || profileTitle == nil {
		return "", fmt.Errorf("error finding title element in profile %s: %s", dsProfileID, err)
	}
	return profileTitle.InnerText(), nil
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
