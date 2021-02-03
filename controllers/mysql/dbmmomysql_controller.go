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

package mysql

import (
	"context"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/mysql/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
)

// DBMMOMySQLReconciler reconciles a DBMMOMySQL object
type DBMMOMySQLReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DBMMOMySQL object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *DBMMOMySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.Background()
	log := r.Log.WithValues(constants.MysqlControllerName, req.NamespacedName)

	// Fetch the Memcached instance
	mysql := &cachev1alpha1.DBMMOMySQL{}
	result := ctrl.Result{}
	err := r.Get(ctx, req.NamespacedName, mysql)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not foundDeployment, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("dbmmomysql resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get dbmmomysql")
		return ctrl.Result{}, err
	}

	// Create list options
	listOpts := []client.ListOption{
		client.InNamespace(mysql.Namespace),
		client.MatchingLabels(model.GetLabels(mysql.Name)),
	}

	result, err = r.reconcileMysqlPVC(ctx, mysql, listOpts)
	if err != nil {
		return result, err
	}

	result, err = r.reconcileMysqlService(ctx, mysql, listOpts)
	if err != nil {
		return result, err
	}

	result, err = r.reconcileMysqlDeployment(ctx, mysql, listOpts)
	if err != nil {
		return result, err
	}

	return ctrl.Result{}, nil
}

func (r *DBMMOMySQLReconciler) reconcileMysqlDeployment(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	// Check if the deployment already exists, if not create a new one
	foundDeployment := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{Name: constants.MysqlDeploymentName, Namespace: mysql.Namespace}, foundDeployment)

	if err != nil && errors.IsNotFound(err) {

		// Define a new deployment
		dep := model.GetMysqlDeployment(mysql)
		// Set Mysql instance as the owner and controller
		_ = ctrl.SetControllerReference(mysql, dep, r.Scheme)

		r.Log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

		err = r.Create(ctx, dep)
		if err != nil {
			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}

		r.Log.Info("Deployment created", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)

		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := mysql.Spec.Size
	if *foundDeployment.Spec.Replicas != size {
		foundDeployment.Spec.Replicas = &size
		err = r.Update(ctx, foundDeployment)
		if err != nil {
			r.Log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the mysql status with the pod names
	// List the pods for this mysql's deployment
	podList := &corev1.PodList{}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list pods", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	podNames := model.GetPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, mysql.Status.Nodes) {
		mysql.Status.Nodes = podNames
		err := r.Status().Update(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) reconcileMysqlService(ctx context.Context, m *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	// Check if the service already exists, if not create a new one
	foundService := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: constants.MysqlServiceName, Namespace: m.Namespace}, foundService); err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := model.GetMysqlService(m)

		_ = ctrl.SetControllerReference(m, service, r.Scheme)

		r.Log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		if err := r.Create(ctx, service); err != nil {
			r.Log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return ctrl.Result{}, err
		}
		r.Log.Info("Service created", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Update the mysql status with the service names
	// List the services for this mysql's deployment
	serviceList := &corev1.ServiceList{}
	if err := r.List(ctx, serviceList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list services", "Mysql.Namespace", m.Namespace, "Mysql.Name", m.Name)
		return ctrl.Result{}, err
	}
	serviceNames := model.GetServiceNames(serviceList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(serviceNames, m.Status.Services) {
		m.Status.Services = serviceNames
		err := r.Status().Update(ctx, m)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{Requeue: true}, nil

}

func (r *DBMMOMySQLReconciler) reconcileMysqlPVC(ctx context.Context, m *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, types.NamespacedName{Name: constants.MysqlClaimName, Namespace: m.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		// Define a new PersistentVolume
		pvc := model.GetMysqlPvc(m)

		_ = ctrl.SetControllerReference(m, pvc, r.Scheme)
		r.Log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		if err = r.Create(ctx, pvc); err != nil {
			r.Log.Error(err, "Failed to create new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			return ctrl.Result{}, err
		}

		r.Log.Info("PersistentVolumeClaim created", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)

		// PrivateVolume created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get PersistentVolumeClaim")
		return ctrl.Result{}, err
	}

	// Update the mysql status with the PersistentVolumeClaim names
	// List the PersistentVolumeClaims for this mysql's deployment
	pvcList := &corev1.PersistentVolumeClaimList{}
	if err = r.List(ctx, pvcList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list PersistentVolumeClaim", "Mysql.Namespace", m.Namespace, "Mysql.Name", m.Name)
		return ctrl.Result{}, err
	}
	pvcNames := model.GetPvcNames(pvcList.Items)

	// Update status.PersistentVolume if needed
	if !reflect.DeepEqual(pvcNames, m.Status.PersistentVolumeClaims) {
		m.Status.PersistentVolumeClaims = pvcNames
		err := r.Status().Update(ctx, m)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBMMOMySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.DBMMOMySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
