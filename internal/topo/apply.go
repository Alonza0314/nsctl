package topo

import (
	"github.com/Alonza0314/nsctl/internal/namespace"
	"github.com/Alonza0314/nsctl/internal/veth"
	"github.com/pterm/pterm"
	"github.com/vishvananda/netlink"
)

func Apply(topo *Topology) error {
	g, err := checkTopo(topo)
	if err != nil {
		return err
	}

	if err := checkExist(topo); err != nil {
		return err
	}

	if err := addBridges(topo.Networks); err != nil {
		return err
	}

	if err := addNamespaces(g.getSortedNamespaces(topo.Namespaces, false)); err != nil {
		return err
	}

	return nil
}

func addBridges(networks []Network) error {
	for _, network := range networks {
		spinner, err := pterm.DefaultSpinner.Start("Creating bridge " + network.Name + "...")
		if err != nil {
			return err
		}

		br := &netlink.Bridge{
			LinkAttrs: netlink.LinkAttrs{
				Name: network.Name,
			},
		}

		if err := netlink.LinkAdd(br); err != nil {
			spinner.Fail("Failed to create bridge " + network.Name + ": " + err.Error())
			return err
		}

		if err := netlink.LinkSetUp(br); err != nil {
			spinner.Fail("Failed to set bridge " + network.Name + " up: " + err.Error())
			return err
		}

		spinner.Success("Bridge " + network.Name + " created")
	}

	return nil
}

func addNamespaces(namespaces []Namespace) error {
	for _, ns := range namespaces {
		spinner, err := pterm.DefaultSpinner.Start("Creating namespace " + ns.Name + "...")
		if err != nil {
			return err
		}

		if err := namespace.Create(ns.Name); err != nil {
			spinner.Fail("Failed to create namespace " + ns.Name + ": " + err.Error())
			return err
		}

		if err := addVethPair(ns); err != nil {
			spinner.Fail("Failed to add veth pair for namespace " + ns.Name + ": " + err.Error())
			return err
		}

		spinner.Success("Namespace " + ns.Name + " created")
	}

	return nil
}

func addVethPair(ns Namespace) error {
	for _, network := range ns.Networks {
		vethLink := &netlink.Veth{
			LinkAttrs: netlink.LinkAttrs{
				Name: network.Name,
			},
			PeerName: "m-" + network.Name,
		}

		if err := netlink.LinkAdd(vethLink); err != nil {
			return err
		}

		brLink, err := netlink.LinkByName(network.Bridge)
		if err != nil {
			return err
		}
		vethNs, err := netlink.LinkByName(network.Name)
		if err != nil {
			return err
		}
		peer, err := netlink.LinkByName("m-" + network.Name)
		if err != nil {
			return err
		}

		if err := netlink.LinkSetMaster(peer, brLink); err != nil {
			return err
		}
		if err := netlink.LinkSetUp(peer); err != nil {
			return err
		}

		nsFd, nsCloseFunc, err := namespace.GetNsFd(ns.Name)
		if err != nil {
			return err
		}

		if err := netlink.LinkSetNsFd(vethNs, int(nsFd)); err != nil {
			return err
		}

		if err := veth.SetIp(ns.Name, network.Name, network.Ipv4); err != nil {
			return err
		}

		if err := veth.UpDown(ns.Name, network.Name, true); err != nil {
			return err
		}

		nsCloseFunc()
	}

	return nil
}
