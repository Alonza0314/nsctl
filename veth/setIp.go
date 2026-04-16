package veth

import (
	"fmt"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func SetIp(ns, ifName, ip string) error {
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

	addr, err := netlink.ParseAddr(ip)
	if err != nil {
		return fmt.Errorf("failed to parse IP address %s: %v", ip, err)
	}

	if err := netlink.AddrAdd(link, addr); err != nil {
		return fmt.Errorf("failed to add IP address %s to interface %s in namespace %s: %v", ip, ifName, ns, err)
	}

	return nil
}
