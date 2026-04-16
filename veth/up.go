package veth

import (
	"fmt"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func Up(ns, ifName string) error {
	_, originCloseFunc, err := namespace.GetOriginNs()
	if err != nil {
		return err
	}
	defer originCloseFunc()

	targetNsFd, targetNsCloseFunc, err := namespace.GetNsFd(ns)
	if err != nil {
		return fmt.Errorf("failed to namespace %s file descriptor: %v", ns, err)
	}
	defer targetNsCloseFunc()

	if err := netns.Set(targetNsFd); err != nil {
		return fmt.Errorf("failed to set to namespace %s: %v", ns, err)
	}

	link, err := netlink.LinkByName(ifName)
	if err != nil {
		return fmt.Errorf("failed to get link %s in namespace %s: %v", ifName, ns, err)
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return fmt.Errorf("failed to set link %s up in namespace %s: %v", ifName, ns, err)
	}

	return nil
}
