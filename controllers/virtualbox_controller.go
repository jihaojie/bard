/*
Copyright 2023 jihaojie.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1beta1 "github.com/jihaojie/bard/api/v1"
)

// VirtualboxReconciler reconciles a Virtualbox object
type VirtualboxReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=virtualbox.tiantian.cn,resources=virtualboxes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=virtualbox.tiantian.cn,resources=virtualboxes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=virtualbox.tiantian.cn,resources=virtualboxes/finalizers,verbs=update
//+kubebuilder:rbac:groups="app",resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Virtualbox object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *VirtualboxReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	rLog := log.FromContext(ctx)
	rLog.Info("VirtualBox start Reconciling.")

	var virtualBox v1beta1.Virtualbox

	if err := r.Client.Get(ctx, req.NamespacedName, &virtualBox); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//如果标记了删除 就不做处理.
	if virtualBox.DeletionTimestamp != nil {
		return ctrl.Result{}, nil
	}

	handle := Handler{Client: r.Client}

	stsChange, err := handle.CreateOrUpdateStatefulSet(ctx, req, &virtualBox)
	if err != nil {
		return ctrl.Result{}, err
	}
	rLog.Info("CreateOrUpdate StatefulSet Finished.", "change", stsChange)

	scChange, err := handle.CreateStorageClass(ctx, req, &virtualBox)
	if err != nil {
		return ctrl.Result{}, err
	}
	rLog.Info("CreateIfNotExist StorageClass Finished.", "change", scChange)

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *VirtualboxReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.Virtualbox{}).
		Complete(r)
}
