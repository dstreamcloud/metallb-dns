package main

import (
	"fmt"
	"github.com/dstreamcloud/metallb-dns/internal/controller"
	"github.com/dstreamcloud/metallb-dns/internal/dns"
	"github.com/dstreamcloud/metallb-dns/internal/version"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
)

func main() {
	println(version.String())
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	metalLBConfig := &controller.MetalLBConfig{}

	{
		cm, err := clientset.CoreV1().ConfigMaps("metallb-system").Get("config", corev1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		if err := yaml.Unmarshal([]byte(cm.Data["config"]), metalLBConfig); err != nil {
			panic(err.Error())
		}
	}

	ctrl := controller.New(clientset, metalLBConfig)
	stopper := make(chan struct{})
	ctrl.Watch(stopper)

	dnsServer := dns.New(ctrl)
	go func() {
		if err := dnsServer.StartTCP(); err != nil {
			fmt.Println("unable to start DNS server at TCP port 53: ", err.Error())
			stopper <- struct{}{}
		}
	}()
	go func() {
		if err := dnsServer.StartUDP(); err != nil {
			fmt.Println("unable to start DNS server at UDP port 53: ", err.Error())
			stopper <- struct{}{}
		}
	}()

	<-stopper
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
