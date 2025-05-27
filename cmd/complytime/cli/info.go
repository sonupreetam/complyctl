// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/oscal-compass/oscal-sdk-go/models/components"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/oscal-compass/oscal-sdk-go/validation"

	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
	"github.com/complytime/complytime/internal/terminal"
)

const (
	defaultTableLimit      = 10
	colWidthControlID      = 12
	colWidthControlTitle   = 40
	colWidthImplStatus     = 21
	colWidthPluginsUsed    = 24
	colWidthRulesInControl = 55
	colWidthPluginUsed     = 31
)

var (
	tableHeaderStyle = table.DefaultStyles().Header.
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")). // Light grey border
				BorderBottom(true).
				Bold(false)

	tableCellStyle = table.DefaultStyles().Cell.
			Foreground(lipgloss.Color("250")). // Very light grey text
			Bold(true)

	// Style for key-value pair keys
	keyStyle = lipgloss.NewStyle().
			PaddingRight(1)

	// Style for key-value pairs values
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			Bold(true)

	// Container style for key-value information blocks
	infoContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240")). // Light grey border for containers
				Padding(0, 1).
				Width(90) // Fixed width for consistent formatting
)

// Rule represents details about a rule for easy
// mapping of rules to plugins.
type Rule struct {
	ID     string
	Plugin string
}

// Control repsents details about a control across
// component sources.
type Control struct {
	ID                  string
	Title               string
	Description         string
	ImplemenationStatus string
	Rules               []Rule
}

// RuleInfo represents details about a rule
type RuleInfo struct {
	ID               string
	Description      string
	CheckID          string
	CheckDescription string
	Parameters       []string
}

type infoOptions struct {
	*option.Common
	complyTimeOpts *option.ComplyTime
	controlID      string // show info for a specific control ID
	ruleID         string // show info for a specific rule ID
	limit          int    // limit number for table rows shown in terminal
	plain          bool   // print plain table only
}

func infoCmd(common *option.Common) *cobra.Command {
	infoOpts := &infoOptions{
		Common:         common,
		complyTimeOpts: &option.ComplyTime{},
	}
	cmd := &cobra.Command{
		Use:     "info <framework-id> [flags]",
		Short:   "Import OSCAL data from API",
		Example: " complytime info anssi_bp28_minimal\n complytime info anssi_bp28_minimal --control r31\n complytime info anssi_bp28_minimal --rule enable_authselect",
		Args:    cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				infoOpts.complyTimeOpts.FrameworkID = filepath.Clean(args[0])
			}
		},
		RunE: func(_ *cobra.Command, _ []string) error { return runInfo(infoOpts) },
	}
	cmd.Flags().StringVarP(&infoOpts.controlID, "control", "c", "", "show info for a specific control ID")
	cmd.Flags().StringVarP(&infoOpts.ruleID, "rule", "r", "", "show info for a specific rule ID")
	cmd.Flags().IntVarP(&infoOpts.limit, "limit", "l", 10, fmt.Sprintf("limit the number of table rows (default is %d)", defaultTableLimit))
	cmd.Flags().BoolVarP(&infoOpts.plain, "plain", "p", false, "print the table with minimal formatting")
	infoOpts.complyTimeOpts.BindFlags(cmd.Flags())
	return cmd
}

// runInfo executes the info command using the provided optoions.
func runInfo(opts *infoOptions) error {

	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return fmt.Errorf("failed to initialize application directory: %w", err)
	}
	logger.Debug(fmt.Sprintf("Using application directory: %s", appDir.AppDir()))

	validator := validation.NewSchemaValidator()

	compDefs, err := complytime.FindComponentDefinitions(appDir.BundleDir(), validator)
	if err != nil {
		return fmt.Errorf("failed to find component definitions: %w", err)
	}

	components := filterFrameworkComponents(compDefs, opts.complyTimeOpts.FrameworkID)
	if len(components) == 0 {
		return fmt.Errorf("no components found for framework ID '%s'", opts.complyTimeOpts.FrameworkID)
	}

	// Use maps for easy lookups of data extracted from component definitions
	ruleRemarksMap, remarksPropsMap := processComponentProperties(components)
	controlMap, setParameters := processControlImplementations(components, appDir, validator)

	// Display info based on controlID or ruleID flag being passed at CLI
	if opts.controlID != "" {
		return displayControlInfo(opts, controlMap)
	} else if opts.ruleID != "" {
		return displayRuleInfo(opts, opts.ruleID, ruleRemarksMap, remarksPropsMap, setParameters)
	} else {
		return displayAllControls(opts, controlMap)
	}

}

