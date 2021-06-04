package webhooks

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type SidecarHandler struct {
	AgentImage string
	Url        string
	Bucket     string
	Org        string
	Token      string
}

var (
	GoMonitor                      = "monitor.online.daq.ihep/monitor"
	GoMonitorAgentImage            = "monitor.online.daq.ihep/image"
	GoMonitorServiceName           = "monitor.online.daq.ihep/serviceName"
	GoMonitorHostPid               = "monitor.online.daq.ihep/hostPid"
	GoMonitorInterval              = "monitor.online.daq.ihep/interval"
	GoMonitorShareProcessNamespace = "monitor.online.daq.ihep/shareProcessNamespace"
	REPORT_DBURL                   = "REPORT_DBURL"
	REPORT_DBBUCKET                = "REPORT_DBBUCKET"
	REPORT_DBORG                   = "REPORT_DBORG"
	REPORT_DBTOKEN                 = "REPORT_DBTOKEN"
	NODE_NAME                      = "NODE_NAME"
	NODE_IP                        = "NODE_IP"
	POD_NAME                       = "POD_NAME"
	POD_IP                         = "POD_IP"
	MONITOR_IP                     = "MONITOR_IP"
	MONITOR_SERVICE                = "MONITOR_SERVICE"
)

func (h *SidecarHandler) shouldAddSidecar(pod *corev1.Pod) bool {
	handlerLog := ctrl.Log.WithName("inject-handler")

	if podHasContainerName(pod, "gomonitor-agent") {
		handlerLog.Info("pod already has a gomonitor-agent", "podname", pod.Name)
		return false
	}
	anno, found := pod.Annotations[GoMonitor]
	if !found {
		handlerLog.Info("pod has not a GoMonitor annotation", "podname", pod.Name)
		return false
	}

	if anno != "true" {
		handlerLog.Info("pod's GoMonitor annotation is not true", "podname", pod.Name)
		return false
	}

	return true
}

func (h *SidecarHandler) addSidecar(pod *corev1.Pod, containerName string) error {
	//设置agent image
	var agentImage string
	if customImage, ok := pod.Annotations[GoMonitorAgentImage]; ok {
		agentImage = customImage
	} else {
		agentImage = h.AgentImage
	}

	//修改Pod shareProcessNamespace, HostPid级别高于ShareProcessNamespace
	var yes bool = true
	if hostPid, ok := pod.Annotations[GoMonitorHostPid]; ok {
		yes, _ = strconv.ParseBool(hostPid)
		pod.Spec.HostPID = yes
	} else if shareProcessNamespace, ok := pod.Annotations[GoMonitorShareProcessNamespace]; ok {
		yes, _ = strconv.ParseBool(shareProcessNamespace)
		pod.Spec.ShareProcessNamespace = &yes
	} else {
		pod.Spec.HostPID = yes
	}

	//修改SecurityContext, 这样才能共享命名空间
	agentSecurityContext := corev1.SecurityContext{
		Capabilities: &corev1.Capabilities{
			Add:  []corev1.Capability{corev1.Capability("SYS_PTRACE")},
			Drop: []corev1.Capability{},
		},
	}

	//创建container
	var envs []corev1.EnvVar
	envs = append(envs, corev1.EnvVar{Name: REPORT_DBURL, Value: h.Url})
	envs = append(envs, corev1.EnvVar{Name: REPORT_DBBUCKET, Value: h.Bucket})
	envs = append(envs, corev1.EnvVar{Name: REPORT_DBORG, Value: h.Org})
	envs = append(envs, corev1.EnvVar{Name: REPORT_DBTOKEN, Value: h.Token})

	envs = append(envs, corev1.EnvVar{Name: NODE_NAME, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "spec.nodeName"}}})
	envs = append(envs, corev1.EnvVar{Name: NODE_IP, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.hostIP"}}})
	envs = append(envs, corev1.EnvVar{Name: POD_NAME, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}}})
	envs = append(envs, corev1.EnvVar{Name: POD_IP, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}}})

	envs = append(envs, corev1.EnvVar{Name: MONITOR_IP, ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}}})
	//查看servicename
	if serviceName, ok := pod.Annotations[GoMonitorServiceName]; ok {
		envs = append(envs, corev1.EnvVar{Name: MONITOR_SERVICE, Value: serviceName})
	} else {
		envs = append(envs, corev1.EnvVar{Name: MONITOR_SERVICE, Value: "DefaultServiceName"})
	}

	agentContainer := corev1.Container{
		Name:            containerName,
		Image:           agentImage,
		SecurityContext: &agentSecurityContext,
		Env:             envs,
		ImagePullPolicy: corev1.PullAlways,
	}

	pod.Spec.Containers = append(pod.Spec.Containers, agentContainer)
	return nil
}

func podHasContainerName(pod *corev1.Pod, name string) bool {
	for _, container := range pod.Spec.Containers {
		if container.Name == name {
			return true
		}
	}
	return false
}
