package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

func Delete(name string) error {
	if found, err := GetNs(name); err != nil {
		return err
	} else {
		if !found {
			return fmt.Errorf("namespace %s not found", name)
		}
	}

	if err := netns.DeleteNamed(GetNsName(name)); err != nil {
		return fmt.Errorf("failed to delete ns: %v", err)
	}

	return nil
}
