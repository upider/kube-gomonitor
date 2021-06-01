package controllers

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

type SidecarHandler struct {
	AgentImage string
	Url        string
	Bucket     string
	Org        string
	Token      string
}

var (
	GoMonitorAgentImage            = "monitor.online.daq.ihep/image"
	GoMonitorHostPid               = "monitor.online.daq.ihep/hostPid"
	GoMonitorShareProcessNamespace = "monitor.online.daq.ihep/shareProcessNamespace"
	REPORT_DBURL                   = "REPORT_DBURL"
	REPORT_DBBUCKET                = "REPORT_DBBUCKET"
	REPORT_DBORG                   = "REPORT_DBORG"
	REPORT_DBTOKEN                 = "REPORT_DBTOKEN"
)

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

	agentContainer := corev1.Container{
		Name:            containerName,
		Image:           agentImage,
		SecurityContext: &agentSecurityContext,
		Env:             envs,
	}

	pod.Spec.Containers = append(pod.Spec.Containers, agentContainer)
	return nil
}
