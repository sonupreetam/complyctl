package scan

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/complytime/complytime/cmd/openscap-plugin/config"
	"github.com/complytime/complytime/cmd/openscap-plugin/oscap"
)

func isXMLFile(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, fmt.Errorf("invalid XML: %w", err)
		}
	}
}

func validateDataStream(path string) (string, error) {
	datastream, err := config.ValidatePath(path, false)
	if err != nil {
		return "", err
	}

	if _, err := isXMLFile(datastream); err != nil {
		return "", err
	}
	return datastream, nil
}

func validateTailoringFile(path string) (string, error) {
	tailoringFile, err := config.ValidatePath(path, false)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		} else {
			return "", err
		}
	}

	if _, err := isXMLFile(tailoringFile); err != nil {
		return "", err
	}
	return tailoringFile, nil
}

func ScanSystem(cfg *config.Config, profile string) ([]byte, error) {
	openscapFiles, err := config.DefineFilesPaths(cfg)
	if err != nil {
		return nil, err
	}

	_, err = validateDataStream(openscapFiles["datastream"])
	if err != nil {
		return nil, err
	}

	policy, err := validateTailoringFile(openscapFiles["policy"])
	if err != nil {
		return nil, err
	}
	if policy == "" {
		openscapFiles["policy"] = ""
	}

	output, err := oscap.OscapScan(openscapFiles, profile)
	if err != nil {
		if output == nil {
			return nil, err
		} else {
			return output, err
		}
	}

	return output, nil
}
