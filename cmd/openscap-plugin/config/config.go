// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

const PluginDir = "openscap"

type Config struct {
	Files struct {
		Workspace  string `config:"workspace"`
		Datastream string `config:"datastream"`
		Results    string `config:"results"`
		ARF        string `config:"arf"`
		Policy     string `config:"policy"`
	}
	Parameters struct {
		Profile string `config:"profile"`
	}
}

// NewConfig creates a new, empty Config.
func NewConfig() *Config {
	return &Config{}
}

// LoadSettings sets the values in the Config from a given config map and
// performs validation.
func (c *Config) LoadSettings(config map[string]string) error {
	filesVal := reflect.ValueOf(&c.Files).Elem()
	if err := setConfigStruct(filesVal, config); err != nil {
		return err
	}
	paramVal := reflect.ValueOf(&c.Parameters).Elem()
	if err := setConfigStruct(paramVal, config); err != nil {
		return err
	}
	return c.validate()
}

func (c *Config) validate() error {
	// String values to sanitize
	inputValues := []*string{
		&c.Files.Policy,
		&c.Files.Results,
		&c.Files.ARF,
		&c.Parameters.Profile,
	}

	for _, inputValue := range inputValues {
		sanitized, err := SanitizeInput(*inputValue)
		if err != nil {
			return err
		}
		*inputValue = sanitized
	}

	_, err := SanitizeAndValidatePath(c.Files.Datastream, false)
	if err != nil {
		return fmt.Errorf("invalid datastream path: %s: %w", c.Files.Datastream, err)
	}

	if err := defineFilesPaths(c); err != nil {
		return err
	}
	return nil
}

func SanitizeInput(input string) (string, error) {
	safePattern := regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`)
	if !safePattern.MatchString(input) {
		return "", fmt.Errorf("input contains unexpected characters: %s", input)
	}
	return input, nil
}

func expandPath(path string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err != nil {
			return "", fmt.Errorf("failed to identify current user: %w", err)
		}
		homeDir := usr.HomeDir
		// Replace "~" with the home directory
		return filepath.Join(homeDir, path[1:]), nil
	}
	return path, nil
}

func validatePath(path string, shouldBeDir bool) (string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to confirm path existence: %w", err)
	}

	if shouldBeDir && !stat.IsDir() {
		return "", fmt.Errorf("expected a directory, but found a file at path: %s", path)
	}
	if !shouldBeDir && stat.IsDir() {
		return "", fmt.Errorf("expected a file, but found a directory at path: %s", path)
	}

	return path, nil
}

func SanitizePath(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	expandedPath, err := expandPath(cleanPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand path: %w", err)
	}
	return expandedPath, nil
}

func SanitizeAndValidatePath(path string, shouldBeDir bool) (string, error) {
	cleanPath, err := SanitizePath(path)
	if err != nil {
		return "", err
	}
	validPath, err := validatePath(cleanPath, shouldBeDir)
	if err != nil {
		return "", err
	}
	return validPath, nil
}

func ensureDirectory(path string) error {
	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(path, 0750)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		log.Printf("Directory created: %s\n", path)
	} else if err != nil {
		return fmt.Errorf("error checking directory: %w", err)
	}

	return nil
}

func ensureWorkspace(cfg *Config) (map[string]string, error) {
	workspacePath, err := SanitizePath(cfg.Files.Workspace)
	if err != nil {
		return nil, fmt.Errorf("failed to sanitize workspace path %s: %w", cfg.Files.Workspace, err)
	}

	workspace, err := validatePath(workspacePath, true)
	if err != nil {
		log.Printf("Informed workspace was not found. It will be created.")
		workspace = workspacePath
	}

	directories := map[string]string{
		"workspace":  workspace,
		"pluginDir":  filepath.Join(workspace, PluginDir),
		"policyDir":  filepath.Join(workspace, PluginDir, "policy"),
		"resultsDir": filepath.Join(workspace, PluginDir, "results"),
	}

	for key, dir := range directories {
		if err := ensureDirectory(dir); err != nil {
			return nil, fmt.Errorf("failed to ensure directory %s (%s): %w", dir, key, err)
		}
	}

	return directories, nil
}

func defineFilesPaths(cfg *Config) error {
	directories, err := ensureWorkspace(cfg)
	if err != nil {
		return err
	}

	cfg.Files.Policy = filepath.Join(directories["policyDir"], cfg.Files.Policy)
	cfg.Files.Results = filepath.Join(directories["resultsDir"], cfg.Files.Results)
	cfg.Files.ARF = filepath.Join(directories["resultsDir"], cfg.Files.ARF)

	return nil
}

// setConfigStruct populates struct fields with matching tags to values
// in a given config map.
func setConfigStruct(val reflect.Value, config map[string]string) error {
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldType := t.Field(i)
		key := fieldType.Tag.Get("config")
		value, ok := config[key]
		if !ok {
			return fmt.Errorf("missing configuration value for option %q (field: %s)", key, fieldType.Name)
		}

		fieldVal := val.Field(i)
		fieldVal.SetString(value)
	}
	return nil
}
