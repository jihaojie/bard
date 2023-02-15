package controllers

import (
	v1beta1 "github.com/jihaojie/bard/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strconv"
)

const (
	StorageClassName = "sc-nfs"
	ProvisionerName  = "fuseim.pri/ifs"

	NFSServer    = "10.1.32.159"
	NFSMountPath = "/data/nfs"
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
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				NewPVC(virtualbox),
			},
		},
	}
}

func NewContainers(virtualbox *v1beta1.Virtualbox) []corev1.Container {

	CPUStr := strconv.Itoa(virtualbox.Spec.CPU)
	return []corev1.Container{
		{
			Name:  virtualbox.Name,
			Image: "busybox:latest",
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse(virtualbox.Spec.DiskSize),
					corev1.ResourceCPU:    resource.MustParse(CPUStr),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceMemory: resource.MustParse(virtualbox.Spec.DiskSize),
					corev1.ResourceCPU:    resource.MustParse(CPUStr),
				},
			},
			Command: []string{"/bin/sh", "-c", "echo hello word && tail -f /dev/null"},
		},
	}
}

func NewStorageClass(virtualbox *v1beta1.Virtualbox) storagev1.StorageClass {
	return storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: StorageClassName,
		},
		Provisioner: ProvisionerName,
		Parameters: map[string]string{
			"server": NFSServer,
			"path":   NFSMountPath,
		},
	}
}

func NewPVC(virtualbox *v1beta1.Virtualbox) corev1.PersistentVolumeClaim {
	sc := NewStorageClass(virtualbox)

	return corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:            virtualbox.Name,
			Namespace:       virtualbox.Namespace,
			OwnerReferences: MakeReference(virtualbox), // 这考虑一下是否要做关联？
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("5Mi"), //  这改成参数
				},
			},
			StorageClassName: &sc.Name,
		},
	}
}
