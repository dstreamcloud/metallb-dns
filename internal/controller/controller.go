package controller

import (
	"k8s.io/client-go/kubernetes"
	"net"
	"strings"
)

type Controller struct {
	client            *kubernetes.Clientset
	serviceAddressMap map[string]string
	hostAddressMap    map[string]string
	addressHostMap    map[string]string
	versionCache      map[string]string

	ipRanges [][]net.IP
}

func New(client *kubernetes.Clientset, metalLBConfig *MetalLBConfig) *Controller {
	var ipRanges [][]net.IP
	for _, pool := range metalLBConfig.AddressPools {
		for _, address := range pool.Addresses {
			ips := strings.Split(address, "-")
			ipRanges = append(ipRanges, []net.IP{net.ParseIP(ips[0]), net.ParseIP(ips[1])})
		}
	}

	return &Controller{
		client:            client,
		serviceAddressMap: make(map[string]string),
		hostAddressMap:    make(map[string]string),
		addressHostMap:    make(map[string]string),
		versionCache:      make(map[string]string),
		ipRanges:          ipRanges,
	}
}

func (c *Controller) Resolve(domain string) string {
	return c.hostAddressMap[domain]
}
