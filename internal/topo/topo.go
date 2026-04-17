package topo

import (
	"fmt"
	"net"

	"github.com/Alonza0314/nsctl/internal/namespace"
	"github.com/pterm/pterm"
)

func checkTopo(topo *Topology) (*graph, error) {
	spinner, err := pterm.DefaultSpinner.Start("Checking topology namespaces...")
	if err != nil {
		return nil, err
	}
	g, err := checkTopoNamespace(topo)
	if err != nil {
		return nil, err
	}
	spinner.Success("Topology namespaces check passed")

	spinner, err = pterm.DefaultSpinner.Start("Checking topology networks...")
	if err != nil {
		return nil, err
	}
	if err := checkTopoNetwork(topo); err != nil {
		return nil, err
	}
	spinner.Success("Topology networks check passed")

	return g, nil
}

func checkTopoNamespace(topo *Topology) (*graph, error) {
	nsNames := make(map[string]struct{})

	for _, ns := range topo.Namespaces {
		if _, exists := nsNames[ns.Name]; exists {
			return nil, fmt.Errorf("duplicate namespace name: %s", ns.Name)
		}
		nsNames[ns.Name] = struct{}{}

		if err := checkNamespaceNetwork(&ns); err != nil {
			return nil, err
		}
	}

	for _, ns := range topo.Namespaces {
		for _, dep := range ns.DependsOn {
			if _, exists := nsNames[dep]; !exists {
				return nil, fmt.Errorf("namespace %s depends on non-existent namespace: %s", ns.Name, dep)
			}
		}
	}

	g, err := existCycle(topo.Namespaces)
	if err != nil {
		return nil, err
	}

	return g, nil
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

type graph struct {
	aTob     map[string][]string
	inDegree map[string]int
	index    map[string]int
	sorted   []string
}

func newGraph() *graph {
	return &graph{
		aTob:     make(map[string][]string),
		inDegree: make(map[string]int),
		index:    make(map[string]int),
		sorted:   nil,
	}
}

func (g *graph) make(nss []Namespace) {
	for i, ns := range nss {
		g.aTob[ns.Name], g.inDegree[ns.Name], g.index[ns.Name] = make([]string, 0), 0, i
	}

	for _, ns := range nss {
		for _, dep := range ns.DependsOn {
			g.aTob[dep] = append(g.aTob[dep], ns.Name)
			g.inDegree[ns.Name] += 1
		}
	}

}

func (g *graph) getNeighbors(ns string) []string {
	return g.aTob[ns]
}

func (g *graph) topologicalSort() {
	queue := make([]string, 0)
	for ns, degree := range g.inDegree {
		if degree == 0 {
			queue = append(queue, ns)
		}
	}

	sorted := make([]string, 0)
	for len(queue) > 0 {
		ns := queue[0]
		queue, sorted = queue[1:], append(sorted, ns)

		for _, neighbor := range g.getNeighbors(ns) {
			g.inDegree[neighbor] -= 1
			if g.inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	g.sorted = sorted
}

func (g *graph) getSortedNamespaces(nss []Namespace, reverse bool) []Namespace {
	nsList := make([]Namespace, len(nss))
	for i, ns := range g.sorted {
		nsList[i] = nss[g.index[ns]]
	}

	if reverse {
		for i, j := 0, len(nsList)-1; i < j; i, j = i+1, j-1 {
			nsList[i], nsList[j] = nsList[j], nsList[i]
		}
	}

	return nsList
}

func existCycle(nss []Namespace) (*graph, error) {
	g := newGraph()
	g.make(nss)

	g.topologicalSort()
	if len(g.sorted) != len(nss) {
		return nil, fmt.Errorf("circular dependency detected among namespaces")
	}

	return g, nil
}
