package topo

import "testing"

var testCheckTopoCases = []struct {
	name              string
	topo              *Topology
	expectedErr       bool
	expectedErrDetail string
}{
	{
		name: "valid topo",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
					DependsOn: []string{"test-2"},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       false,
		expectedErrDetail: "",
	},
	{
		name: "duplicate ns name",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate namespace name: test-1",
	},
	{
		name: "duplicate network name",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.3/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate network name in namespace test-1: br-test-1",
	},
	{
		name: "duplicate IP address",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
						{
							Name:   "br-test-3",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate IP address in namespace test-1: 10.0.0.1/24",
	},
	{
		name: "duplicate network name",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
				{
					Name:   "br-test-1-2",
					Subnet: "11.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate network name: br-test-1-2",
	},
	{
		name: "invalid subnet CIDR",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/100",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "invalid subnet CIDR for network br-test-1-2: 10.0.0.0/100",
	},
	{
		name: "duplicate subnet",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
				{
					Name:   "br-test-3-4",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate subnet: 10.0.0.0/24",
	},
	{
		name: "duplicate IP address across namespaces",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "duplicate IP address 10.0.0.1/24 for bridge br-test-1-2",
	},
	{
		name: "invalid subnet for network bridge",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.1.1/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "invalid IP address 10.0.1.1/24 for bridge br-test-1-2: IPv4 address 10.0.1.1/24 is not in subnet 10.0.0.0/24",
	},
	{
		name: "invalid IP address",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.999/24",
						},
					},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "invalid IP address 10.0.0.999/24 for bridge br-test-1-2: invalid IP address: 10.0.0.999/24",
	},
	{
		name: "depends on non-existent namespace",
		topo: &Topology{
			Namespaces: []Namespace{
				{
					Name:     "test-1",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-1",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.1/24",
						},
					},
					DependsOn: []string{"test-3"},
				},
				{
					Name:     "test-2",
					Commands: []string{"echo"},
					Networks: []Net{
						{
							Name:   "br-test-2",
							Bridge: "br-test-1-2",
							Ipv4:   "10.0.0.2/24",
						},
					},
				},
			},
			Networks: []Network{
				{
					Name:   "br-test-1-2",
					Subnet: "10.0.0.0/24",
				},
			},
		},
		expectedErr:       true,
		expectedErrDetail: "namespace test-1 depends on non-existent namespace: test-3",
	},
}

func TestCheckTopo(t *testing.T) {
	for _, tc := range testCheckTopoCases {
		t.Run(tc.name, func(t *testing.T) {
			err := checkTopo(tc.topo)
			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tc.expectedErrDetail {
					t.Errorf("expected error detail '%s' but got '%s'", tc.expectedErrDetail, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}
