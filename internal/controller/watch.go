package controller

import (
	"bytes"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"net"
)

func (c *Controller) Watch(stopper chan struct{}) {
	factory := informers.NewSharedInformerFactory(c.client, 0)
	informer := factory.Core().V1().Services().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.add(obj.(*corev1.Service))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.remove(oldObj.(*corev1.Service))
			c.add(newObj.(*corev1.Service))
		},
		DeleteFunc: func(obj interface{}) {
			c.remove(obj.(*corev1.Service))
		},
	})

	go informer.Run(stopper)
}

func (c *Controller) remove(svc *corev1.Service) {
	name := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
	address, ok := c.serviceAddressMap[name]
	if !ok {
		return
	}
	host := c.addressHostMap[address]
	delete(c.serviceAddressMap, name)
	delete(c.addressHostMap, address)
	delete(c.hostAddressMap, host)
}

func (c *Controller) add(svc *corev1.Service) {
	name := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
	if c.versionCache[name] == svc.ResourceVersion {
		return
	}

	if svc.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return
	}

	domain, ok := svc.Annotations["metallb-dns/domain"]
	if !ok {
		return
	}

	if len(svc.Status.LoadBalancer.Ingress) != 1 {
		return
	}

	found := false
	ingress := svc.Status.LoadBalancer.Ingress[0]
	ip := net.ParseIP(ingress.IP)
	for _, ranges := range c.ipRanges {
		if bytes.Compare(ip, ranges[0]) >= 0 && bytes.Compare(ip, ranges[1]) <= 0 {
			if c.addressHostMap[ingress.IP] == domain {
				return
			}
			found = true
			break
		}
	}

	if !found {
		return
	}
	c.serviceAddressMap[name] = ingress.IP
	c.addressHostMap[ingress.IP] = domain
	c.hostAddressMap[domain] = ingress.IP
	if ingress.Hostname == domain {
		return
	}
	svc.Status = corev1.ServiceStatus{
		LoadBalancer: corev1.LoadBalancerStatus{
			Ingress: []corev1.LoadBalancerIngress{{
				IP:       ip.String(),
				Hostname: domain,
			}},
		},
	}
	result, err := c.client.CoreV1().Services(svc.Namespace).UpdateStatus(svc);
	if err != nil {
		fmt.Printf("unable to update ignress status of svc %s: %s\n", name, err.Error())
	} else {
		fmt.Printf("successfully updated ignress status of svc %s\n", name)
		c.versionCache[name] = result.ResourceVersion
	}
}
