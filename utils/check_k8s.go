package utils

import (
	"k8s.io/client-go/rest"
)

//CheckK8s check if running on k8s
func CheckInK8s() (*rest.Config, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		return nil, err
	}

	if err != nil {
		panic(err.Error())
	}

	return config, nil
}
