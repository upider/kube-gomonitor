package main

import (
	"kube-gomonitor/pkg"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/client-go/rest"
)

func main() {
	_, err := pkg.GetKubeConfig()
	if err == rest.ErrNotInCluster {
		bareRegistry()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go busy(signalChan)
	<-signalChan
}
