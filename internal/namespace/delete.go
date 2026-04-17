package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

func Delete(name string) error {
	found, err := GetNs(name)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("namespace %s not found", name)
	}

	if err := netns.DeleteNamed(GetNsName(name)); err != nil {
		return fmt.Errorf("failed to delete ns: %v", err)
	}

	return nil
}
