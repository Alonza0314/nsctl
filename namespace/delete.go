package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

func Delete(name string) error {
	nsList, err := getNsList()
	if err != nil {
		return err
	}

	found := false
	for _, ns := range nsList {
		if ns == name {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("namespace %s not found", name)
	}

	if err := netns.DeleteNamed(GetNsName(name)); err != nil {
		return fmt.Errorf("failed to delete ns: %v", err)
	}

	return nil
}
