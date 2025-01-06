// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
	"github.com/complytime/complytime/cmd/openscap-plugin/scan"
)

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./openscap-plugin.yml", "Path to config file")
	flag.Parse()

	configFile, err := config.SanitizeAndValidatePath(configPath, false)
	if err != nil {
		return "", err
	}

	return configFile, nil
}

func initializeConfig() (*config.Config, error) {
	configFile, err := parseFlags()
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
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

	output, err := scan.ScanSystem(config, "cis")
	if err != nil {
		log.Printf("%v", err)
	}

	if output != nil {
		fmt.Printf("Scan command output:\n%s", output)
	}
}
