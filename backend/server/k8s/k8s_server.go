package k8s

import (
	"context"
	"gomonitor/backend/server/k8s/controllers"
	"gomonitor/backend/server/k8s/validatorwebhook"
	"os"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type KServer struct {
	url    string
	bucket string
	org    string
	token  string
}

var (
	// scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	log.SetLogger(zap.New())
}

func (server *KServer) Start(ctx context.Context) {
	entryLog := setupLog.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up webhook server")
	mgr, err := manager.New(ctrl.GetConfigOrDie(), ctrl.Options{
		// Scheme:  scheme,
		CertDir: "/etc/certs"})
	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup a new controller to reconcile ReplicaSets
	entryLog.Info("Setting up controller")
	c, err := controller.New("gomonitor-controller", mgr, controller.Options{
		Reconciler: &controllers.PodReconciler{Client: mgr.GetClient(),
			Handler: &controllers.SidecarHandler{
				AgentImage: "1445277435/gomonitor-agent:v0.0.1",
				Url:        server.url,
				Bucket:     server.bucket,
				Org:        server.org,
				Token:      server.token,
			}},
	})
	if err != nil {
		entryLog.Error(err, "unable to set up individual controller")
		os.Exit(1)
	}

	// Watch Pods
	if err := c.Watch(&source.Kind{Type: &corev1.Pod{}},
		&handler.EnqueueRequestForOwner{OwnerType: &corev1.Pod{}, IsController: true}); err != nil {
		entryLog.Error(err, "unable to watch Pods")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/validate-v1-pod", &webhook.Admission{Handler: &validatorwebhook.PodValidator{Client: mgr.GetClient()}})

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func NewKServer(url string, bucket string, org string, token string) *KServer {
	return &KServer{
		url:    url,
		bucket: bucket,
		org:    org,
		token:  token,
	}
}
