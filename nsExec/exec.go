package nsExec

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Alonza0314/nsctl/namespace"
)

const netnsPath = "/var/run/netns"

func Run(ns string, command []string) error {
	found, err := namespace.GetNs(ns)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("namespace %s not found", ns)
	}

	nsPath := netnsPath + "/" + namespace.GetNsName(ns)

	args := []string{"--net=" + nsPath, "--"}
	args = append(args, command...)

	cmd := exec.Command("nsenter", args...)
	cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin

	var execErr error
	if execErr = cmd.Run(); execErr != nil {
		if exitErr, ok := execErr.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 130 {
				return nil
			}

		}
	}
	return execErr
}
