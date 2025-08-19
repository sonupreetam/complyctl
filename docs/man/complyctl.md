% COMPLYCTL(1) Complyctl Manual
% Marcus Burghardt <maburgha@redhat.com>
% April 2025

# NAME

complyctl - Complyctl CLI perform compliance assessment activities using plugins for different underlying technologies.

# SYNOPSIS

**complyctl** [command] [flags]

# DESCRIPTION

Complyctl CLI leverages OSCAL to perform compliance assessment activities, using plugins for each stage of the lifecycle.

Complyctl can be extended to support desired policy engines (PVPs) by the use of plugins.
The plugin acts as the integration between complyctl and the PVPs native interface.
Each plugin is responsible for converting the policy content described in OSCAL into the input format expected by the PVP.
In addition, the plugin converts the raw results provided by the PVP into the schema used by complyctl to generate OSCAL output.

Plugins communicate with complyctl via gRPC and can be authored using any preferred language. The plugin acts as the gRPC server while the complyctl CLI acts as the client. When a complyctl command is run, it invokes the appropriate method served by the plugin.

Complyctl is built on https://github.com/oscal-compass/compliance-to-policy-go which provides a flexible plugin framework for leveraging OSCAL with various PVPs.

# COMMANDS

**completion**
Generate the autocompletion script for the specified shell.

**generate**
Generate PVP policy from an assessment plan.

**help**
Display help about any command.

**list**
List information about supported frameworks and components.

**info**
Display information about a framework's controls and rules.

**plan**
Generate a new assessment plan for a given compliance framework ID.

**scan**
Scan environment with assessment plan.

**version**
Print the version.

# OPTIONS

**-d**, **--debug**
Output debug logs.

**-h**, **--help**
Show help for complyctl.

Run **complyctl [command] --help** for more information about a specific command.

# EXAMPLES

## Leveraging the `info` command

Use the `info` command functionality to explore the relationship between controls, rules, and parameters within a framework.

```bash
complyctl info anssi_bp28_minimal --limit 5 
# List the controls in the framework
complyctl info anssi_bp28_minimal --control 
# List details about a control and rules used by that control
complyctl info anssi_bp28_minimal --rule 
# List the parameters used be a specific rule
complyctl info anssi_bp28_minimal --parameter 
# See the current set value and valid alternatives for a parameter
```

## Assessment Scoping using the `plan` command

The `plan` command is used for scoping an OSCAL Assessment Plan. To populate defaults for your assessment plan, use the `./bin/complyctl plan <framework-id> --dry-run --out config.yml` command. The `config.yml` will be populated with the default values from the selected `framework-id`. The fields of the `config.yml` can be updated to scope controls, rules, and parameters.

```bash
complyctl plan anssi_bp28_minimal --dry-run --out config.yml 
# Populate defaults in config.yml

complyctl plan anssi_bp28_minimal --scope-config config.yml  
# Configure plan based on updates made in config.yml
```

## Configuring the Assessment Plan 

### Excluding Controls and Rules

Excluding components of the assessment plan works the same way for both controls and rules. To exclude a control from the assessment plan, simply delete the entire group of YAML keys associated with the `controlId`. The deleted fields will exclude that control from the generated `assessment-plan.json`. Excluding rules from an assessment plan requires indication of the rule name at the `controlId` level or the `includedControls` level. To exclude a rule specific to a controlId, use the `excludeRules` YAML key. To globally exclude a rule for all included controls use the `globalExcludeRules` YAML key. The default state of the `config.yml` indicates all rules to be included by default (using the global wildcard "*").

### Selecting Parameter Values

Selecting parameter values in the `selectParameters` field will allow you to update the value that is set for a parameter. If the update is not a valid parameter alternative, the `assessment-plan.json` will not be written and an error will be produced with valid alternatives printed to the screen.

Default values for the `config.yml` are populated based on the `framework-id` selected. The `config.yml` can be updated to scope controls, rules, and parameters.

```yaml
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30 # Delete the entire group associated with r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  excludeRules: # Use to exclude a rule specific to the controlId
  - "accounts_password_set_max_life_root"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "1" # Update parameter value to an available alternative
  - name: var_password_pam_unix_rounds
    value: "11"
```
## Controls

To configure the controls included in the assessment plan, use the `controlId` and `controlTitle` YAML keys. The `includeRules` YAML key is used to include rules for a specific controlId. The `selectParameters` YAML key is used to configure the parameters applied to the rules for a specific controlId.

### Excluding Control(s)

