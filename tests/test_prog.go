package main

import (
	"gomonitor/utils"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/client-go/rest"
)

func main() {
	_, err := utils.CheckInK8s()
	if err == rest.ErrNotInCluster {
		bareRegistry()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go busy(signalChan)
	<-signalChan
}
