# openscap-plugin

## Overview

NOTE: The development of this plugin is in progress and therefore it should only be used for testing purposes at this point.

**openscap-plugin** is a plugin which extends the **ComplyTime** capabilities to use OpenSCAP. The plugin communicates with **ComplyTime** via gRPC, providing a standard and consistent communication mechanism that gives independence for plugin developers to choose their preferred languages. This plugin is structured to allow modular development, ease of packaging, and maintainability.

For now, this plugin is developed together with ComplyTime for better collaboration during this phase of the project. In the future, this plugin will likely be decoupled into its own repository.

## Plugin Structure

```
openscap-plugin/
├── config/               # Package for plugin configuration
│ ├── config_test.go      # Tests for functions in config.go
│ └── config.go           # Main code used to process plugin configuration
├── oscap/                # Package to interact with oscap command
│ ├── oscap_test.go       # Tests for functions in oscap.go
│ └── oscap.go            # Main code used to interact with oscap command
├── scan/                 # Package to process system scan instructions
│ ├── scan_test.go        # Tests for functions in scan.go
│ └── scan.go             # Main code used to process scan instructions
├── server/               # Package to process server functions. Here is where the plugin communicates with ComplyTime CLI
│ ├── server_test.go      # Tests for functions in server.go
│ └── server.go           # Main code used to process server functions
├── xccdf/                # Package to process SCAP Datastreams
│ ├── datastream_test.go  # Tests for functions in datastream.go
│ ├── datastream.go       # Main code used to process Datastream files
│ ├── tailoring_test.go   # Tests for functions in tailoring.go
│ └── tailoring.go        # Main code used to generate tailoring files based on OSCAL and available Datastreams.
├── openscap-config.yml   # Example of plugin configuration file (still in development)
└── README.md             # This file
```

## Features
### Scan
When the plugin receives the scan command from ComplyTime, it will use the configured Datastream and Policy (tailoring file) to:
* Validate the SCAP files
* Assembly the `oscap` command
* Scan the system saving results in an ARF file according to the plugin configuration
* Process the results and return observations to ComplyTime.

### Generate
This feature is currently under development.

## Installation

### Prerequisites

- **Go** version 1.20 or higher
- **Make** (optional, for using the `Makefile` if included)
- **scap-security-guide** package installed

### Clone the repository

```bash
git clone https://github.com/complytime/complytime.git
cd complytime
```

## Build Instructions
To compile complytime and openscap-plugin:

```bash
make build
```

## Running
Install the plugin
```bash
mkdir -p ~/.config/complytime/plugins
cp -rp bin/openscap-plugin ~/.config/complytime/plugins
cp -rp cmd/openscap-plugin/openscap-plugin.yml ~/.config/complytime/plugins
```

Create the manifest
```bash
checksum=$(sha256sum ~/.config/complytime/plugins/openscap-plugin| cut -d ' ' -f 1 )
cat > ~/.config/complytime/plugins/c2p-openscap-manifest.json << EOF
{
  "metadata": {
    "id": "openscap",
    "description": "My openscap plugin",
    "version": "0.0.1",
    "types": ["pvp"]
  },
  "executablePath": "openscap-plugin",
  "sha256": "$checksum"
}
EOF
```

### Scan
Run ComplyTime scan command to use plugin
```bash
bin/complytime scan
```

After the scan, check the files in "user_workspace" directory.

### Testing
Tests are organized within each package. Whenever possible a unit test is created for every function.

Run tests using:

```bash
make test-units
```
