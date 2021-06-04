package webhooks

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/storage/names"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

// podInjector inject gomonitor Pods
type PodInjector struct {
	Client  client.Client
	Handler *SidecarHandler
	decoder *admission.Decoder
	Logger  logr.Logger
}

// podInjector adds an container to every incoming pods.
func (a *PodInjector) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if a.Handler.shouldAddSidecar(pod) {
		name := pod.GetName()
		if name == "" {
			name = names.SimpleNameGenerator.GenerateName(pod.GetGenerateName())
			pod.SetName(name)
		}
		a.Logger.Info("inject gomonitor agent ", "podName", pod.GetName(), "namespace", req.Namespace)
		a.Handler.addSidecar(pod, "gomonitor-agent")

		marshaledPod, err := json.Marshal(pod)
		if err != nil {
			a.Logger.Error(err, "Pod Marshal Error")
			return admission.Errored(http.StatusInternalServerError, err)
		}
		return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)

	} else {
		a.Logger.Info("skipping pod as gomonitor-injector should not handle it")
		return admission.Allowed("gomonitor-injector has no power over this pod")
	}
}

// podInjector implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (a *PodInjector) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
