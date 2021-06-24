module kube-gomonitor

go 1.15

require (
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/go-logr/logr v0.4.0
	github.com/google/gopacket v1.1.19
	github.com/influxdata/influxdb-client-go/v2 v2.4.0
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/shirou/gopsutil v3.21.5+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/tklauser/go-sysconf v0.3.6 // indirect
	golang.org/x/net v0.0.0-20210428140749-89ef3d95e781
	k8s.io/api v0.21.2
	k8s.io/apimachinery v0.21.2
	k8s.io/apiserver v0.21.2
	k8s.io/client-go v0.21.2
	sigs.k8s.io/controller-runtime v0.9.1
)