// processControlImplementations extracts control details and set parameters from component definitions.
func processControlImplementations(compDefs []oscalTypes.DefinedComponent, appDir complytime.ApplicationDirectory, validator *validation.SchemaValidator) (
	map[string]Control,
	map[string][]string, // setParameters: maps param ID to its values
) {
	controlMap := make(map[string]Control)
	setParameters := make(map[string][]string)

	for _, comp := range compDefs {
		pluginForComponent := comp.Title // Using Title as plugin name

		if comp.ControlImplementations == nil {
			continue
		}

		for _, ci := range *comp.ControlImplementations {
			// Process implemented requirements
			if ci.ImplementedRequirements != nil {
				for _, ir := range ci.ImplementedRequirements {
					controlDetails, ok := controlMap[ir.ControlId]
					if !ok {
						// Initialize controlDetails if not already present
						controlTitle, err := getControlTitle(ir.ControlId, ci, appDir, validator)
						if err != nil {
							logger.Warn("could not get title for control %s: %v", ir.ControlId, err)
							controlTitle = "N/A"
						}

						controlDetails = Control{
							ID:          ir.ControlId,
							Title:       controlTitle,
							Description: ir.Description,
							Rules:       []Rule{},
						}
					}

					if ir.Props != nil {
						for _, p := range *ir.Props {
							if p.Name == "Rule_Id" && p.Value != "" {
								controlDetails.Rules = append(controlDetails.Rules, Rule{ID: p.Value, Plugin: pluginForComponent})
							} else if p.Name == "implementation-status" && p.Value != "" {
								controlDetails.ImplemenationStatus = p.Value
							}
						}
					}
					controlMap[ir.ControlId] = controlDetails // Update map with potentially modified control
				}
			}

			// Process set parameters
			if ci.SetParameters != nil {
				for _, sp := range *ci.SetParameters {
					if sp.ParamId != "" && len(sp.Values) > 0 {
						setParameters[sp.ParamId] = sp.Values
					}
				}
			}
		}
	}
	return controlMap, setParameters
}

// processComponentProperties extracts rule and property information from
// component definitions info a maps for easy lookups.
func processComponentProperties(compDefs []oscalTypes.DefinedComponent) (map[string]string, map[string][]oscalTypes.Property) {
	ruleRemarksMap := make(map[string]string)                 // maps rule ID to its remarks value
	remarksPropsMap := make(map[string][]oscalTypes.Property) // maps remarks value to list of properties

	for _, comp := range compDefs {
		if comp.Props == nil {
			continue
		}
		for _, prop := range *comp.Props {
			remarksPropsMap[prop.Remarks] = append(remarksPropsMap[prop.Remarks], prop)
			// only populate ruleRemarksMap if both Value and Remarks are present
			if prop.Name == "Rule_Id" && prop.Value != "" && prop.Remarks != "" {
				ruleRemarksMap[prop.Value] = prop.Remarks
			}
		}
	}
	return ruleRemarksMap, remarksPropsMap
}

