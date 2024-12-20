// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
	"github.com/complytime/complytime/cmd/openscap-plugin/scan"
)

func getConfigFile() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "openscap-plugin.yml")

	// Construct the full path to the file
	configFile, err := config.SanitizeAndValidatePath(configPath, false)
	if err != nil {
		return "", fmt.Errorf("failed to sanitize or validate config file: %w", err)
	}

	return configFile, nil
}

func initializeConfig() (*config.Config, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error locating config file: %w", err)
	}

	config, err := config.ReadConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config from %s: %w", configFile, err)
	}

	return config, nil
}

func main() {
	config, err := initializeConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	output, err := scan.ScanSystem(config, config.Parameters.Profile)
	if err != nil {
		log.Printf("%v", err)
	}

	if output != nil {
		fmt.Printf("Scan command output:\n%s", output)
	}
}
