// SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-2"
	"github.com/oscal-compass/compliance-to-policy-go/v2/framework"
	"github.com/oscal-compass/oscal-sdk-go/extensions"
	"github.com/oscal-compass/oscal-sdk-go/generators"
	"github.com/oscal-compass/oscal-sdk-go/settings"
	"github.com/spf13/cobra"

	"github.com/complytime/complytime/cmd/complytime/option"
	"github.com/complytime/complytime/internal/complytime"
)

const assessmentResultsLocationJson = "assessment-results.json"
const assessmentResultsLocationMd = "assessment-results.md"

// scanOptions defined options for the scan subcommand.
type scanOptions struct {
	*option.Common
	complyTimeOpts *option.ComplyTime
}

// scanCmd creates a new cobra.Command for the version subcommand.
func scanCmd(common *option.Common) *cobra.Command {
	scanOpts := &scanOptions{
		Common:         common,
		complyTimeOpts: &option.ComplyTime{},
	}
	cmd := &cobra.Command{
		Use:          "scan [flags]",
		Short:        "Scan environment with assessment plan",
		Example:      "complytime scan",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runScan(cmd, scanOpts)
		},
	}
	cmd.Flags().BoolP("with-md", "m", false, "If true, assessement-result markdown will be generated")
	scanOpts.complyTimeOpts.BindFlags(cmd.Flags())
	return cmd
}

func runScan(cmd *cobra.Command, opts *scanOptions) error {

	// Load settings from assessment plan
	ap, apCleanedPath, err := loadPlan(opts.complyTimeOpts)
	if err != nil {
		return err
	}

	planSettings, err := getPlanSettings(opts.complyTimeOpts, ap)
	if err != nil {
		return err
	}

	// Set the framework ID from state (assessment plan)
	frameworkProp, valid := extensions.GetTrestleProp(extensions.FrameworkProp, *ap.Metadata.Props)
	if !valid {
		return fmt.Errorf("error reading framework property from assessment plan")
	}
	opts.complyTimeOpts.FrameworkID = frameworkProp.Value
	logger.Debug(fmt.Sprintf("Framework property was successfully read from the assessment plan: %v.", frameworkProp))

	// Create the application directory if it does not exist
	appDir, err := complytime.NewApplicationDirectory(true)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Using application directory: %s", appDir.AppDir()))

	cfg, err := complytime.Config(appDir)
	if err != nil {
		return err
	}

	// set config logger to CLI charm logger
	cfg.Logger = logger

	manager, err := framework.NewPluginManager(cfg)
	if err != nil {
		return fmt.Errorf("error initializing plugin manager: %w", err)
	}

	pluginOptions := opts.complyTimeOpts.ToPluginOptions()
	plugins, cleanup, err := complytime.Plugins(manager, pluginOptions)
	if err != nil {
		return fmt.Errorf("errors launching plugins: %w", err)
	}
	defer cleanup()
	logger.Info(fmt.Sprintf("Successfully loaded %v plugin(s).", len(plugins)))

	allResults, err := manager.AggregateResults(cmd.Context(), plugins, planSettings)
	if err != nil {
		return err
	}

	// Collect results in a single report
	r, err := framework.NewReporter(cfg)
	if err != nil {
		return err
	}

	var allImplementations []oscalTypes.ControlImplementationSet
	var profileHref string
	for _, compDef := range cfg.ComponentDefinitions {
		for _, component := range *compDef.Components {
			if component.ControlImplementations == nil {
				continue
			}
			for _, implementation := range *component.ControlImplementations {
				frameworkShortName, found := settings.GetFrameworkShortName(implementation)
				fmt.Println(frameworkShortName)
				fmt.Println(frameworkProp.Value)
				// If the framework property value match the assessment plan framework property values
				// this is the correct control source.
				if found && frameworkShortName == frameworkProp.Value {
					profileHref = implementation.Source
					fmt.Println(profileHref)
					// The implementations would have been filtered later in settings.Framework, but no need to add extra
					// implementations that are not needed to the slice.
					allImplementations = append(allImplementations, *component.ControlImplementations...)
				}
			}
		}
	}

	implementationSettings, err := settings.Framework(opts.complyTimeOpts.FrameworkID, allImplementations)
	if err != nil {
		return err
	}

	planHref := fmt.Sprintf("file://%s", apCleanedPath)
	assessmentResults, err := r.GenerateAssessmentResults(cmd.Context(), planHref, implementationSettings, allResults)
	if err != nil {
		return err
	}
	filePath := filepath.Join(opts.complyTimeOpts.UserWorkspace, assessmentResultsLocationJson)
	cleanedPath := filepath.Clean(filePath)
	err = complytime.WriteAssessmentResults(&assessmentResults, cleanedPath)
	if err != nil {
		return err
	}
	outputFlag, _ := cmd.Flags().GetBool("with-md")
	if outputFlag {
		profile, err := complytime.LoadProfile(appDir, profileHref)
		if err != nil {
			return err
		}

		if len(profile.Imports) != 1 {
			return errors.New("profile imports must be one")
		}
		catalog, err := complytime.LoadCatalogSource(appDir, profile.Imports[0].Href)
		if err != nil {
			return err
		}
		filePath := filepath.Join(opts.complyTimeOpts.UserWorkspace, assessmentResultsLocationMd)
		cleanedPath := filepath.Clean(filePath)
		templateValues, err := framework.CreateTemplateValues(*catalog, *assessmentPlan, assessmentResults)
		if err != nil {
			return err
		}
		assessmentResultsMd, err := templateValues.GenerateAssessmentResultsMd(cleanedPath)
		if err != nil {
			return err
		}
		err = os.WriteFile(cleanedPath, assessmentResultsMd, 0600)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No assessment result markdown will be generated.")
	}
	logger.Info(fmt.Sprintf("The assessment results were successfully written to %v.", assessmentResultsLocation))
	return nil
}
