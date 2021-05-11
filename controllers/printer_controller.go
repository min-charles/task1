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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;

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

			//지워지면 여기
			log.Info("Printer resource not found. Ignoring since object must be deleted")

			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Printer")
		return ctrl.Result{}, err
	}

	log.Info("SPEC NAME !!!!!! " + Printer.Spec.Name)
	log.Info("SPEC SIZE !!!!!! " + strconv.FormatInt(int64(Printer.Spec.Size), 10))

	// Check if the deployment already exists, if not create a new one
	return ctrl.Result{}, nil
}

// deploymentForPrinter returns a Printer Deployment object
// func (r *PrinterReconciler) deploymentForPrinter(m *taskv1.Printer) *appsv1.Deployment {
// 	ls := labelsForPrinter(m.Name)
// 	replicas := m.Spec.Size

// 	dep := &appsv1.Deployment{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      m.Name,
// 			Namespace: m.Namespace,
// 		},
// 		Spec: appsv1.DeploymentSpec{
// 			Replicas: &replicas,
// 			Selector: &metav1.LabelSelector{
// 				MatchLabels: ls,
// 			},
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: ls,
// 				},
// 				Spec: corev1.PodSpec{
// 					Containers: []corev1.Container{{
// 						Image:   "Printer:1.4.36-alpine",
// 						Name:    "Printer",
// 						Command: []string{"Printer", "-m=64", "-o", "modern", "-v"},
// 						Ports: []corev1.ContainerPort{{
// 							ContainerPort: 11211,
// 							Name:          "Printer",
// 						}},
// 					}},
// 				},
// 			},
// 		},
// 	}
// 	// Set Printer instance as the owner and controller
// 	ctrl.SetControllerReference(m, dep, r.Scheme)
// 	return dep
// }

// labelsForPrinter returns the labels for selecting the resources
// belonging to the given Printer CR name.
// func labelsForPrinter(name string) map[string]string {
// 	return map[string]string{"app": "Printer", "Printer_cr": name}
// }

// // getPodNames returns the pod names of the array of pods passed in
// func getPodNames(pods []corev1.Pod) []string {
// 	var podNames []string
// 	for _, pod := range pods {
// 		podNames = append(podNames, pod.Name)
// 	}
// 	return podNames
// }

func (r *PrinterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&taskv1.Printer{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
