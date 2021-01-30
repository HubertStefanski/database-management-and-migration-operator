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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Check if the deployment already exists, if not create a new one
	foundDeployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep := r.getMysqlDeployment(mysql)
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

	// Ensure the deployment size is the same as the spec
	size := mysql.Spec.Size
	if *foundDeployment.Spec.Replicas != size {
		foundDeployment.Spec.Replicas = &size
		err = r.Update(ctx, foundDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the mysql status with the pod names
	// List the pods for this mysql's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(mysql.Namespace),
		client.MatchingLabels(getLabels(mysql.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, mysql.Status.Nodes) {
		mysql.Status.Nodes = podNames
		err := r.Status().Update(ctx, mysql)
		if err != nil {
			log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}

	// Check if the service already exists, if not create a new one
	foundService := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		service := r.getMysqlService(mysql)
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	// Update the mysql status with the service names
	// List the services for this mysql's deployment
	serviceList := &corev1.ServiceList{}
	if err = r.List(ctx, serviceList, listOpts...); err != nil {
		log.Error(err, "Failed to list services", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	serviceNames := getServiceNames(serviceList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(serviceNames, mysql.Status.Services) {
		mysql.Status.Nodes = serviceNames
		err := r.Status().Update(ctx, mysql)
		if err != nil {
			log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}
	//------------------------------------------------------------------------------------------------------------------
	// Check if the service already exists, if not create a new one
	foundPV := &corev1.PersistentVolume{}
	err = r.Get(ctx, types.NamespacedName{Name: mysql.Name, Namespace: mysql.Namespace}, foundPV)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		pv := r.getMysqlService(mysql)
		log.Info("Creating a new PrivateVolume", "PrivateVolume.Namespace", pv.Namespace, "Service.Name", pv.Name)
		err = r.Create(ctx, pv)
		if err != nil {
			log.Error(err, "Failed to create new PrivateVolume", "PrivateVolume.Namespace", pv.Namespace, "PrivateVolume.Name", pv.Name)
			return ctrl.Result{}, err
		}
		// PrivateVolume created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get PrivateVolume")
		return ctrl.Result{}, err
	}

	// Update the mysql status with the service names
	// List the services for this mysql's deployment
	pvList := &corev1.PersistentVolumeList{}
	if err = r.List(ctx, pvList, listOpts...); err != nil {
		log.Error(err, "Failed to list privateVolumes", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	pvNames := getPvNames(pvList.Items)

	// Update status.PrivateVolumes if needed
	if !reflect.DeepEqual(pvNames, mysql.Status.Services) {
		mysql.Status.PrivateVolumes = pvNames
		err := r.Status().Update(ctx, mysql)
		if err != nil {
			log.Error(err, "Failed to update Mysql status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *DBMMOMySQLReconciler) getMysqlService(m *cachev1alpha1.DBMMOMySQL) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.MysqlServiceName,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: constants.MysqlContainerPort,
				},
			},
		},
	}
	// Set Mysql instance as the owner and controller
	_ = ctrl.SetControllerReference(m, service, r.Scheme)
	return service

}

// getMysqlDeployment returns a mysql Deployment object
func (r *DBMMOMySQLReconciler) getMysqlDeployment(m *cachev1alpha1.DBMMOMySQL) *appsv1.Deployment {
	ls := getLabels(m.Name)
	replicas := m.Spec.Size

	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: constants.MysqlStrategyType,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: constants.MysqlPVName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: constants.MysqlClaimName,
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: constants.MysqlContainerImage,
						Name:  constants.MysqlContainerName,
						Env: []corev1.EnvVar{
							{
								Name:  constants.MysqlSecretEnvName,
								Value: constants.MysqlSecretEnvVal,
								// TODO add Secret here
								//ValueFrom: nil,
							},
						},
						//Command: []string{"memcached", "-m=64", "-o", "modern", "-v"},
						Ports: []corev1.ContainerPort{{
							ContainerPort: constants.MysqlContainerPort,
							Name:          constants.MysqlContainerPortName,
						}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      constants.MysqlVolumeMountName,
								MountPath: constants.MysqlVolumeMountPath,
							},
						},
					}},
				},
			},
		},
	}
	// Set Mysql instance as the owner and controller
	_ = ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// getLabels returns the labels for selecting the resources
// belonging to the given mysql CR name.
func getLabels(name string) map[string]string {
	return map[string]string{"app": constants.MysqlAppSelector, "mysql_cr": "mysql"}
}

// getPvNames returns the pv names of mysql
func getPvNames(pvs []corev1.PersistentVolume) []string {
	var persistentVolumesNames []string
	for _, pv := range pvs {
		persistentVolumesNames = append(persistentVolumesNames, pv.Name)
	}
	return persistentVolumesNames
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// getServiceNames returns the pod names of the array of pods passed in
func getServiceNames(services []corev1.Service) []string {
	var serviceNames []string
	for _, pod := range services {
		serviceNames = append(serviceNames, pod.Name)
	}
	return serviceNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBMMOMySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.DBMMOMySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
