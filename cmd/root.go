package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var nsctlCmd = &cobra.Command{
	Use:   "nsctl",
	Short: "A quick namespace management tool.",
	Long:  "Build your own network topology with pure Linux namespaces — no containers, no overhead, full control.",
}

func Execute() {
	if err := nsctlCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
