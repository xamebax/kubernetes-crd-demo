package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	v1 "github.com/xamebax/kubernetes-crd-demo/pkg/apis/gooslo.io/v1"
	clientset "github.com/xamebax/kubernetes-crd-demo/pkg/client/clientset/versioned"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ApplicationClient holds information on Kubernetes clients. It's needed to establish
// a client connection and watch resources.
type ApplicationClient struct {
	// Clientset is necessary if we want to operate on built-in Kubernetes resources.
	// APIExtensionsClientset is necessary to register a new CRD kind.
	// ApplicationClientset is a clientset used to make operations specific to this CRD's API Group.
	Clientset              kubernetes.Interface
	APIExtensionsClientset apiextensionsclientset.Interface
	ApplicationClientset   clientset.Interface
}

var kubeconfig string

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
}

func main() {
	var err error

	cfg, err := getClusterConfig()
	if err != nil {
		// We can't do anything without configuration
		log.Panicf("error building kubeconfig: %v", err)
	}

	applicationClient := ApplicationClient{
		Clientset:              createGenericClientset(cfg),
		APIExtensionsClientset: createAPIExtensionsClientset(cfg),
		ApplicationClientset:   createApplicationClientset(cfg),
	}

	// Initialize our CRD kind.
	_, err = v1.CreateCRD(applicationClient.APIExtensionsClientset)
	if err != nil {
		log.Panic(err)
	}

	// We instantiate a watcher using generated code.
	// Note: "default" is the namespace.
	applicationWatcher, err := applicationClient.
		ApplicationClientset.
		GoosloV1().
		Applications("default").
		Watch(metav1.ListOptions{})
	if err != nil {
		log.Panicf("cannot create watcher: %v", err)
	}

	// This infinite loop is what we need to go through events related to the
	// watched CRD (Application). The applicationWatcher has a result channel which
	// receives all events. If an error occurs, this channel will be closed. If there
	// is no new activity, the channel will eventually time out and close, too.
	for {
		event := <-applicationWatcher.ResultChan()
		// We break if event.Object is nil because otherwise we face a panic:
		// `interface conversion: runtime.Object is nil, not *v1.Application`.
		// break means the program eventually exit when there are no new events.
		// Using `continue` won't work here because the watch connection times out
		// quietly, and new events won't be registered or handled.
		if event.Object == nil {
			break
		}
		// Here, we convert the Object interface to the v1.Application interface
		// to be able to work with it further.
		var application *v1.Application = event.Object.(*v1.Application)

		// This is where we handle each object. In this demo, we just print the event,
		// but there's plenty of possibilities.
		log.Printf("Received event [%v] for application [%s]. It's using the [%v] image",
			event.Type, application.Name, application.Spec.Image)
	}
}

// createApplicationClientset returns a clientset specific to our custom resources.
func createApplicationClientset(cfg *rest.Config) *clientset.Clientset {
	clientset, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Panicf("unable to create new clientset: %s", err)
	}

	return clientset
}

// createAPIExtensionsClientset returns a clientset that allows registering new
// custom resources.
func createAPIExtensionsClientset(cfg *rest.Config) *apiextensionsclientset.Clientset {
	clientset, err := apiextensionsclientset.NewForConfig(cfg)
	if err != nil {
		log.Panicf("unable to create new clientset: %s", err)
	}

	return clientset
}

// CreateGenericClientset return a clientset that allows performing operations
// on built-in Kubernetes resources.
func createGenericClientset(cfg *rest.Config) *kubernetes.Clientset {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Panicf("unable to create new clientset: %s", err)
	}

	return clientset
}

// getClusterConfig is a helper function that recognizes if you're inside the
// cluster or not.
func getClusterConfig() (*rest.Config, error) {
	if kubeconfig == "" {
		log.Info("using in-cluster configuration")
		return rest.InClusterConfig()
	}
	log.Infof("using configuration from '%s'", kubeconfig)
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}
