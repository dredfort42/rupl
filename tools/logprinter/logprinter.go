package logprinter

import (
	"fmt"
)

// Color definitions
const (
	RED    = "\033[1;31m"
	GREEN  = "\033[1;32m"
	YELLOW = "\033[1;33m"
	RESET  = "\033[0m"
)

// PrintSuccess prints a success message
func PrintSuccess(msg string, info string) {
	if info == "" {
		fmt.Printf("%s[I] %s%s\n", GREEN, msg, RESET)
	} else {
		fmt.Printf("%s[I] %s: %s%s\n", GREEN, msg, info, RESET)
	}
}

// PrintWarning prints a warning message
func PrintWarning(msg string, info string) {
	if info == "" {
		fmt.Printf("%sW] %s%s\n", YELLOW, msg, RESET)
	} else {
		fmt.Printf("%s[W] %s: %s%s\n", YELLOW, msg, info, RESET)
	}
}

// PrintError prints an error message
func PrintError(msg string, err error) {
	if err == nil {
		fmt.Printf("%sE] %s%s\n", RED, msg, RESET)
	} else {
		fmt.Printf("%s[E] %s: %s%s\n", RED, msg, err.Error(), RESET)
	}
}
