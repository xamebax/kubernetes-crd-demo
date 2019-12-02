# kubernetes-crd-demo

This code accompanies a talk I gave at [Go Oslo User Group](https://www.meetup.com/Go-Oslo-User-Group/) in December 2019 on extending Kubernetes and code generation.

The Application CRD used in this demo is a simplified version of what we use in the [FIAAS go client](https://github.com/fiaas/fiaas-go-client).

## What's in this repository

- `pkg/apis/gooslo.io/v1/types.go`: annotated human-made code that defines the data structure of your custom resource.
- `pkg/apis/gooslo.io/v1/[doc.go|register.go]`: necessary boilerplate code.
- `pkg/apis/gooslo.io/v1/zz_generated.deepcopy.go`: generated deepcopy functions.
- `pkg/client` generated clientset, informers, and listers code to handle that custom resource.
- `cmd/…/main.go`: a simple client that watches events on our CRD and prints them out.

## Usage

### Prerequisites

- This code was written in Go 1.12. Might not work with Go 1.13 due to library incompatibilities. Haven't had time to test/fix this.
- You need a Kubernetes cluster. You can set one on your machine locally using either [kind](https://github.com/kubernetes-sigs/kind) or [minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/).

### Usage

1. Clone this repository and run `go get` to fetch dependencies.
2. Start your local Kubernetes cluster.
3. Run the executable. The easiest way is to run `KUBECONFIG=<path to your local cluster's kubeconfig> make local`

## Troubleshooting

If the event channel doesn't receive any new events (= there's no related activity in the cluster), it will time out and exit quietly. This is ok for the demo, but for production code you probably want another sort of queue. You can use generated informers for that (like in the somewhat confusing [sample controller](https://github.com/kubernetes/sample-controller)), or you can explore [`client-go`'s cache package](https://github.com/kubernetes/client-go/tree/master/tools/cache).

If you get an error that goes like `not enough arguments in call to watch.NewStreamWatcher`, it means you're using an incompatible `client-go` version in `go.mod`. There is a number of solutions to this problem [in this apimachinery github issue](https://github.com/kubernetes/apimachinery/issues/63).

## Further reading
