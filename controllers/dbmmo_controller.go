/*
Copyright 2020 HubertStefanski.

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
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/models"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1"
)

// DBMMOReconciler reconciles a DBMMO object
type DBMMOReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmoes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DBMMO object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *DBMMOReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.Background()
	log := r.Log.WithValues("dbmmo", req.NamespacedName)
	//Fetch dbmmo instance
	dbmmo := &v1.DBMMO{}
	if err := r.Get(ctx, req.NamespacedName, dbmmo); err != nil {
		if errors.IsNotFound(err) {
			//Object not found, return and don't requeue
			log.Info("DBMMO not found, ignoring")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get DBMMO")
		return ctrl.Result{}, err
	}
	// Check if the deployment already exists
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: dbmmo.Name, Namespace: dbmmo.Namespace}, found)
	// If the deployment doesn't exist, create a new one
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.getDeployment(dbmmo)
		if err != nil {
			log.Error(err, "Failed to get deployment")
		}
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBMMOReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.DBMMO{}).
		Complete(r)
}

func (r *DBMMOReconciler) getDeployment(d *v1.DBMMO) (*appsv1.Deployment, error) {
	labels := getLabels(d.Name)
	replicas := d.Spec.Size

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  models.OperatorName,
						Image: models.OperatorImage + ":" + models.OperatorVersion,
					},
					},
				},
			},
		},
	}
	// set the DBMMO instance as the owner and controller
	if err := ctrl.SetControllerReference(d, deployment, r.Scheme); err != nil {
		return nil, err
	}
	return deployment, nil
}

func getLabels(name string) map[string]string {
	return map[string]string{"app": "dbmmo", "dbmmo_cr": name}

}

func (r *DBMMOReconciler) getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
