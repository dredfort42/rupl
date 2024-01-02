package configreader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	log "github.com/dredfort42/rupl/tools/logprinter"
)

// ConfigMap is a map containing configuration properties.
type ConfigMap map[string]string

// ReadConfig reads a configuration file to a ConfigMap and returns an error if it fails.
func ReadConfig(path string, config *ConfigMap) error {
	file, err := os.Open(path)

	if err != nil {
		log.PrintError("Failed to open file", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			before, after, found := strings.Cut(line, "=")
			if found {
				parameter := strings.TrimSpace(before)
				value := strings.TrimSpace(after)
				(*config)[parameter] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	log.PrintSuccess("Successfully read configuration from file", path)

	return nil
}

// Get configuration from global and local .cfg files and returns a ConfigMap and an error if it fails.
func GetConfig() (ConfigMap, error) {
	success := false
	config := make(ConfigMap)

	// Read global config file
	if err := ReadConfig("./global.cfg", &config); err == nil {
		success = true
	}

	// Read local config file
	if err := ReadConfig("./local.cfg", &config); err == nil {
		success = true
	}

	if !success {
		log.PrintError("Failed to read configuration", nil)
		return nil, fmt.Errorf("Failed to read configuration")
	} else {
		log.PrintSuccess("Successfully read configuration", "")
		return config, nil
	}
}

// PrintConfig prints a ConfigMap to stdout.
func PrintConfig(config ConfigMap) {
	log.PrintSuccess("Configuration", "")
	for key, value := range config {
		log.PrintSuccess(key, value)
	}
}
