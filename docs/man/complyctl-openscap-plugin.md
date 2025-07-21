% COMPLYCTL_OPENSCAP-PLUGIN(7) Complyctl OpenSCAP Plugin
% Qingmin Duanmu<qduanmu@redhat.com>
% July 2025

# NAME

complyctl-openscap-plugin - a plugin which extends the complyctl capabilities to use OpenSCAP.

# DESCRIPTION

The plugin is not meant to be executed directly, it communicates with complyctl via gRPC. It has configurable options that can be configured via a manifest file, complyctl processes the manifest file and sends the configuration values to the plugin.

When the plugin receives the **generate** command from complyctl, it will generate a tailing policy file and remediation files for bash, ansible, and imagebuilder. The generated files are placed in the **openscap** directory under user workspace.

When the plugin receives the **scan** command from complyctl, it will scan the system with **oscap** and return the observations to complyctl.

The generated remediation files from complyctl are based on the whole policy, it's not targeted to remediate specific findings. **oscap** could be used to manually generate remediation scripts containing only remediation for failed rules based on complyctl scan result.

# FILES

**/usr/share/complytime/plugins/c2p-openscap-manifest.json**
Default plugin manifest file.

**/etc/complytime/config.d/c2p-openscap-manifest.json**
Optional drop-in manifest file with customized plugin configurations.

# EXAMPLES

The following steps create a remediation Ansible Playbook that contains only the remediations required to align your system with a specific baseline, based on the complyctl scan result.

Step 1: generate a new assessment plan for the specified compliance framework, e.g., anssi_bp28_minimal

$ complyctl plan anssi_bp28_minimal

Suppose the user workspace is ~/complytime, the generated assessment plan file would be ~/complytime/assessment-plan.json.

Step2: generate PVP policy from the generated assessment plan

$ complyctl generate

With this step, the generated policy file could be found here: ~/complytime/openscap/policy/tailoring_policy.xml.

Step 3: scan environment with assessment plan to get the scan result

$ complyctl scan

After scanning, the generated openscap scan result is in ~/complytime/openscap/results/results.xml.

Step 4: find value of the result ID in scan result file

$ oscap info ~/complytime/openscap/results/results.xml

The result ID in the example would be xccdf_org.open-scap_testresult_xccdf_complytime.openscapplugin_profile_anssi_bp28_minimal_complytime.

Step 5: generate a remediation Ansible Playbook based on the scan result

$ oscap xccdf generate fix --fix-type ansible --result-id xccdf_org.open-scap_testresult_xccdf_complytime.openscapplugin_profile_anssi_bp28_minimal_complytime --tailoring-file ~/complytime/openscap/policy/tailoring_policy.xml --output output_remediations.yml ~/complytime/openscap/results/results.xml

# SEE ALSO

complyctl(1), c2p-openscap-manifest.json(5), oscap(8)

See the upstream projects at https://github.com/complytime/complyctl and https://github.com/OpenSCAP/openscap for more detailed documentation.

# COPYRIGHT

© 2025 Red Hat, Inc. complyctl-openscap-plugin is released under the terms of the Apache-2.0 license.
