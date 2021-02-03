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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
)

// DBMMOMySQLReconciler reconciles a DBMMOMySQL object
type DBMMOMySQLReconciler struct {
	client client.Client
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
	err := r.client.Get(ctx, req.NamespacedName, mysql)
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

	result, err = r.reconcileMysqlDeployment(ctx, mysql)
	if err != nil {
		return result, err
	}

	return ctrl.Result{}, nil
}


func (r *DBMMOMySQLReconciler) reconcileMysqlStatus(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error)  {
	// Update the mysql status with the pod names
	// List the pods for this mysql's deployment
	podList := &corev1.PodList{}
	if err := r.client.List(ctx, podList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list pods", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	podNames := model.GetPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, mysql.Status.Nodes) {
		mysql.Status.Nodes = podNames
		err := r.client.Status().Update(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) reconcileMysqlDeployment(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	// Check if the deployment already exists, if not create a new one

	replicas := mysql.Spec.Size
	// Define a new deployment
	dep := model.GetMysqlDeployment(mysql)
	// Set Mysql instance as the owner and controller
	_ = ctrl.SetControllerReference(mysql, dep, r.Scheme)

	_, err := controllerutil.CreateOrUpdate(ctx, r.client, dep, func() error {
		r.Log.Info("Reconciling deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		dep.Spec = appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: model.GetLabels(mysql.Name),
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: constants.MysqlStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: model.GetLabels(mysql.Name),
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: constants.MysqlClaimName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: constants.MysqlClaimName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Name:  constants.MysqlContainerName,
						Image: constants.MysqlContainerImage,
						Ports: []corev1.ContainerPort{{
							ContainerPort: constants.MysqlContainerPort,
							Name:          constants.MysqlContainerPortName,
						}},
						EnvFrom: nil,
						Env: []corev1.EnvVar{
							{
								Name:  constants.MysqlSecretEnvName,
								Value: constants.MysqlSecretEnvVal,
								// TODO add Secret here
								//ValueFrom: nil,
							},
						},
						ImagePullPolicy: "IfNotPresent",
					},
					}},
			},
		}
		return nil
	})
	if err != nil {
		r.Log.Error(err, "Failed to reconcile Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) reconcileMysqlService(ctx context.Context, m *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	// Check if the service already exists, if not create a new one
	foundService := &corev1.Service{}
	if err := r.client.Get(ctx, types.NamespacedName{Name: constants.MysqlServiceName, Namespace: m.Namespace}, foundService); err != nil && errors.IsNotFound(err) {
		// Define a new service
		service := model.GetMysqlService(m)

		_ = ctrl.SetControllerReference(m, service, r.Scheme)

		r.Log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		if err := r.client.Create(ctx, service); err != nil {
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
	if err := r.client.List(ctx, serviceList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list services", "Mysql.Namespace", m.Namespace, "Mysql.Name", m.Name)
		return ctrl.Result{}, err
	}
	serviceNames := model.GetServiceNames(serviceList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(serviceNames, m.Status.Services) {
		m.Status.Services = serviceNames
		err := r.client.Status().Update(ctx, m)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{Requeue: true}, nil

}

func (r *DBMMOMySQLReconciler) reconcileMysqlPVC(ctx context.Context, m *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	foundPVC := &corev1.PersistentVolumeClaim{}
	err := r.client.Get(ctx, types.NamespacedName{Name: constants.MysqlClaimName, Namespace: m.Namespace}, foundPVC)
	if err != nil && errors.IsNotFound(err) {
		// Define a new PersistentVolume
		pvc := model.GetMysqlPvc(m)

		_ = ctrl.SetControllerReference(m, pvc, r.Scheme)
		r.Log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		if err = r.client.Create(ctx, pvc); err != nil {
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
	if err = r.client.List(ctx, pvcList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list PersistentVolumeClaim", "Mysql.Namespace", m.Namespace, "Mysql.Name", m.Name)
		return ctrl.Result{}, err
	}
	pvcNames := model.GetPvcNames(pvcList.Items)

	// Update status.PersistentVolume if needed
	if !reflect.DeepEqual(pvcNames, m.Status.PersistentVolumeClaims) {
		m.Status.PersistentVolumeClaims = pvcNames
		err := r.client.Status().Update(ctx, m)
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
