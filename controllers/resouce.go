package controllers

import (
	v1beta1 "github.com/jihaojie/bard/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func MakeReference(virtualbox *v1beta1.Virtualbox) []metav1.OwnerReference {
	return []metav1.OwnerReference{
		*metav1.NewControllerRef(virtualbox, schema.GroupVersionKind{
			Group:   v1beta1.GroupVersion.Group,
			Version: v1beta1.GroupVersion.Version,
			Kind:    v1beta1.Kind,
		}),
	}
}

func NewStatefulSet(virtualbox *v1beta1.Virtualbox) *appsv1.StatefulSet {
	labels := map[string]string{"virtualbox": virtualbox.Name}
	selector := &metav1.LabelSelector{
		MatchLabels: labels,
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:            virtualbox.Name,
			Namespace:       virtualbox.Namespace,
			OwnerReferences: MakeReference(virtualbox),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: virtualbox.Spec.Size,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: virtualbox.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: NewContainers(virtualbox),
				},
			},
		},
	}
}

func NewContainers(virtualbox *v1beta1.Virtualbox) []corev1.Container {
	return []corev1.Container{
		{
			Name:  virtualbox.Name,
			Image: "busybox:latest",
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("10Mi"),
					corev1.ResourceCPU:    resource.MustParse("1"),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse("10Mi"),
					corev1.ResourceCPU:    resource.MustParse("1"),
				},
			},
			Command: []string{"/bin/sh", "-c", "echo hello word && tail -f /dev/null"},
		},
	}
}

func NewStorageClass(virtualbox *v1beta1.Virtualbox) storagev1.StorageClass {
	return storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs",
			Namespace: virtualbox.Namespace,
		},
		Provisioner: "kubernetes.io/nfs",
		Parameters: map[string]string{
			"server": "10.1.32.159",
			"path":   "/data/nfs",
		},
	}
}

//TODO:这个先放放 最后再加
func NewPVC(virtualbox *v1beta1.Virtualbox) corev1.PersistentVolumeClaim {
	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      virtualbox.Name,
			Namespace: virtualbox.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{},
	}
}
