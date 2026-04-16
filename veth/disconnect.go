package veth

import (
	"fmt"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func Disconnect(ns1, ns2 string) error {
	originNs, err := netns.Get()
	if err != nil {
		return fmt.Errorf("failed to get origin ns: %v", err)
	}
	defer func() {
		if err := netns.Set(originNs); err != nil {
			fmt.Printf("failed to re-set to origin ns: %v\n", err)
		}
		if err := originNs.Close(); err != nil {
			fmt.Printf("failed to close origin ns: %v\n", err)
		}
	}()

	ns1Fd, ns1CloseFunc, err := namespace.GetNsFd(ns1)
	if err != nil {
		return fmt.Errorf("failed to namespace %s file descriptor: %v", ns1, err)
	}
	defer ns1CloseFunc()
	_, ns2CloseFunc, err := namespace.GetNsFd(ns2)
	if err != nil {
		return fmt.Errorf("failed to namespace %s file descriptor: %v", ns2, err)
	}
	defer ns2CloseFunc()

	if err := netns.Set(ns1Fd); err != nil {
		return fmt.Errorf("failed to set to namespace %s: %v", ns1, err)
	}

	vethName := GetVethName(ns1, ns2)
	link, err := netlink.LinkByName(vethName)
	if err != nil {
		return fmt.Errorf("failed to get link %s: %v", vethName, err)
	}

	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("failed to delete link %s: %v", vethName, err)
	}

	return nil
}
