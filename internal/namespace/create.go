package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

func Create(name string) error {
	found, err := GetNs(name)
	if err != nil {
		return err
	}
	if found {
		return fmt.Errorf("namespace %s already exists", name)
	}

	_, originCloseFunc, err := GetOriginNs()
	if err != nil {
		return err
	}
	defer originCloseFunc()

	newNs, err := netns.NewNamed(GetNsName(name))
	if err != nil {
		return fmt.Errorf("failed to create new ns: %v", err)
	}
	defer func() {
		if err := newNs.Close(); err != nil {
			fmt.Printf("failed to close new ns: %v\n", err)
		}
	}()

	return nil
}
