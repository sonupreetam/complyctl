// SPDX-License-Identifier: Apache-2.0

package oscap

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func constructScanCommand(openscapFiles map[string]string, profile string) ([]string, error) {
	datastream := openscapFiles["datastream"]
	tailoringFile := openscapFiles["policy"]
	resultsFile := openscapFiles["results"]
	arfFile := openscapFiles["arf"]

	cmd := []string{
		"oscap",
		"xccdf",
		"eval",
		"--profile",
		profile,
		"--results",
		resultsFile,
		"--results-arf",
		arfFile,
	}

	if _, err := os.Stat(tailoringFile); err == nil {
		cmd = append(cmd, "--tailoring-file", tailoringFile)
	}

	cmd = append(cmd, datastream)
	return cmd, nil
}

func OscapScan(openscapFiles map[string]string, profile string) ([]byte, error) {
	command, err := constructScanCommand(openscapFiles, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to construct command %s: %w", command, err)
	}

	cmdPath, err := exec.LookPath(command[0])
	if err != nil {
		return nil, fmt.Errorf("command not found: %s", command[0])
	}

	log.Printf("Executing the command: '%v'", command)
	cmd := exec.Command(cmdPath, command[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if err.Error() == "exit status 1" {
			return output, fmt.Errorf("%s: oscap error during evaluation", err)
		} else if err.Error() == "exit status 2" {
			log.Printf("%s: at least one rule resulted in fail or unknown", err)
			return output, nil
		} else {
			log.Printf("%s", err)
			return output, nil
		}
	}

	return output, nil
}
