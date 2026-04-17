package cmd

import (
	"github.com/Alonza0314/nsctl/internal/nsExec"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:                   "exec <namespace> -- <command>",
	Short:                 "Execute a command in a network namespace",
	Long:                  "Execute a command in a network namespace",
	DisableFlagsInUseLine: true,
	Run:                   execFunc,
}

func init() {
	nsctlCmd.AddCommand(execCmd)
}

func execFunc(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		errFormat(args)
	}

	if err := nsExec.Run(args[0], args[1:]); err != nil {
		errPrint(err)
	}
}
