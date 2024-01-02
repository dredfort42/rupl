package configreader

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"web/internal/colors"
)

// ConfigMap is a map containing configuration properties.
type ConfigMap map[string]string

// ReadConfig reads a configuration file to a ConfigMap and returns an error if it fails.
func ReadConfig(path string, configMap *map[string]string) error {
	file, err := os.Open(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, colors.RED+"Failed to open file: %s\n"+colors.RESET, err)
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
				(*configMap)[parameter] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf(colors.GREEN+"Successfully read configuration from file: %s\n"+colors.RESET, path)

	return nil
}

// Get configuration from main and extra .config files and returns a ConfigMap and an error if it fails.
func GetConfig() (ConfigMap, error) {
	success := false
	configMap := make(map[string]string)

	// Read main config file
	if err := ReadConfig("/app/config/service.cfg", &configMap); err == nil {
		success = true
	}

	// Read extra config file
	if err := ReadConfig("/app/local.cfg", &configMap); err == nil {
		success = true
	}

	if !success {
		fmt.Println(colors.RED + "Failed to read configuration" + colors.RESET)
		return nil, fmt.Errorf("Failed to read configuration")
	} else {
		fmt.Println(colors.GREEN + "Configuration read successfully!" + colors.RESET)
		return configMap, nil
	}
}

// PrintConfig prints a ConfigMap to stdout.
func PrintConfig(configMap ConfigMap) {
	for key, value := range configMap {
		fmt.Printf("%s: %s\n", key, value)
	}
}
