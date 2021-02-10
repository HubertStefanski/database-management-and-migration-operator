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

	if mysql.Spec.DeploymentType != "" {
		switch depType := mysql.Spec.DeploymentType; depType {
		case constants.MysqlDeploymentTypeOnCluster:
			if result, err := r.onClusterReconcileMysqlPVC(ctx, mysql); err != nil {
				return result, err
			}

			if result, err := r.onClusterReconcileMysqlService(ctx, mysql); err != nil {
				return result, err
			}

			if result, err := r.onClusterReconcileMysqlDeployment(ctx, mysql); err != nil {
				return result, err
			}

			if result, err := r.onClusterReconcileMysqlStatus(ctx, mysql, listOpts); err != nil {
				return result, err
			}
			// If the object is being deleted then delete all sub resources
			if mysql.DeletionTimestamp != nil {
				if result, err := r.OnClusterCleanup(ctx, mysql); err != nil {
					return result, err
				}
			}
		case constants.MysqlDeploymentTypeAzure:
			r.Log.Info("Specified deployment type is not currently supported, please enter a supported type",
				"DeploymentType", mysql.Spec.DeploymentType)
		default:
			r.Log.Error(fmt.Errorf("%v", "Unrecognized deployment type"), "ensure correct spelling or supported type",
				"DeploymentType", mysql.Spec.DeploymentType)

		}
		return ctrl.Result{Requeue: true, RequeueAfter: constants.ReconcilerRequeueDelay}, nil
	}
	r.Log.Error(
		fmt.Errorf("%v", "Spec.DeploymentType empty"),
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
