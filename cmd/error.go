package cmd

import "fmt"

func errEmptyAction() {
	fmt.Println("Please provide an action, use '--help' to check command")
}

func errInvalidAction(action string) {
	fmt.Printf("Invalid action: %s, use '--help' to check command\n", action)
}

func errSystemError(err error) {
	fmt.Printf("System error: %v\n", err)
}