package cmd

import (
	"fmt"

	"github.com/Alonza0314/nsctl/internal/topo"
	"github.com/free-ran-ue/util"
	"github.com/spf13/cobra"
)

var topoCmd = &cobra.Command{
	Use:   "topo <apply|delete> [args]",
	Short: "Manage the topology of namespaces and networks",
	Long:  "Manage the topology of namespaces and networks",
	Run:   topoFunc,
}

func init() {
	nsctlCmd.AddCommand(topoCmd)
}

func topoFunc(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		errEmptyAction()
	}

	if !util.FileExists(args[1]) {
		errFileNotExist(args[1])
	}

	var topoStruct topo.Topology
	if err := util.LoadFromYaml(args[1], &topoStruct); err != nil {
		errYamlFormat(args[1], err)
	}

	switch args[0] {
	case "apply":
		if err := topo.Apply(&topoStruct); err != nil {
			errPrint(err)
		} else {
			fmt.Printf("Topology applied successfully from file: %s\n", args[1])
		}
	case "delete":
	default:
		errFormat(args)
	}
}
