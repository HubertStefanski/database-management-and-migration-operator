package mysql

import (
	"context"
	"fmt"
	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/util"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *DBMMOMySQLReconciler) azureReconcileMysql(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	if util.ValidateAzureConfig(mysql.Spec.Deployment) {
		r.Log.Info("Reconciling MySQL on Azure", "Mysql.ServerName", mysql.Spec.Deployment.ServerName)
		// If Azure state doesn't indicate an error and hasn't been created, then create it
		if !mysql.Status.AzureStatus.Created {
			r.Log.Info("Mysql Azure instance creating, please wait", "mysql.ServerName", *mysql.Spec.Deployment.ServerName, "Config:", *mysql.Spec.Deployment.AzureConfig)
			server, err := util.CreateServer(ctx, mysql)
			if err != nil {
				mysql.Status.AzureStatus.State = cachev1alpha1.AzureError
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
			}
			//Update the status for future reference to the server
			mysql.Status.AzureStatus.ServerInfo = cachev1alpha1.ServerInfo{
				Tags:     server.Tags,
				Location: server.Location,
				ID:       server.ID,
				Name:     server.Name,
				Type:     server.Type,
			}
			mysql.Status.AzureStatus.State = cachev1alpha1.AzureCreated
			mysql.Status.AzureStatus.Created = true

			// update the status to prevent next creation loop
			if result, err := r.azureReconcileStatus(ctx, mysql); err != nil {
				return result, err
			}
			return ctrl.Result{}, nil
		}

	} else {
		r.Log.Error(fmt.Errorf("%v", "Spec.Deployment Azure field misconfiguration"), "ensure data is valid",
			"Deployment", mysql.Spec.Deployment)
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, nil
	}
	r.Log.Info("Reconciled MySQL on Azure ", "Mysql.ServerName", mysql.Spec.Deployment.ServerName)
	return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelay}, nil
}

func (r *DBMMOMySQLReconciler) azureReconcileStatus(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	r.Log.Info("Reconciling Mysql Status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)

	if err := r.Client.Status().Update(ctx, mysql); err != nil {
		r.Log.Error(err, "Failed to update Mysql status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
		return ctrl.Result{}, err
	}

	r.Log.Info("Mysql status reconciled", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
	return ctrl.Result{Requeue: true}, nil
}
