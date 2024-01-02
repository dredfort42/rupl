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
func PrintSuccess(msg string) {
	fmt.Printf("[I] %s%s%s\n", GREEN, msg, RESET)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Printf("[W] %s%s%s\n", YELLOW, msg, RESET)
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Printf("[E] %s%s%s\n", RED, msg, RESET)
}
