package topo

import (
	"errors"

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

	errFlag := false
	if err := deleteNamespaces(g.getSortedNamespaces(topo.Namespaces, true)); err != nil {
		errFlag = true
	}

	if err := deleteBridges(topo.Networks); err != nil {
		errFlag = true
	}

	if errFlag {
		return errors.New("Error")
	}
	return nil
}

func deleteNamespaces(nss []Namespace) error {
	errFlag := false

	for _, ns := range nss {
		spinner, err := pterm.DefaultSpinner.Start("Deleting namespace " + ns.Name + "...")
		if err != nil {
			spinner.Fail("Failed to start spinner for namespace " + ns.Name + ": " + err.Error())
			errFlag = true
			continue
		}

		found, err := namespace.GetNs(ns.Name)
		if err != nil {
			spinner.Fail("Failed to get namespace " + ns.Name + ": " + err.Error())
			errFlag = true
			continue
		}
		if !found {
			spinner.Fail("Namespace " + ns.Name + " does not exist")
			errFlag = true
			continue
		}

		_, originCloseFunc, err := namespace.GetOriginNs()
		if err != nil {
			spinner.Fail("Failed to get origin namespace file descriptor: " + err.Error())
			errFlag = true
			continue
		}
		nsFd, nsCloseFunc, err := namespace.GetNsFd(ns.Name)
		if err != nil {
			spinner.Fail("Failed to get namespace " + ns.Name + " file descriptor: " + err.Error())
			originCloseFunc()
			errFlag = true
			continue
		}

		if err := netns.Set(nsFd); err != nil {
			spinner.Fail("Failed to set namespace " + ns.Name + ": " + err.Error())
			nsCloseFunc()
			originCloseFunc()
			errFlag = true
			continue
		}

		for _, network := range ns.Networks {
			if err := veth.UpDown(ns.Name, network.Name, false); err != nil {
				spinner.Fail("Failed to bring down veth for namespace " + ns.Name + " and network " + network.Name + ": " + err.Error())
				errFlag = true
				continue
			}

			link, err := netlink.LinkByName(network.Name)
			if err != nil {
				spinner.Fail("Failed to get link " + network.Name + " in namespace " + ns.Name + ": " + err.Error())
				errFlag = true
				continue
			}
			if err := netlink.LinkDel(link); err != nil {
				spinner.Fail("Failed to delete link " + network.Name + " in namespace " + ns.Name + ": " + err.Error())
				errFlag = true
			}
		}

		nsCloseFunc()
		originCloseFunc()

		if err := namespace.Delete(ns.Name); err != nil {
			spinner.Fail("Failed to delete namespace " + ns.Name + ": " + err.Error())
			errFlag = true
		}
		spinner.Success("Namespace " + ns.Name + " deleted")
	}

	if errFlag {
		return errors.New("Error")
	}
	return nil
}

func deleteBridges(networks []Network) error {
	errFlag := false

	for _, network := range networks {
		spinner, err := pterm.DefaultSpinner.Start("Deleting bridge " + network.Name + "...")
		if err != nil {
			spinner.Fail("Failed to start spinner for bridge " + network.Name + ": " + err.Error())
			errFlag = true
			continue
		}

		link, err := netlink.LinkByName(network.Name)
		if err != nil {
			spinner.Fail("Failed to get bridge " + network.Name + ": " + err.Error())
			errFlag = true
			continue
		}
		if err := netlink.LinkDel(link); err != nil {
			spinner.Fail("Failed to delete bridge " + network.Name + ": " + err.Error())
			errFlag = true
		}
		spinner.Success("Bridge " + network.Name + " deleted")
	}

	if errFlag {
		return errors.New("Error")
	}
	return nil
}
