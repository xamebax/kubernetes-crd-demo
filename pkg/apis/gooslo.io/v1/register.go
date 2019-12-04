// Most of this is boilerplate code. The important parts are:
// - constants that define naming
// - addKnownTypes() function.

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName contains the api name
// GroupVersion contains the api version
const (
	GroupName                    = "gooslo.io"
	Kind                         = "application"
	GroupVersion                 = "v1"
	SingularName                 = "application"
	PluralName                   = "applications"
	CustomResourceDefinitionName = PluralName + "." + GroupName
)

var SchemeGroupVersion = schema.GroupVersion{
	Group:   GroupName,
	Version: GroupVersion,
}

// localSchemeBuilder will stay in k8s.io/kubernetes.
// AddToScheme will stay in k8s.io/kubernetes.
var (
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
	SchemeBuilder      runtime.SchemeBuilder
)

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to api.Scheme
// so it can be used with the go-client library.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Application{},
		&ApplicationList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&metav1.Status{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
