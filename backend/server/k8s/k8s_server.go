package k8s

import (
	"kube-gomonitor/backend/server"
	"kube-gomonitor/backend/server/k8s/webhooks"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type KServer struct {
	url    string
	bucket string
	org    string
	token  string
}

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	log.SetLogger(zap.New())
	clientgoscheme.AddToScheme(scheme)

	// +kubebuilder:scaffold:scheme
}

func (server *KServer) Start() {
	entryLog := setupLog.WithName("entrypoint")

	// Setup a Manager
	entryLog.Info("setting up webhook server")
	mgr, err := manager.New(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:         scheme,
		LeaderElection: false,
		Port:           9443,
		CertDir:        "/etc/certs",
	})

	if err != nil {
		entryLog.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering mutating webhooks to the webhook server")
	hookServer.Register("/mutate-v1-pod", &webhook.Admission{Handler: &webhooks.PodInjector{
		Client: mgr.GetClient(),
		Logger: setupLog.WithName("PodInjector"),
		Handler: &webhooks.SidecarHandler{
			AgentImage: "1445277435/gomonitor-agent:v0.0.1",
			Url:        server.url,
			Bucket:     server.bucket,
			Org:        server.org,
			Token:      server.token,
		}}})

	entryLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

func NewKServer(flags *server.ServerFlags) *KServer {
	return &KServer{
		url:    flags.DBUrl,
		bucket: flags.Bucket,
		org:    flags.Organization,
		token:  flags.Token,
	}
}
