package cmd

import (
	"fmt"

	"github.com/Alonza0314/nsctl/internal/veth"
	"github.com/spf13/cobra"
)

var netCmd = &cobra.Command{
	Use:                   "net <connect|disconnect|list|set-ip|up|down> [args]",
	Short:                 "Manage the networks between namespace",
	Long:                  "Manage the networks between namespace",
	DisableFlagsInUseLine: true,
	Run:                   netFunc,
}

func init() {
	nsctlCmd.AddCommand(netCmd)
}

func netFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		errEmptyAction()
	}

	switch args[0] {
	case "connect":
		if len(args) != 3 {
			errFormat(args)
		} else {
			if err := veth.Connect(args[1], args[2]); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Network connect between %s and %s successfully, named: %s and %s\n", args[1], args[2], veth.GetVethName(args[1], args[2]), veth.GetVethName(args[2], args[1]))
			}
		}
	case "disconnect":
		if len(args) != 3 {
			errFormat(args)
		} else {
			if err := veth.Disconnect(args[1], args[2]); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Network disconnect between %s and %s successfully\n", args[1], args[2])
			}
		}
	case "list":
		if len(args) != 2 {
			errFormat(args)
		} else {
			if list, err := veth.List(args[1]); err != nil {
				errPrint(err)
			} else {
				fmt.Print(list)
			}
		}
	case "set-ip":
		if len(args) != 4 {
			errFormat(args)
		} else {
			if err := veth.SetIp(args[1], args[2], args[3]); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Interface %s in namespace %s set IP with %s successfully\n", args[2], args[1], args[3])
			}
		}
	case "up":
		if len(args) != 3 {
			errFormat(args)
		} else {
			if err := veth.UpDown(args[1], args[2], true); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Interface %s in namespace %s set up successfully\n", args[2], args[1])
			}
		}
	case "down":
		if len(args) != 3 {
			errFormat(args)
		} else {
			if err := veth.UpDown(args[1], args[2], false); err != nil {
				errPrint(err)
			} else {
				fmt.Printf("Interface %s in namespace %s set down successfully\n", args[2], args[1])
			}
		}
	default:
		errInvalidAction(args[0])
	}
}
