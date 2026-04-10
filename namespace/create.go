package namespace

import (
	"fmt"

	"github.com/vishvananda/netns"
)

func Create(name string) error {
	nsList, err := getNsList()
	if err != nil {
		return err
	}

	for _, ns := range nsList {
		if ns == name {
			return fmt.Errorf("namespace %s already exists", name)
		}
	}

	originNs, err := netns.Get()
	if err != nil {
		return fmt.Errorf("failed to get origin ns: %v", err)
	}
	defer func() {
		if err := originNs.Close(); err != nil {
			fmt.Printf("failed to close origin ns: %v\n", err)
		}
	}()

	newNs, err := netns.NewNamed(NS_PREFIX + name)
	if err != nil {
		return fmt.Errorf("failed to create new ns: %v", err)
	}
	defer func() {
		if err := newNs.Close(); err != nil {
			fmt.Printf("failed to close new ns: %v\n", err)
		}
	}()

	if err := netns.Set(originNs); err != nil {
		return fmt.Errorf("failed to re-set to origin ns: %v", err)
	}

	return nil
}