// getControlTitle retrieves the title for a given control ID from the associated catalog.
func getControlTitle(controlID string, controlImplementation oscalTypes.ControlImplementationSet, appDir complytime.ApplicationDirectory, validator *validation.SchemaValidator) (string, error) {
	profile, err := complytime.LoadProfile(appDir, controlImplementation.Source, validator)
	if err != nil {
		return "", fmt.Errorf("failed to load profile from source '%s': %w", controlImplementation.Source, err)
	}

	if profile.Imports == nil {
		return "", fmt.Errorf("profile '%s' has no imports", controlImplementation.Source)
	}

	for _, imp := range profile.Imports {
		catalog, err := complytime.LoadCatalogSource(appDir, imp.Href, validator)
		if err != nil {
			logger.Warn("failed to load catalog from href '%s': %v", imp.Href, err)
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

// filterFrameworkComponents filters component definitions by framework ID.
func filterFrameworkComponents(componentDefinitions []oscalTypes.ComponentDefinition, frameworkID string) []oscalTypes.DefinedComponent {
	var validationComponents []oscalTypes.DefinedComponent

	for _, compDef := range componentDefinitions {
		if compDef.Components == nil {
			continue
		}

		for _, component := range *compDef.Components {
			if component.Type == string(components.Validation) {
				if component.ControlImplementations == nil {
					continue
				}
				for _, controlImp := range *component.ControlImplementations {
					framework, ok := settings.GetFrameworkShortName(controlImp)
					if ok && framework == frameworkID {
						validationComponents = append(validationComponents, component)
						break // Component belongs to this framework, move to next compDef
					}
				}
			}
		}
	}
	return validationComponents
}

// extractRuleDetails parses properties to fill a RuleInfo struct.
func extractRuleDetails(props []oscalTypes.Property) RuleInfo {
	info := RuleInfo{}
	for _, prop := range props {
		switch prop.Name {
		case "Rule_Description":
			info.Description = prop.Value
		case "Check_Id":
			info.CheckID = prop.Value
		case "Check_Description":
			info.CheckDescription = prop.Value
		case "Parameter_Id":
			info.Parameters = append(info.Parameters, prop.Value)
		}
	}
	return info
}

// runBubbleTeaProgram executes the model
func runBubbleTeaProgram(model terminal.Model, output io.Writer) error {

	if _, err := tea.NewProgram(model, tea.WithOutput(output)).Run(); err != nil {
		return fmt.Errorf("failed to display information: %w", err)
	}
	return nil
}

// removeDuplicates removes duplicate items from a slice
func removeDuplicates[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

// renderKeyValuePair renders a single key-value pair using predefined lipgloss styles.
func renderKeyValuePair(key, value string) string {
	return keyStyle.Render(key) + ": " + valueStyle.Render(value)
}

// getControlRulesColumnsAndRows prepares columns and rows for a specific control's rules table.
func getControlRulesColumnsAndRows(control Control) ([]table.Column, []table.Row) {
	// Sort rules by ID for logical ordering in the table
	sort.Slice(control.Rules, func(i, j int) bool {
		return control.Rules[i].ID < control.Rules[j].ID
	})

	var rows []table.Row
	for _, rule := range control.Rules {
		row := table.Row{rule.ID, rule.Plugin}
		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "Rules In Control", Width: colWidthRulesInControl},
		{Title: "Plugin Used", Width: colWidthPluginUsed},
	}

	// Calculate dynamic column width based on content.
	for i := range columns {
		maxLength := columns[i].Width
		for _, row := range rows {
			if i < len(row) {
				cellLength := lipgloss.Width(row[i])
				if cellLength > maxLength {
					maxLength = cellLength
				}
			}
		}
		columns[i].Width = maxLength
	}

	return columns, rows
}

// getRuleParametersColumnsAndRows prepares columns and rows for the rule parameters table.
func getRuleParametersColumnsAndRows(ruleDetails RuleInfo, setParameters map[string][]string) ([]table.Column, []table.Row) {
	var rows []table.Row

	// Sort parameters by ID for consistent ordering in the table
	sort.Slice(ruleDetails.Parameters, func(i, j int) bool {
		return ruleDetails.Parameters[i] < ruleDetails.Parameters[j]
	})

	for _, paramID := range ruleDetails.Parameters {
		paramValue := "N/A"
		if values, ok := setParameters[paramID]; ok && len(values) > 0 {
			paramValue = strings.Join(values, ", ")
		}
		rows = append(rows, table.Row{paramID, paramValue})
	}

	columns := []table.Column{
		{Title: "Parameter ID", Width: 56},
		{Title: "Set Value(s)", Width: 30},
	}

	// Calculate dynamic column width based on content
	for i := range columns {
		maxLength := columns[i].Width
		for _, row := range rows {
			if i < len(row) {
				cellLength := lipgloss.Width(row[i])
				if cellLength > maxLength {
					maxLength = cellLength
				}
			}
		}
		columns[i].Width = maxLength
	}

	return columns, rows
}

func getControlListColumnsAndRows(controls []Control) ([]table.Column, []table.Row) {
	// Sort controls by ID for logical ordering in the table
	sort.Slice(controls, func(i, j int) bool {
		return controls[i].ID < controls[j].ID
	})

	var rows []table.Row
	for _, control := range controls {
		var plugins []string
		for _, rule := range control.Rules {
			plugins = append(plugins, rule.Plugin)
		}

		row := table.Row{
			control.ID,
			control.Title,
			control.ImplemenationStatus,
			strings.Join(removeDuplicates(plugins), ", "),
		}
		rows = append(rows, row)
	}

	columns := []table.Column{
		{Title: "Control ID", Width: colWidthControlID},
		{Title: "Control Title", Width: colWidthControlTitle},
		{Title: "Implementation Status", Width: colWidthImplStatus},
		{Title: "Plugins Used", Width: colWidthPluginsUsed},
	}

	// Calculate dynamic column width based on content.
	for i := range columns {
		maxLength := columns[i].Width // Start with defined width
		for _, row := range rows {
			if i < len(row) {
				cellLength := lipgloss.Width(row[i])
				if cellLength > maxLength {
					maxLength = cellLength
				}
			}
		}
		columns[i].Width = maxLength
	}
	return columns, rows
}

// newControlInfoModel creates a Tea model for displaying specific control details.
func newControlInfoModel(control Control, rowLimit int) terminal.Model {
	// Prepare the header message with control details
	wrappedDescription := terminal.WrapText(control.Description, 60)
	headerFields := strings.Join([]string{
		renderKeyValuePair("Control ID", control.ID),
		renderKeyValuePair("Title", control.Title),
		renderKeyValuePair("Implementation Status", control.ImplemenationStatus),
		keyStyle.Render("Description") + ":\n" + valueStyle.Render(wrappedDescription),
	}, "\n")

	finalHeaderOutput := infoContainerStyle.Render(headerFields)

	columns, rows := getControlRulesColumnsAndRows(control)

	effectiveRuleLimit := rowLimit
	if effectiveRuleLimit > len(control.Rules) {
		effectiveRuleLimit = len(control.Rules)
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(effectiveRuleLimit+1), // +1 for header row
	)

	tbl.SetStyles(table.Styles{
		Header: tableHeaderStyle,
		Cell:   tableCellStyle,
	})

	return terminal.Model{
		Table:     tbl,
		HeaderMsg: finalHeaderOutput,
		HelpMsg:   fmt.Sprintf("Showing %d of %d available rules. Use --limit to change max table rows.", effectiveRuleLimit, len(control.Rules)),
	}
}

// newRuleInfoModel creates a Bubble Tea model for displaying specific rule details.
func newRuleInfoModel(ruleDetails RuleInfo, setParameters map[string][]string, rowLimit int) terminal.Model {

	headerFields := strings.Join([]string{
		renderKeyValuePair("Rule ID", ruleDetails.ID),
		renderKeyValuePair("Rule Description", ruleDetails.Description),
		renderKeyValuePair("Check ID", ruleDetails.CheckID),
		renderKeyValuePair("Check Description", ruleDetails.CheckDescription),
	}, "\n")

	finalHeaderOutput := infoContainerStyle.Render(headerFields)

	// Prepare the parameters table
	columns, rows := getRuleParametersColumnsAndRows(ruleDetails, setParameters)

	effectiveParamLimit := rowLimit
	if effectiveParamLimit > len(rows) {
		effectiveParamLimit = len(rows)
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(effectiveParamLimit+1), // +1 for header row
	)

	tbl.SetStyles(table.Styles{
		Header: tableHeaderStyle,
		Cell:   tableCellStyle,
	})

	helpMsg := fmt.Sprintf("Showing %d of %d available parameters. Use --limit to change max table rows.", effectiveParamLimit, len(rows))
	if len(rows) == 0 {
		helpMsg = "No parameters found for this rule."
	}

	return terminal.Model{
		Table:     tbl,
		HeaderMsg: finalHeaderOutput,
		HelpMsg:   helpMsg,
	}
}

// newControlListModel creates a Bubble Tea model for displaying the list of controls.
func newControlListModel(controls []Control, rowLimit int) terminal.Model {
	columns, rows := getControlListColumnsAndRows(controls)

	// Ensure row limit does not exceed the total length of the controls list
	effectiveRowLimit := rowLimit
	if effectiveRowLimit > len(controls) {
		effectiveRowLimit = len(controls)
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(effectiveRowLimit+1), // +1 for header row
	)

	tbl.SetStyles(table.Styles{
		Header: tableHeaderStyle,
		Cell:   tableCellStyle,
	})

	return terminal.Model{
		Table:   tbl,
		HelpMsg: fmt.Sprintf("Showing %d of %d available controls. Use --limit to change max table rows.", effectiveRowLimit, len(controls)),
	}
}

// displayControlInfo handles displaying information for a specific control.
func displayControlInfo(opts *infoOptions, controlMap map[string]Control) error {
	control, ok := controlMap[opts.controlID]
	if !ok {
		return fmt.Errorf("control '%s' does not exist in workspace", opts.controlID)
	}

	if opts.plain {
		cols, rows := getControlRulesColumnsAndRows(control)

		_, _ = fmt.Fprintf(opts.Out, "Control ID: %s \n", control.ID)
		_, _ = fmt.Fprintf(opts.Out, "Title: %s \n", control.Title)
		_, _ = fmt.Fprintf(opts.Out, "Implementation Status: %s \n", control.ImplemenationStatus)
		_, _ = fmt.Fprintf(opts.Out, "Description: %s \n", control.Description)
		_, _ = fmt.Fprintln(opts.Out)
		terminal.ShowPlainTable(opts.Out, cols, rows)
		return nil
	} else {
		model := newControlInfoModel(control, opts.limit)
		return runBubbleTeaProgram(model, opts.Out)
	}

}

// displayRuleInfo handles displaying information for a specific rule.
func displayRuleInfo(opts *infoOptions, ruleID string, ruleRemarksMap map[string]string, remarksPropsMap map[string][]oscalTypes.Property, setParameters map[string][]string) error {
	remarksForRule, ok := ruleRemarksMap[ruleID]
	if !ok || remarksForRule == "" {
		return fmt.Errorf("rule '%s' remarks not found in workspace", ruleID)
	}

	propsForRule, ok := remarksPropsMap[remarksForRule]
	if !ok || len(propsForRule) == 0 {
		return fmt.Errorf("properties for rule '%s' (remarks '%s') not found", ruleID, remarksForRule)
	}

	ruleDetails := extractRuleDetails(propsForRule)
	ruleDetails.ID = ruleID // Ensure ID is set for consistency

	if opts.plain {
		_, _ = fmt.Fprintf(opts.Out, "Rule ID: %s \n", ruleDetails.ID)
		_, _ = fmt.Fprintf(opts.Out, "Rule Description: %s \n", ruleDetails.Description)
		_, _ = fmt.Fprintf(opts.Out, "Check ID: %s \n", ruleDetails.CheckID)
		_, _ = fmt.Fprintf(opts.Out, "Check Description: %s \n", ruleDetails.CheckDescription)
		_, _ = fmt.Fprintln(opts.Out)
		cols, rows := getRuleParametersColumnsAndRows(ruleDetails, setParameters)
		terminal.ShowPlainTable(opts.Out, cols, rows)
		return nil
	} else {
		model := newRuleInfoModel(ruleDetails, setParameters, opts.limit)
		return runBubbleTeaProgram(model, opts.Out)
	}

}

// displayAllControls handles displaying a list controls in the framework.
func displayAllControls(opts *infoOptions, controlMap map[string]Control) error {
	var controls []Control
	for _, control := range controlMap {
		controls = append(controls, control)
	}

	if opts.plain {
		cols, rows := getControlListColumnsAndRows(controls)
		terminal.ShowPlainTable(opts.Out, cols, rows)
		return nil
	} else {
		model := newControlListModel(controls, opts.limit)
		return runBubbleTeaProgram(model, opts.Out)
	}
}
