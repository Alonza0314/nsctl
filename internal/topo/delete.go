package topo

import (
	"github.com/Alonza0314/nsctl/internal/namespace"
	"github.com/Alonza0314/nsctl/internal/veth"
	"github.com/pterm/pterm"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func Delete(topo *Topology) error {
	g, err := checkTopo(topo)
	if err != nil {
		return err
	}

	if err := deleteNamespaces(g.getSortedNamespaces(topo.Namespaces, true)); err != nil {
		return err
	}

	if err := deleteBridges(topo.Networks); err != nil {
		return err
	}

	return nil
}

func deleteNamespaces(nss []Namespace) error {
	for _, ns := range nss {
		spinner, err := pterm.DefaultSpinner.Start("Deleting namespace " + ns.Name + "...")
		if err != nil {
			spinner.Fail("Failed to start spinner for namespace " + ns.Name + ": " + err.Error())
			continue
		}

		found, err := namespace.GetNs(ns.Name)
		if err != nil {
			spinner.Fail("Failed to get namespace " + ns.Name + ": " + err.Error())
			continue
		}
		if !found {
			spinner.Fail("Namespace " + ns.Name + " does not exist")
			continue
		}

		_, originCloseFunc, err := namespace.GetOriginNs()
		if err != nil {
			spinner.Fail("Failed to get origin namespace file descriptor: " + err.Error())
			continue
		}
		nsFd, nsCloseFunc, err := namespace.GetNsFd(ns.Name)
		if err != nil {
			spinner.Fail("Failed to get namespace " + ns.Name + " file descriptor: " + err.Error())
			originCloseFunc()
			continue
		}

		if err := netns.Set(nsFd); err != nil {
			spinner.Fail("Failed to set namespace " + ns.Name + ": " + err.Error())
			nsCloseFunc()
			originCloseFunc()
			continue
		}

		for _, network := range ns.Networks {
			if err := veth.UpDown(ns.Name, network.Name, false); err != nil {
				spinner.Fail("Failed to bring down veth for namespace " + ns.Name + " and network " + network.Name + ": " + err.Error())
				continue
			}

			link, err := netlink.LinkByName(network.Name)
			if err != nil {
				spinner.Fail("Failed to get link " + network.Name + " in namespace " + ns.Name + ": " + err.Error())
				continue
			}
			if err := netlink.LinkDel(link); err != nil {
				spinner.Fail("Failed to delete link " + network.Name + " in namespace " + ns.Name + ": " + err.Error())
			}
		}

		nsCloseFunc()
		originCloseFunc()

		if err := namespace.Delete(ns.Name); err != nil {
			spinner.Fail("Failed to delete namespace " + ns.Name + ": " + err.Error())
		}
		spinner.Success("Namespace " + ns.Name + " deleted")
	}

	return nil
}

func deleteBridges(networks []Network) error {
	for _, network := range networks {
		spinner, err := pterm.DefaultSpinner.Start("Deleting bridge " + network.Name + "...")
		if err != nil {
			spinner.Fail("Failed to start spinner for bridge " + network.Name + ": " + err.Error())
			continue
		}

		link, err := netlink.LinkByName(network.Name)
		if err != nil {
			spinner.Fail("Failed to get bridge " + network.Name + ": " + err.Error())
			continue
		}
		if err := netlink.LinkDel(link); err != nil {
			spinner.Fail("Failed to delete bridge " + network.Name + ": " + err.Error())
		}
		spinner.Success("Bridge " + network.Name + " deleted")
	}

	return nil
}
