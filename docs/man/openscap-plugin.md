% OPENSCAP-PLUGIN(3) Complyctl OpenSCAP Plugin
% Qingmin Duanmu<qduanmu@redhat.com>
% July 2025

# NAME

openscap-plugin - a plugin which extends the complyctl capabilities to use OpenSCAP.

# DESCRIPTION

The plugin is not meant to be executed directly, it communicates with complyctl via gRPC. It has configurable options that can be configured via a manifest file, Complyctl processes the manifest file and sends the configuration values to the plugin.

When the plugin receives the **generate** command from complyctl, it will generate a tailing policy file and remediation files for bash, ansible, and imagebuilder. The generated files are placed in the **openscap** directory under workspace.

When the plugin receives the **scan** command from complyctl, it will scan the system with **oscap** and return the observations to complyctl.

The generated remediation files from complyctl are based on the whole policy, it's not targeted to remediate specific findings. **oscap** could be used to manually generate remediation scripts containing only remediation for failed rules based on complyctl scan result.

# FILES

**/usr/share/complytime/plugins/c2p-openscap-manifest.json**
Default plugin manifest file.

**/etc/complytime/config.d/c2p-openscap-manifest.json**
Drop in manifest file with customized plugin configurations.

# EXAMPLES

The following commands create a remediation Ansible Playbook that contains only the remediations required to align your system with a specific baseline, based on the complyctl scan result, you need to know the ID of the profile according to which you want to remediate your system.

Generate PVP policy from an existing assessment plan:

$ complyctl generate

Scan environment with assessment plan to get the scan result:

$ complyctl scan

Find value of the result ID in scan result file:

$ oscap info <scan_result.xml>

Generate a remediation Ansible Playbook based on the scan result:

$ oscap xccdf generate fix --fix-type ansible --result-id <result_ID> --tailoring-file <tailoring_policy.xml> --output <profile_remediations.yml> <scan_result.xml>

# SEE ALSO

complyctl(1), c2p-openscap-manifest.json(5), oscap(8)

See the upstream project at https://github.com/complytime/complyctl for more detailed documentation.

# COPYRIGHT

© 2025 Red Hat, Inc. openscap-plugin is released under the terms of the Apache-2.0 license.
