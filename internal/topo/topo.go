package topo

import (
	"fmt"
	"net"

	"github.com/Alonza0314/nsctl/internal/namespace"
	"github.com/pterm/pterm"
)

func checkTopo(topo *Topology) error {
	spinner, err := pterm.DefaultSpinner.Start("Checking topology namespaces...")
	if err != nil {
		return err
	}
	if err := checkTopoNamespace(topo); err != nil {
		return err
	}
	spinner.Success("Topology namespaces check passed")

	spinner, err = pterm.DefaultSpinner.Start("Checking topology networks...")
	if err != nil {
		return err
	}
	if err := checkTopoNetwork(topo); err != nil {
		return err
	}
	spinner.Success("Topology networks check passed")

	return nil
}

func checkTopoNamespace(topo *Topology) error {
	nsNames := make(map[string]struct{})

	for _, ns := range topo.Namespaces {
		if _, exists := nsNames[ns.Name]; exists {
			return fmt.Errorf("duplicate namespace name: %s", ns.Name)
		}
		nsNames[ns.Name] = struct{}{}

		if err := checkNamespaceNetwork(&ns); err != nil {
			return err
		}

	}

	return nil
}

func checkNamespaceNetwork(ns *Namespace) error {
	netNames, ips := make(map[string]struct{}), make(map[string]struct{})

	for _, network := range ns.Networks {
		if _, exists := netNames[network.Name]; exists {
			return fmt.Errorf("duplicate network name in namespace %s: %s", ns.Name, network.Name)
		}
		netNames[network.Name] = struct{}{}

		if _, exists := ips[network.Ipv4]; exists {
			return fmt.Errorf("duplicate IP address in namespace %s: %s", ns.Name, network.Ipv4)
		}
		ips[network.Ipv4] = struct{}{}
	}

	return nil
}

func checkTopoNetwork(topo *Topology) error {
	bridges, subnets, bridgeToSubnet := make(map[string][]Net), make(map[string]struct{}), make(map[string]*net.IPNet)

	for _, network := range topo.Networks {
		if _, exists := bridges[network.Name]; exists {
			return fmt.Errorf("duplicate network name: %s", network.Name)
		}

		_, subnet, err := net.ParseCIDR(network.Subnet)
		if err != nil {
			return fmt.Errorf("invalid subnet CIDR for network %s: %s", network.Name, network.Subnet)
		}

		if _, exists := subnets[subnet.String()]; exists {
			return fmt.Errorf("duplicate subnet: %s", network.Subnet)
		}

		bridges[network.Name], subnets[subnet.String()], bridgeToSubnet[network.Name] = make([]Net, 0), struct{}{}, subnet
	}

	for _, ns := range topo.Namespaces {
		for _, network := range ns.Networks {
			if _, exists := bridges[network.Bridge]; !exists {
				return fmt.Errorf("network bridge %s in namespace %s is not defined in topo networks", network.Bridge, ns.Name)
			}
			bridges[network.Bridge] = append(bridges[network.Bridge], network)
		}
	}

	for bdg, nets := range bridges {
		ips := make(map[string]struct{})
		for _, ip := range nets {
			if _, exists := ips[ip.Ipv4]; exists {
				return fmt.Errorf("duplicate IP address %s for bridge %s", ip.Ipv4, bdg)
			}
			if err := checkSubnet(bridgeToSubnet[bdg], ip.Ipv4); err != nil {
				return fmt.Errorf("invalid IP address %s for bridge %s: %v", ip.Ipv4, bdg, err)
			}
			ips[ip.Ipv4] = struct{}{}
		}
	}

	return nil
}

func checkSubnet(subnet *net.IPNet, targetIp string) error {
	ip, _, err := net.ParseCIDR(targetIp)
	if err != nil {
		return fmt.Errorf("invalid IP address: %s", targetIp)
	}

	if !subnet.Contains(ip) {
		return fmt.Errorf("IPv4 address %s is not in subnet %s", targetIp, subnet)
	}

	return nil
}

func checkExist(topo *Topology) error {
	if err := checkNamespaceExist(topo.Namespaces); err != nil {
		return err
	}

	return nil
}

func checkNamespaceExist(nss []Namespace) error {
	for _, ns := range nss {
		found, err := namespace.GetNs(namespace.GetNsName(ns.Name))
		if err != nil {
			return err
		}
		if found {
			return fmt.Errorf("namespace %s already exists", ns.Name)
		}
	}

	return nil
}
