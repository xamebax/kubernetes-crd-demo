// Note: In the real world, you might want to have a yaml specification that
// you use when you bootstrap your Kubernetes cluster instead of programatically
// creating a CRD.

package v1

import (
	"reflect"
	"time"

	"github.com/prometheus/common/log"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// PollPeriod defines how often we poll to check if CRD is there
// PollTimeout defines after how much time we'll stop waiting to check if CRD was created
const (
	PollPeriod  = 5 * time.Second
	PollTimeout = 30 * time.Second
)

// CreateCRD registers a new CRD kind and then periodically polls the API server
// to check if the CRD was created.
func CreateCRD(clientset apiextensionsclientset.Interface) (*apiextensionsv1beta1.CustomResourceDefinition, error) {

	crdSchema := createCRDSchema(clientset)

	_, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crdSchema)
	switch {
	case err == nil:
		log.Infof("creating %s CRD…", crdSchema.Name)
	case apierrors.IsAlreadyExists(err):
		log.Infof("CRD %s already exists, not creating", crdSchema.Name)
	default:
		return nil, err
	}

	err = wait.Poll(PollPeriod, PollTimeout, func() (bool, error) {
		crd, err := clientset.ApiextensionsV1beta1().CustomResourceDefinitions().Get(CustomResourceDefinitionName, metav1.GetOptions{})
		if err != nil {
			log.Errorf("error when trying to wait for CRD creation: %s", err)
			return false, err
		}

		for _, condition := range crd.Status.Conditions {
			switch condition.Type {
			case apiextensionsv1beta1.Established:
				if condition.Status == apiextensionsv1beta1.ConditionTrue {
					return true, nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				if condition.Status == apiextensionsv1beta1.ConditionFalse {
					return false, nil
				}
			}
		}
		return false, err
	})

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return crdSchema, nil
}

func createCRDSchema(clientset apiextensionsclientset.Interface) *apiextensionsv1beta1.CustomResourceDefinition {
	crdSchema := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CustomResourceDefinitionName,
			Namespace: "default",
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   GroupName,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: PluralName,
				Kind:   reflect.TypeOf(Application{}).Name(),
			},
		},
	}
	return crdSchema
}
