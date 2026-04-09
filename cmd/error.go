package cmd

import (
	"fmt"
	"os"
	"strings"
)

func errEmptyAction() {
	fmt.Println("Please provide an action, use '--help' to check command")
	os.Exit(1)
}

func errInvalidAction(action string) {
	fmt.Printf("Invalid action: %s, use '--help' to check command\n", action)
	os.Exit(1)
}

func errFormat(args []string) {
	fmt.Printf("Invalid format: '%s', use '--help' to check command\n", strings.Join(args, " "))
	os.Exit(1)
}

func errPrint(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}