To exclude a control in your assessment plan, delete the entire group of YAML keys associated with the `controlId`. The deleted fields will exclude that control from the generated `assessment-plan.json`. The activities associated with that `controlId` will be marked "skipped" if there are no other controls in scope of the activity.

```yaml
# Example deletion of r31: User Password Strength
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
```

Once the `config.yml` is updated, the `assessment-plan.json` will be generated with the updated list of included controls.

## Rules

All rules associated with a controlId are included by default and are indicated by the "*" wildcard. To exclude a rule specific to a controlId, use the `excludeRules` YAML key. To globally exclude a rule across all controls, use the `globalExcludeRules` YAML key.

### Exclude Rules for a Control

To exclude a rule specific to a `controlId`, use the `excludeRules` YAML key. The `excludeRules` YAML key takes priority over `includeControls` within a control. When the rule is excluded, the `selectParameters` values that are associated with the rule will not be considered as part of the plan.

```yaml
# Example excludeRules for controlId: r31
# excludeRules for account_password_set_max_life_root
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules: # Initial global wildcard
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "1"
  - name: var_password_pam_unix_rounds
    value: "11"
  excludeRules: # Use to exclude a rule specific to the controlId
  - "accounts_password_set_max_life_root"
```

### Global Exclude Rules

To exclude a rule across all controls, use the `globalExcludeRules` YAML key. The example below excludes all rules for all controlIds in the `config.yml`. The `globalExcludeRules` YAML key takes priority over `includeRules` globally.

```yaml
# Example using globalExcludeRules for all rules "*"
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "1"
  - name: var_password_pam_unix_rounds
    value: "11"
globalExcludeRules:
- "*"
```
One passing the `config.yml` with the `--scope-config config.yml` flag, the assessment plan will be generated with the updated list of rules. 

## Parameters

The parameters of the assessment-plan are grouped by remarks value. To configure the `selectParameters` field, update the second-level YAML key `value` with a valid alternative.

If you update the value of a parameter to an invalid alternative, you will receive an error that populates the available alternatives. 

### Initial Set Parameters

The content below reflects the set-parameter values for the frameworkId in the `config.yml`. The `selectParameters` are included underneath each `controlId` based on remarks grouping with rules from the OSCAL Component Definition implemented requirements. The parameter value can be configured based on the available alternatives.

```yaml
# Example default selectParameters value for controlId: r30 and r31.
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "1" # Initially set parameter value
  - name: var_password_pam_unix_rounds
    value: "11"
```

### Invalid Alternative Parameter Value Update

An invalid update to the `selectParameters` field with "test-error" will not write the `assessment-plan.json` with the updated value. Below, the controlId `r31` has an invalid parameter update to "var_password_pam_ucredit" in `r31`. When passing the `--scope-config config.yml` flag, the `assessment-plan.json` will not be written and an error will be produced with valid alternatives printed to the screen.

```yaml
# Example incorrect update to selectParameters value for controlId: r31.
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "test-error" # Update to "test-error"
  - name: var_password_pam_unix_rounds
    value: "11"
```

### Valid Alternative Parameter Value Update

A valid update to the `selectParameters` field with an available alternative will write the `assessment-plan.json` with the updated value. Below, the controlId `r31` has a valid parameter update to "var_password_pam_ucredit" in `r31`. When passing the `--scope-config config.yml` flag, the `assessment-plan.json` will be written reflecting that update.

```yaml
# Example update to selectParameters value for controlId: r31.
frameworkId: anssi_bp28_minimal
includeControls:
- controlId: r30
  controlTitle: Removal Of Unused User Accounts
  includeRules:
  - "*"
  selectParameters:
  - name: N/A
    value: N/A
- controlId: r31
  controlTitle: User Password Strength
  includeRules:
  - "*"
  selectParameters:
  - name: var_password_pam_ucredit
    value: "0" # Update to available alternative
  - name: var_password_pam_unix_rounds
    value: "11"
```

## Assessment Plan Scope Inheritance

When excluding a `controlId` from the `config.yml`, the initial "*" `includeRules` values will be skipped and not assessed for the `controlId` in the assessment plan. The activities of the assessment plan will indicate "skipped" for a rule that is globally excluded. Therefore, all parameters associated with a globally excluded rule will not be used in the generated `assessment-plan.json`.

# SEE ALSO

complyctl-openscap-plugin(7)

See the Upstream project at https://github.com/complytime/complyctl for more detailed documentation.

See https://github.com/oscal-compass/compliance-to-policy-go project.

# COPYRIGHT

© 2025 Red Hat, Inc. Complyctl is released under the terms of the Apache-2.0 license.
