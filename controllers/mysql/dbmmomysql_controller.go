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
	"fmt"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/model"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
)

// DBMMOMySQLReconciler reconciles a DBMMOMySQL object
type DBMMOMySQLReconciler struct {
	Client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.my.domain,resources=dbmmomysqls/finalizers,verbs=update
func (r *DBMMOMySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx = context.Background()
	log := r.Log.WithValues(constants.MysqlControllerName, req.NamespacedName)
	result := ctrl.Result{}

	// Fetch the Memcached instance
	mysql := &cachev1alpha1.DBMMOMySQL{}
	err := r.Client.Get(ctx, req.NamespacedName, mysql)
	if err != nil {
		if k8serr.IsNotFound(err) {
			// Request object not foundDeployment, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("dbmmomysql resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get dbmmomysql", "Request.namespacedName", req.NamespacedName)
		return ctrl.Result{}, err
	}

	// Create list options
	listOpts := []client.ListOption{
		client.InNamespace(mysql.Namespace),
		client.MatchingLabels(model.GetLabels(mysql.Name)),
	}

	if mysql.Spec.Deployment != nil && mysql.Spec.Deployment.DeploymentType != nil && *mysql.Spec.Deployment.DeploymentType != "" {
		switch depType := *mysql.Spec.Deployment.DeploymentType; depType {
		case constants.MysqlDeploymentTypeOnCluster:
			if result, err = r.onClusterReconcileMysqlPVC(ctx, mysql); err != nil {
				return result, err
			}

			if result, err = r.onClusterReconcileMysqlService(ctx, mysql); err != nil {
				return result, err
			}

			// Only create ingress if directly specified to do so
			if mysql.Spec.Deployment.Ingress != nil && mysql.Spec.Deployment.Ingress.Enabled != nil && *mysql.Spec.Deployment.Ingress.Enabled != false {
				if result, err = r.onClusterReconcileIngress(ctx, mysql); err != nil {
					return result, err
				}
			}

			if mysql.Spec.Deployment.Ingress != nil && mysql.Spec.Deployment.Ingress.Enabled != nil && *mysql.Spec.Deployment.Ingress.Enabled != true { // If an ingress exists but is not enabled, delete it
				if !k8serr.IsNotFound(r.Client.Get(ctx,
					types.NamespacedName{
						Namespace: mysql.Namespace,
						Name:      model.GetMysqlIngress(mysql).Name,
					}, mysql)) {
					result, err = r.cleanUpIngress(ctx, mysql)
					if err != nil {
						return result, err
					}
				}
			}

			if result, err = r.onClusterReconcileMysqlDeployment(ctx, mysql); err != nil {
				return result, err
			}

			// wait for resources to be ready before updating status
			ready, err := r.getCollectiveReadiness(ctx, mysql)
			if err != nil {
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
			}
			if !ready {
				// Give the resource some time to reach readiness
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelay}, nil
			}
			if result, err = r.onClusterReconcileMysqlStatus(ctx, mysql, listOpts); err != nil {
				return result, err
			}
			// If the object is being deleted then delete all sub resources
			if mysql.DeletionTimestamp != nil {
				r.Log.Info("Detected deletion timestamp, starting cleanup", "mysql.Name", mysql.Name)
				if result, err := r.OnClusterCleanup(ctx, mysql); err != nil {
					return result, err
				}
			}
		case constants.MysqlDeploymentTypeAzure:
			if result, err = r.azureReconcileMysql(ctx, mysql); err != nil {
				return result, err
			}
			if result, err = r.azureReconcileStatus(ctx, mysql); err != nil {
				return result, err
			}
			// If the object is being deleted then delete all sub resources
			if mysql.DeletionTimestamp != nil {
				r.Log.Info("Detected deletion timestamp, starting cleanup", "mysql.Name", mysql.Name)
				if result, err = r.azureCleanup(ctx, mysql); err != nil {
					return result, err
				}
			}
		default:
			r.Log.Error(fmt.Errorf("%v", "Unrecognized deployment type"), "ensure correct spelling or supported type",
				"DeploymentType", mysql.Spec.Deployment.DeploymentType)

		}
		return ctrl.Result{Requeue: true, RequeueAfter: constants.ReconcilerRequeueDelay}, nil
	}
	r.Log.Error(
		fmt.Errorf("%v", "Spec.Deployment.DeploymentType empty"),
		"no deployment type selected",
	)

	return ctrl.Result{Requeue: true, RequeueAfter: constants.ReconcilerRequeueDelay}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DBMMOMySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.DBMMOMySQL{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
