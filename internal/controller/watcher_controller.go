/*
Copyright 2024 Malek Zaag.

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

package controller

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/Malek-Zaag/MyNewOperator/api/v1beta1"
	watchersv1beta1 "github.com/Malek-Zaag/MyNewOperator/api/v1beta1"
)

// WatcherReconciler reconciles a Watcher object
type WatcherReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=watchers.malek.dev,resources=watchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=watchers.malek.dev,resources=watchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=watchers.malek.dev,resources=watchers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Watcher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *WatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log.Log.Info("Reconcile loop here")
	watcher := &v1beta1.Watcher{}
	err := r.Get(ctx, req.NamespacedName, watcher)
	if err != nil {
		return ctrl.Result{}, err
	}
	startTime := watcher.Spec.Start
	endTime := watcher.Spec.End
	replicas := watcher.Spec.Replicas
	currentHour := time.Now().UTC().Hour()
	log.Log.Info(fmt.Sprintf("The time now is %d", currentHour))
	log.Log.Info(fmt.Sprintf("The start time is %d", startTime))
	// TODO(user): your logic here
	for _, deploy := range watcher.Spec.Deployments {
		log.Log.Info("here")
		if currentHour >= startTime && currentHour <= endTime {
			err := ScaleDeployment(ctx, r, deploy, replicas)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}
	return ctrl.Result{RequeueAfter: 30}, nil
}

func ScaleDeployment(ctx context.Context, r *WatcherReconciler, deploy watchersv1beta1.NamespacedName, replicas int32) error {
	log.Log.Info("Scaling up replicas")
	deployment := &v1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: deploy.Namespace,
		Name:      deploy.Name,
	}, deployment)
	if err != nil {
		return err
	}
	if deployment.Spec.Replicas != &replicas {
		deployment.Spec.Replicas = &replicas
		err := r.Update(ctx, deployment)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&watchersv1beta1.Watcher{}).
		Complete(r)
}
