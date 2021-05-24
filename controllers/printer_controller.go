/*
Copyright 2021.

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
	"strconv"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	taskv1 "github.com/min-charles/task1.git/api/v1"
)

// PrinterReconciler reconciles a Printer object
type PrinterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=task.my.domain,resources=printers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=task.my.domain,resources=printers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=task.my.domain,resources=printers/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;create

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Printer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PrinterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Printer", req.NamespacedName)

	// Fetch the Printer instance
	Printer := &taskv1.Printer{}

	err := r.Get(ctx, req.NamespacedName, Printer)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Printer resource not found. Ignoring since object must be deleted")

			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Printer")
		return ctrl.Result{}, err
	}

	log.Info("SPEC NAME !!!!!! " + Printer.Spec.Name)
	log.Info("SPEC SIZE !!!!!! " + strconv.FormatInt(int64(Printer.Spec.Size), 10))

	// check if printing pod is already existing//
	found := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Name: "log-printer-pod", Namespace: "task1-system"}, found)
	log.Info("Already existing pod: " + found.ObjectMeta.Name)

	if err != nil && errors.IsNotFound(err) {
		// declare pod which prints CR spec//
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "task1-system",
				Name:      "log-printer-pod",
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Image:           "lmch0000/log-printer:latest", // "PrinterName" EnvVar print image
						Name:            "log-printer-ctn",
						ImagePullPolicy: corev1.PullAlways,
						Env: []corev1.EnvVar{ // declare CR spec printed
							{
								Name:  "PrinterName",
								Value: Printer.Spec.Name,
							},
						},
					},
				},
			},
		}

		//Set OwnerReference//
		if err2 := controllerutil.SetOwnerReference(Printer, pod, r.Scheme); err2 != nil {
			return ctrl.Result{}, err2
		}

		/////Pod 생성////////
		err = r.Client.Create(context.Background(), pod)

		if err != nil {
			// Error reading the object - requeue the request.
			log.Error(err, "Fail to create pod")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil

	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *PrinterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1.Printer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
