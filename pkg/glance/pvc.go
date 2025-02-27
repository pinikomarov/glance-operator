package glance

import (
	"fmt"
	glancev1 "github.com/openstack-k8s-operators/glance-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPvc - creates and returns a PVC object for a backing store
func GetPvc(api *glancev1.GlanceAPI, labels map[string]string, pvcType PvcType) (corev1.PersistentVolumeClaim, error) {
	// By default we point to the local storage pvc request
	// that will be customized in case the pvc is requested
	// for cache purposes
	var err error
	requestSize := api.Spec.StorageRequest
	pvcName := ServiceName
	if pvcType == "cache" {
		requestSize = api.Spec.ImageCacheSize
		// append -cache to avoid confusion when listing PVCs
		pvcName = fmt.Sprintf("%s-cache", ServiceName)
	}

	// Build the basic pvc object
	pvc := corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcName,
			Namespace: api.Namespace,
			Labels:    labels,
		},
	}

	// If the StorageRequest is a wrong string, we must return
	// an error. MustParse can't be used in this context as it
	// generates panic() and we can't recover the operator.
	storageSize, err := resource.ParseQuantity(requestSize)
	if err != nil {
		return pvc, err
	}

	pvc.Spec = corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: storageSize,
			},
		},
		StorageClassName: &api.Spec.StorageClass,
	}

	return pvc, err
}
