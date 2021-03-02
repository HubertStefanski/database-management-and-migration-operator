package mysql

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/mysql/mgmt/mysqlflexibleservers"

	"github.com/Azure/go-autorest/autorest/to"
	cachev1alpha1 "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/util"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *DBMMOMySQLReconciler) azureReconcileMysql(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL) (ctrl.Result, error) {
	if util.ValidateAzureConfig(mysql.Spec.Deployment) {
		r.Log.Info("Reconciling MySQL on Azure", "Mysql.ServerName", mysql.Spec.Deployment.ServerName)
		//If Azure state doesn't indicate an error and hasn't been created, then create it
		//exists, err := util.ServerExists(ctx, mysql) //TODO revisit me
		//if err != nil {
		//	r.Log.Error(err, "Could not retrieve server/s", "mysql.ServerName", *mysql.Spec.Deployment.ServerName)
		//	mysql.Status.AzureStatus.State = cachev1alpha1.AzureError
		//	_, _ = r.azureReconcileStatus(ctx, mysql, nil)
		//	return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
		//}

		if !mysql.Status.AzureStatus.Created {
			r.Log.Info("Mysql Azure instance creating, please wait", "mysql.ServerName", *mysql.Spec.Deployment.ServerName, "Config:", *mysql.Spec.Deployment.AzureConfig)
			server, err := util.CreateServer(ctx, mysql)
			r.Log.Info("server", "server", server)
			if err != nil {
				if result, err := r.azureReconcileStatus(ctx, mysql, nil); err != nil {
					return result, err
				}
				return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
			}
			r.Log.Info("Created MySQL on Azure ", "Mysql.ServerName", mysql.Spec.Deployment.ServerName)

			mysql.Status.AzureStatus.State = cachev1alpha1.AzureCreated
			mysql.Status.AzureStatus.Created = true

			//update the status to prevent next creation loop
			if result, err := r.azureReconcileStatus(ctx, mysql, &server); err != nil {
				return result, err
			}
		} else {
			r.Log.Info("MySQL on Azure already exists, nothing to do", "Mysql.ServerName", *mysql.Spec.Deployment.ServerName)
		}
		// update the status to prevent next creation loop
		if result, err := r.azureReconcileStatus(ctx, mysql, nil); err != nil {
			return result, err
		}
		r.Log.Info("Executing query on Azure MySQL", "Mysql.ServerName", *mysql.Spec.Deployment.ServerName)
		err := util.ConnectAndExec(*mysql.Spec.Deployment.TableStatement,
			*mysql.Status.AzureStatus.ServerInfo.AdministratorLogin,
			*mysql.Status.AzureStatus.ServerInfo.AdministratorLogin,
			*mysql.Status.AzureStatus.ServerInfo.FullyQualifiedDomainName,
			*mysql.Spec.Deployment.ServerName)
		if err != nil {
			return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
		}
	} else {
		r.Log.Error(fmt.Errorf("%v", "Spec.Deployment Azure field misconfiguration"), "ensure data is valid",
			"Deployment", mysql.Spec.Deployment)
		return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, nil
	}
	r.Log.Info("Reconciled MySQL on Azure ", "Mysql.ServerName", mysql.Spec.Deployment.ServerName)
	return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelay}, nil
}

func (r *DBMMOMySQLReconciler) azureReconcileStatus(ctx context.Context, mysql *cachev1alpha1.DBMMOMySQL, server *mysqlflexibleservers.Server) (ctrl.Result, error) {
	r.Log.Info("Reconciling Mysql Status", "Mysql.Namespace", mysql.Namespace, "Mysql.Name", mysql.Name)

	if server != nil {
		mysql.Status.AzureStatus.ServerInfo = cachev1alpha1.ServerInfo{
			Tags:                       server.Tags,
			Location:                   server.Location,
			ID:                         server.ID,
			Name:                       server.Name,
			Type:                       server.Type,
			AdministratorLogin:         server.ServerProperties.AdministratorLogin,
			AdministratorLoginPassword: server.ServerProperties.AdministratorLoginPassword,
			//State:                      server.ServerProperties.State,
			FullyQualifiedDomainName: server.FullyQualifiedDomainName,
			ReplicationRole:          server.ReplicationRole,
			ReplicaCapacity:          server.ReplicaCapacity,
			SourceServerID:           server.SourceServerID,
			AvailabilityZone:         server.AvailabilityZone,
		}
	} else {
		srv, err := util.GetServer(ctx, mysql)
		if err != nil {
			return ctrl.Result{RequeueAfter: constants.ReconcilerRequeueDelayOnFail}, err
		}

		mysql.Status.AzureStatus.ServerInfo = cachev1alpha1.ServerInfo{
			Tags:                       srv.Tags,
			Location:                   srv.Location,
			ID:                         srv.ID,
			Name:                       srv.Name,
			Type:                       srv.Type,
			AdministratorLogin:         to.StringPtr(constants.MysqlAdminUser),
			AdministratorLoginPassword: to.StringPtr(constants.MysqlSecretEnvVal),
			//State:                      to.StringPtr(server.ServerProperties.State),
			FullyQualifiedDomainName: srv.FullyQualifiedDomainName,
			//ReplicationRole:          srv.ReplicationRole,
			ReplicaCapacity:  srv.ReplicaCapacity,
			SourceServerID:   srv.SourceServerID,
			AvailabilityZone: srv.AvailabilityZone,
		}
	}
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
