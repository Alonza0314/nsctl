package topo

type Topology struct {
	Namespaces []Namespace `yaml:"namespaces" validate:"required"`
	Networks   []Network   `yaml:"networks" validate:"required"`
}

type Namespace struct {
	Name     string   `yaml:"name" validate:"required"`
	Commands []string `yaml:"commands"`
	Networks []Net    `yaml:"networks"`
}

type Network struct {
	Name   string `yaml:"name" validate:"required"`
	Subnet string `yaml:"subnet" validate:"required,ipv4cidr"`
}

type Net struct {
	Name   string `yaml:"name" validate:"required"`
	Bridge string `yaml:"bridge" validate:"required"`
	Ipv4   string `yaml:"ipv4" validate:"required,ipv4cidr"`
}
