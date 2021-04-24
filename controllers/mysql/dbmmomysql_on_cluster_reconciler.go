package mysql

import (
	"context"
	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/model"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *DBMMOMySQLReconciler) onClusterReconcileMysqlStatus(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL, listOpts []client.ListOption) (ctrl.Result, error) {
	// Update the mysql status with the pod names
	// List the pods for this mysql's deployment
	r.Log.Info("Reconciling Mysql Status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
	podList := &corev1.PodList{}
	if err := r.Client.List(ctx, podList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list pods", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	podNames := model.GetPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, mysql.Status.Nodes) {
		mysql.Status.Nodes = podNames
		err := r.Client.Status().Update(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
			return ctrl.Result{}, err
		}
	}

	// Update the mysql status with the service names
	// List the services for this mysql's deployment
	serviceList := &corev1.ServiceList{}
	if err := r.Client.List(ctx, serviceList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list services", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	serviceNames := model.GetServiceNames(serviceList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(serviceNames, mysql.Status.Services) {
		mysql.Status.Services = serviceNames
		err := r.Client.Status().Update(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
			return ctrl.Result{}, err
		}
	}
	// Update the mysql status with the PersistentVolumeClaim names
	// List the PersistentVolumeClaims for this mysql's deployment
	pvcList := &corev1.PersistentVolumeClaimList{}
	if err := r.Client.List(ctx, pvcList, listOpts...); err != nil {
		r.Log.Error(err, "Failed to list PersistentVolumeClaim", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}
	pvcNames := model.GetPvcNames(pvcList.Items)

	// Update status.PersistentVolume if needed
	if !reflect.DeepEqual(pvcNames, mysql.Status.PersistentVolumeClaims) {
		mysql.Status.PersistentVolumeClaims = pvcNames
		err := r.Client.Status().Update(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Failed to update Mysql status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
			return ctrl.Result{}, err
		}
	}
	r.Log.Info("Reconciled Mysql status ", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
	return ctrl.Result{}, nil
}

func (r *DBMMOMySQLReconciler) onClusterReconcileMysqlDeployment(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	// Check if the deployment already exists, if not create a new one

	replicas := mysql.Spec.Size
	// Define a new deployment
	dep := model.GetMysqlDeployment(mysql)
	// Set Mysql instance as the owner and controller
	_ = ctrl.SetControllerReference(mysql, dep, r.Scheme)

	r.Log.Info("Reconciling deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, dep, func() error {
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
						EnvFrom:         model.MysqlDeploymentGetEnvFrom(mysql),
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

	r.Log.Info("Deployment reconciled", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) onClusterReconcileMysqlService(ctx context.Context, m *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	// Define a new service
	service := model.GetMysqlService(m)

	_ = ctrl.SetControllerReference(m, service, r.Scheme)

	r.Log.Info("Reconciling service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:       constants.MysqlContainerPortName,
				Port:       constants.MysqlContainerPort,
				Protocol:   corev1.ProtocolTCP,
				TargetPort: intstr.FromString(constants.MysqlContainerPortName),
			},
		}
		service.Spec.ClusterIP = corev1.ClusterIPNone
		return nil
	})

	if err != nil {
		r.Log.Error(err, "Failed to reconcile Service", "Service.Namespace", service.Namespace, "service.Name", service.Name)
		return ctrl.Result{}, err
	}

	r.Log.Info("Service reconciled", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
	return ctrl.Result{Requeue: true}, nil

}

func (r *DBMMOMySQLReconciler) onClusterReconcileMysqlPVC(ctx context.Context, m *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	foundPVC := &corev1.PersistentVolumeClaim{}
	if err := r.Client.Get(ctx, types.NamespacedName{Name: constants.MysqlClaimName, Namespace: m.Namespace}, foundPVC); err != nil && k8serr.IsNotFound(err) {
		// Define a new PersistentVolume
		pvc := model.GetMysqlPvc(m)
		r.Log.Info("Reconciling PVC", "Pvc.Namespace", pvc.Namespace, "Pvc.Name", pvc.Name)
		_ = ctrl.SetControllerReference(m, pvc, r.Scheme)
		r.Log.Info("Creating a new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
		if err = r.Client.Create(ctx, pvc); err != nil {
			r.Log.Error(err, "Failed to create new PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)
			return ctrl.Result{}, err
		}

		r.Log.Info("PersistentVolumeClaim created", "PersistentVolumeClaim.Namespace", pvc.Namespace, "PersistentVolumeClaim.Name", pvc.Name)

		// PrivateVolume created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get PersistentVolumeClaim", "PersistentVolumeClaim.Namespace", foundPVC.Namespace, "PersistentVolumeClaim.Name", foundPVC.Name)
		return ctrl.Result{}, err
	}
	r.Log.Info("PVC reconciled", "Pvc.Namespace", foundPVC.Namespace, "Pvc.Name", foundPVC.Name)
	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) onClusterReconcileIngress(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	ingr := model.GetMysqlIngress(mysql)

	_ = ctrl.SetControllerReference(mysql, ingr, r.Scheme)

	r.Log.Info("Reconciling ingress", "Ingress.Namespace", ingr.Namespace, "Ingress.Name", ingr.Name)
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, ingr, func() error {
		specific := netv1.PathTypePrefix
		ingr.Spec = netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: constants.MysqlHostName,
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{
									Path:     constants.MysqlDefaultPath,
									PathType: &specific,
									Backend: netv1.IngressBackend{
										Service: &netv1.IngressServiceBackend{
											Name: constants.MysqlServiceName,
											Port: netv1.ServiceBackendPort{
												Name: constants.MysqlServiceName,
											},
										},
										Resource: nil,
									},
								},
							},
						},
					},
				},
			},
		}
		return nil
	})
	if err != nil {
		r.Log.Error(err, "Failed to reconcile Ingress", "Ingress.Name", ingr.Namespace, "Ingress.Name", ingr.Name)
		return ctrl.Result{}, err
	}

	r.Log.Info("Deployment reconciled", "Ingress.Namespace", ingr.Namespace, "Ingress.Name", ingr.Name)
	return ctrl.Result{Requeue: true}, nil
}

//OnClusterCleanup cleans up the resources for a specific Mysql object
func (r *DBMMOMySQLReconciler) OnClusterCleanup(ctx context.Context, m *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	pvc := model.GetMysqlPvc(m)
	if err := r.Client.Delete(ctx, pvc); err != nil && k8serr.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
	}
	svc := model.GetMysqlService(m)
	if err := r.Client.Delete(ctx, svc); err != nil && k8serr.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
	}
	dep := model.GetMysqlDeployment(m)
	if err := r.Client.Delete(ctx, dep); err != nil && k8serr.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
	}

	// Don't requeue if the cleanup was successful
	return ctrl.Result{}, nil
}

func (r *DBMMOMySQLReconciler) cleanUpIngress(ctx context.Context, m *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	ingr := model.GetMysqlIngress(m)
	if err := r.Client.Delete(ctx, ingr); err != nil && k8serr.IsNotFound(err) {
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
	}
	return ctrl.Result{}, nil
}
