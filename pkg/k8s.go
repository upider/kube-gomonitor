package pkg

import (
	"k8s.io/client-go/rest"
)

//GetKubeConfig get kubernetes config
func GetKubeConfig() (*rest.Config, error) {
	var kubeConfig *rest.Config
	var err error

	// creates the in-cluster config
	kubeConfig, err = rest.InClusterConfig()
	if err == nil {
		return kubeConfig, nil
	} else {
		return nil, err
	}
}
