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
		exists, err := util.ServerExists(ctx, mysql)
		if err != nil {
			r.Log.Error(err, "Could not retrieve server/s", "mysql.ServerName", *mysql.Spec.Deployment.ServerName)
			mysql.Status.AzureStatus.State = cachev1alpha1.AzureError
			_, _ = r.azureReconcileStatus(ctx, mysql)
			return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
		}

		if !mysql.Status.AzureStatus.Created || mysql.Status.AzureStatus.State == cachev1alpha1.AzureError || !exists {
			r.Log.Info("Mysql Azure instance creating, please wait", "mysql.ServerName", *mysql.Spec.Deployment.ServerName, "Config:", *mysql.Spec.Deployment.AzureConfig)
			server, err := util.CreateServer(ctx, mysql)
			if err != nil {
				mysql.Status.AzureStatus.State = cachev1alpha1.AzureError
				if result, err := r.azureReconcileStatus(ctx, mysql); err != nil {
					return result, err
				}
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
			}
			mysql.Status.AzureStatus.State = cachev1alpha1.AzureCreated
			mysql.Status.AzureStatus.Created = true

			mysql.Status.AzureStatus.ServerInfo = cachev1alpha1.ServerInfo{
				Tags:                       server.Tags,
				Location:                   server.Location,
				ID:                         server.ID,
				Name:                       server.Name,
				Type:                       server.Type,
				AdministratorLogin:         server.ServerProperties.AdministratorLogin,
				AdministratorLoginPassword: server.ServerProperties.AdministratorLoginPassword,
				State:                      server.ServerProperties.State,
				FullyQualifiedDomainName:   server.FullyQualifiedDomainName,
				ReplicationRole:            server.ReplicationRole,
				ReplicaCapacity:            server.ReplicaCapacity,
				SourceServerID:             server.SourceServerID,
				AvailabilityZone:           server.AvailabilityZone,
			}
			err = util.ConnectAndExec(*mysql.Spec.Deployment.TableStatement,
				*mysql.Status.AzureStatus.ServerInfo.AdministratorLogin,
				*mysql.Status.AzureStatus.ServerInfo.AdministratorLogin,
				*mysql.Status.AzureStatus.ServerInfo.FullyQualifiedDomainName,
				*mysql.Spec.Deployment.ServerName)
			if err != nil {
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
			}
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

	r.Log.Info("Reconciled Mysql status ", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)
	return ctrl.Result{Requeue: true}, nil
}

func (r *DBMMOMySQLReconciler) azureCleanup(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	mysql.Status.AzureStatus.State = cachev1alpha1.AzureDeleting
	r.Log.Info("Deleting MySQL on Azure", "mysql.ServerName", mysql.Spec.Deployment.ServerName)
	if resp, err := util.DeleteServer(ctx, *mysql.Status.AzureStatus.ServerInfo.Name, mysql); err != nil || resp.StatusCode != 200 {
		r.Log.Error(err, "Failed to delete mysql Azure server", "mysql.ServerName", mysql.Spec.Deployment.ServerName)
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
	}
	r.Log.Info("Deleted MySQL on Azure", "mysql.ServerName", mysql.Spec.Deployment.ServerName)
	return ctrl.Result{}, nil

}
