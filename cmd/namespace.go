package cmd

import (
	"fmt"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/spf13/cobra"
)

var nsCmd = &cobra.Command{
	Use:                   "ns <create|delete|list> [namespace]",
	Short:                 "Manage network namespaces",
	Long:                  "Manage network namespaces",
	DisableFlagsInUseLine: true,
	Run:                   nsFunc,
}

func init() {
	nsctlCmd.AddCommand(nsCmd)
}

func nsFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		errEmptyAction()
	}

	switch args[0] {
	case "create":
		if len(args) != 2 {
			errFormat(args)
		} else {
			if err := namespace.Create(args[1]); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Namespace %s created successfully\n", args[1])
			}
		}
	case "delete":
		if len(args) != 2 {
			errFormat(args)
		} else {
			if err := namespace.Delete(args[1]); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Namespace %s deleted successfully\n", args[1])
			}
		}
	case "list":
		if err := namespace.List(); err != nil {
			errPrint(err)
		}
	default:
		errInvalidAction(args[0])
	}
}
