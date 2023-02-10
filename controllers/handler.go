package controllers

import (
	"context"
	v1beta1 "github.com/jihaojie/bard/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type Handler struct {
	Client client.Client
}

//判断是否被更新过
func (h *Handler) IsUpdated(ctx context.Context, virtualBox v1beta1.Virtualbox) bool {

	return true
}

func (h *Handler) CreateOrUpdateStatefulSet(ctx context.Context, req ctrl.Request, virtualbox *v1beta1.Virtualbox) (controllerutil.OperationResult, error) {
	newSts := NewStatefulSet(virtualbox)

	loadSts := &appsv1.StatefulSet{}
	loadSts.Name = newSts.Name
	loadSts.Namespace = newSts.Namespace

	change, err := controllerutil.CreateOrUpdate(ctx, h.Client, loadSts, func() error {
		deepCopy := newSts.Spec.DeepCopy()
		if deepCopy != nil {
			loadSts.Spec = *deepCopy
		}

		loadSts.Annotations = newSts.Annotations
		loadSts.ObjectMeta.OwnerReferences = newSts.ObjectMeta.OwnerReferences

		return nil
	})

	if err != nil {
		return controllerutil.OperationResultNone, err
	}

	return change, nil
}

func (h *Handler) CreateStorageClass(ctx context.Context, req ctrl.Request, virtualbox *v1beta1.Virtualbox) (controllerutil.OperationResult, error) {
	sc := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nfs",
			Namespace: virtualbox.Namespace,
		},
	}

	if err := h.Client.Get(ctx, req.NamespacedName, &sc); err != nil {
		newSc := NewStorageClass(virtualbox)
		if errors.IsNotFound(err) {
			if err := h.Client.Create(ctx, &newSc); err != nil {
				return controllerutil.OperationResultNone, err
			}
		}
		return controllerutil.OperationResultCreated, nil
	}
	return controllerutil.OperationResultNone, nil
}

func (h *Handler) CreateOrUpdatePvc() {

}
