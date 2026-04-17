package topo

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Alonza0314/nsctl/internal/namespace"
	"github.com/Alonza0314/nsctl/internal/veth"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

func Delete(topo *Topology) error {
	if err := checkTopo(topo); err != nil {
		return err
	}

	if err := deleteNamespaces(topo.Namespaces); err != nil {
		return fmt.Errorf("failed to delete namespaces: %v", err)
	}

	if err := deleteBridges(topo.Networks); err != nil {
		return fmt.Errorf("failed to delete bridges: %v", err)
	}

	return nil
}

func deleteNamespaces(nss []Namespace) error {
	var errBuilder []error
	for _, ns := range nss {
		found, err := namespace.GetNs(ns.Name)
		if err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to get namespace %s: %v", ns.Name, err))
			continue
		}
		if !found {
			errBuilder = append(errBuilder, fmt.Errorf("namespace %s does not exist", ns.Name))
			continue
		}

		_, originCloseFunc, err := namespace.GetOriginNs()
		if err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to get origin namespace file descriptor: %v", err))
			continue
		}
		nsFd, nsCloseFunc, err := namespace.GetNsFd(ns.Name)
		if err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to get namespace %s file descriptor: %v", ns.Name, err))
			originCloseFunc()
			continue
		}

		if err := netns.Set(nsFd); err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to set namespace %s: %v", ns.Name, err))
			nsCloseFunc()
			originCloseFunc()
			continue
		}

		for _, network := range ns.Networks {
			if err := veth.UpDown(ns.Name, network.Name, false); err != nil {
				errBuilder = append(errBuilder, fmt.Errorf("failed to bring down veth for namespace %s and network %s: %v", ns.Name, network.Name, err))
				continue
			}

			link, err := netlink.LinkByName(network.Name)
			if err != nil {
				errBuilder = append(errBuilder, fmt.Errorf("failed to get link %s in namespace %s: %v", network.Name, ns.Name, err))
				continue
			}
			if err := netlink.LinkDel(link); err != nil {
				errBuilder = append(errBuilder, fmt.Errorf("failed to delete link %s in namespace %s: %v", network.Name, ns.Name, err))
			}
		}

		nsCloseFunc()
		originCloseFunc()

		if err := namespace.Delete(ns.Name); err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to delete namespace %s: %v", ns.Name, err))
		}
	}

	return errBuild(errBuilder)
}

func deleteBridges(networks []Network) error {
	var errBuilder []error
	for _, network := range networks {
		link, err := netlink.LinkByName(network.Name)
		if err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to get bridge %s: %v", network.Name, err))
			continue
		}
		if err := netlink.LinkDel(link); err != nil {
			errBuilder = append(errBuilder, fmt.Errorf("failed to delete bridge %s: %v", network.Name, err))
		}
	}

	return errBuild(errBuilder)
}

func errBuild(errBuilder []error) error {
	if len(errBuilder) == 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("Multiple errors occurred:\n")
	for _, err := range errBuilder {
		fmt.Fprintf(&b, "- %v\n", err)
	}

	return errors.New(b.String())
}
