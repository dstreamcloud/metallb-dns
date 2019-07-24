package controller

type MetalLBConfig struct {
	AddressPools []*AddressPool `yaml:"address-pools"`
}

type AddressPool struct {
	Name      string   `yaml:"name"`
	Protocol  string   `yaml:"protocol"`
	Addresses []string `yaml:"addresses"`
}
