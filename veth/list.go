package veth

import (
	"fmt"
	"strings"

	"github.com/Alonza0314/nsctl/namespace"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func List(ns string) (string, error) {
	_, originCloseFunc, err := namespace.GetOriginNs()
	if err != nil {
		return "", err
	}
	defer originCloseFunc()

	targetNsFd, targetNsCloseFunc, err := namespace.GetNsFd(ns)
	if err != nil {
		return "", fmt.Errorf("failed to namespace %s file descriptor: %v", ns, err)
	}
	defer targetNsCloseFunc()

	if err := netns.Set(targetNsFd); err != nil {
		return "", fmt.Errorf("failed to set to namespace %s: %v", ns, err)
	}

	links, err := netlink.LinkList()
	if err != nil {
		return "", fmt.Errorf("failed to get net links in namespace %s: %v", ns, err)
	}

	return printTable(links), nil
}

func printTable(links []netlink.Link) string {
	var b strings.Builder

	ifWidth, stateWidth, ipWidth, macWidth := 18, 5, 18, 17

	b.WriteString("┌" + strings.Repeat("─", ifWidth) + "┬" + strings.Repeat("─", stateWidth) + "┬" + strings.Repeat("─", ipWidth) + "┬" + strings.Repeat("─", macWidth) + "┐\n")
	b.WriteString(fmt.Sprintf("│%-*.*s│%-*.*s│%-*.*s│%-*.*s│\n", ifWidth, ifWidth, "IFACE", stateWidth, stateWidth, "STATE", ipWidth, ipWidth, "IP", macWidth, macWidth, "MAC"))
	b.WriteString("├" + strings.Repeat("─", ifWidth) + "┼" + strings.Repeat("─", stateWidth) + "┼" + strings.Repeat("─", ipWidth) + "┼" + strings.Repeat("─", macWidth) + "┤\n")

	for _, link := range links {
		attrs := link.Attrs()
		if attrs == nil {
			continue
		}

		if attrs.Name == "gre0" || attrs.Name == "gretap0" || attrs.Name == "erspan0" || attrs.Name == "ip6gre0" || attrs.Name == "ip6gretap0" {
			continue
		}

		state := attrs.OperState.String()
		if state == "" {
			state = "-"
		}

		ip := "-"
		if addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL); err == nil && len(addrs) > 0 {
			ipList := make([]string, 0, len(addrs))
			for _, addr := range addrs {
				if addr.IPNet != nil {
					ipList = append(ipList, addr.IP.String())
				}
			}
			if len(ipList) > 0 {
				ip = strings.Join(ipList, ",")
			}
		}

		mac := attrs.HardwareAddr.String()
		if mac == "" {
			mac = "-"
		}

		b.WriteString(fmt.Sprintf("│%-*.*s│%-*.*s│%-*.*s│%-*.*s│\n", ifWidth, ifWidth, attrs.Name, stateWidth, stateWidth, state, ipWidth, ipWidth, ip, macWidth, macWidth, mac))

	}

	b.WriteString("└" + strings.Repeat("─", ifWidth) + "┴" + strings.Repeat("─", stateWidth) + "┴" + strings.Repeat("─", ipWidth) + "┴" + strings.Repeat("─", macWidth) + "┘\n")

	return b.String()
}
