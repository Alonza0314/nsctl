package veth

import (
	"fmt"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/vishvananda/netlink"
)

func Connect(ns1, ns2 string) error {
	ns1Fd, ns1CloseFunc, err := namespace.GetNsFd(ns1)
	if err != nil {
		return fmt.Errorf("failed to namespace %s file descriptor: %v", ns1, err)
	}
	defer ns1CloseFunc()
	ns2Fd, ns2CloseFunc, err := namespace.GetNsFd(ns2)
	if err != nil {
		return fmt.Errorf("failed to namespace %s file descriptor: %v", ns2, err)
	}
	defer ns2CloseFunc()

	vethName1, vethName2 := GetVethName(ns1, ns2), GetVethName(ns2, ns1)

	vethLink := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: vethName1,
		},
		PeerName: vethName2,
	}

	if err := netlink.LinkAdd(vethLink); err != nil {
		return fmt.Errorf("failed to create %s and %s: %v\n", vethName1, vethName2, err)
	}

	link1, err := netlink.LinkByName(vethName1)
	if err != nil {
		return fmt.Errorf("failed to get link %s: %v\n", vethName1, err)
	}
	link2, err := netlink.LinkByName(vethName2)
	if err != nil {
		return fmt.Errorf("failed get link %s: %v\n", vethName2, err)
	}

	if err := netlink.LinkSetNsFd(link1, int(ns1Fd)); err != nil {
		return fmt.Errorf("failed to move %s to %s: %v", vethName1, ns1, err)
	}
	if err := netlink.LinkSetNsFd(link2, int(ns2Fd)); err != nil {
		return fmt.Errorf("failed to move %s to %s: %v", vethName2, ns2, err)
	}

	return nil
}
