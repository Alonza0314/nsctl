package topo

func Apply(topo *Topology) error {
	if err := checkTopo(topo); err != nil {
		return err
	}

	return nil
}