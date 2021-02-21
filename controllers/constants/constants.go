package constants

import "time"

const (
	// Operator/Project Constants ---------------------------------------------

	//CloudName is the name for Azure cloud
	CloudName string = "AzurePublicCloud"

	//ReconcilerRequeueDelayOnFail is the time delay for controllers between failed reconcile loops
	ReconcilerRequeueDelayOnFail = 5 * time.Second

	// ReconcilerRequeueDelay is the time delay for controllers between reconcile loops
	ReconcilerRequeueDelay = 20 * time.Second

	// Prefix is the prefix for all generated resources
	Prefix = "dbmmo"
	// OperatorName is the name of the operator
	OperatorName = "dbmmo-operator"

	// MysqlControllerName is the name of the mysql controller
	MysqlControllerName = "mysql-controller"

	// MYSQL Constants ---------------------------------------------

	//MysqlAzureClientIDEnvName is the constant envar name for AZURE_CLIENT_ID
	MysqlAzureClientIDEnvName = "AZURE_CLIENT_ID"
	//MysqlAzureClientSecretEnvName is the constant envar name for AZURE_CLIENT_SECRET
	MysqlAzureClientSecretEnvName = "AZURE_CLIENT_SECRET"
	//MysqlAzureTenantIDEnvName is the constant envar name for AZURE_TENANT_ID
	MysqlAzureTenantIDEnvName = "AZURE_TENANT_ID"
	//MysqlAzureSubscriptionIDEnvName is the constant envar name for AZURE_SUBSCRIPTION_ID
	MysqlAzureSubscriptionIDEnvName = "AZURE_SUBSCRIPTION_ID"
	//MysqlAzureBaseGroupNameEnvName is the constant envar name for AZURE_BASE_GROUP_NAME
	MysqlAzureBaseGroupNameEnvName = "AZURE_BASE_GROUP_NAME"
	//MysqlAzureLocationDefaultEnvName is the constant envar name for AZURE_LOCATION_DEFAULT
	MysqlAzureLocationDefaultEnvName = "AZURE_LOCATION_DEFAULT"
	// MysqlDeploymentTypeOnCluster Is the deployment type used to indicate that the database should be deployed on the same cluster as the operator
	MysqlDeploymentTypeOnCluster = "OnCluster"
	//MysqlDeploymentTypeAzure Is the deployment type used to indicate that the database should be deployed on azure
	MysqlDeploymentTypeAzure = "Azure"
	//MysqlName is the default name for mysql
	MysqlName = "mysql"
	//MysqlServiceName is the default name of the mysql service
	MysqlServiceName = "mysql-service"
	//MysqlServicePort is the default port on which the service will run
	MysqlServicePort = 3306
	//MysqlAppSelector is the label mapped to the `App:` label
	MysqlAppSelector = "dbmmo-mysql"
	//MysqlDeploymentName is the default name of the mysql deployment
	MysqlDeploymentName = "mysql-deployment"
	//MysqlStrategyType is the default deployment strategy
	MysqlStrategyType = "Recreate"
	//MysqlContainerImage is the default container image from which the container will be created
	MysqlContainerImage = "mysql:5.6"
	//MysqlContainerName is the default container name for the container
	MysqlContainerName = "mysql-container"
	//MysqlSecretEnvName is the default env var from which the mysql password will be retrieved
	MysqlSecretEnvName = "MYSQL_ROOT_PASSWORD"
	//MysqlIngressName is the default name for the mysql ingress
	MysqlIngressName = "dbmmo-mysql-ingress"
	//MysqlSecretEnvVal is the default password with which mysql will be set up
	MysqlSecretEnvVal = "password"
	//MysqlAdminUser is the default admin password
	MysqlAdminUser = "dbmmo_admin"
	//MysqlPathTypePrefix is the default pathtype for the ingress
	MysqlPathTypePrefix = "Prefix"
	//MysqlTargetPort is the default target port for mysql ingress
	MysqlTargetPort = 3306
	//MysqlDefaultPath is teh default path for the mysql ingres
	MysqlDefaultPath = "/"
	//MysqlHostName is the name of the host
	MysqlHostName = "mysql-host-name"
	//MysqlContainerPort is the default port from which the container will run the app
	MysqlContainerPort = 3306
	//MysqlContainerPortName is the default port name for the application port
	MysqlContainerPortName = "mysql"
	//MysqlVolumeMountName is the default name for the volume mount
	MysqlVolumeMountName = "mysql-persistent-storage"
	//MysqlVolumeMountPath is the default path to the volume mount
	MysqlVolumeMountPath = "/var/lib/mysql"
	//MysqlClaimName is the default volume claim name
	MysqlClaimName = "mysql-pv-claim"

	// MYSQL PV ---------------------------------------------

	// MysqlPVName is the default name of the private volume
	MysqlPVName = "mysql-pv-volume"
	// MysqlPVLabelType is the default deployment label
	MysqlPVLabelType = "local"
	// MysqlStorageClassName is the default storage class name
	MysqlStorageClassName = "manual"
	// MysqlCapacityStorage is the default storage capacity for the private volume
	MysqlCapacityStorage = "20Gi"
	// MysqlCapacityStorageTest is the default storage capacity for the private volume, used for testing only
	MysqlCapacityStorageTest = "2Gi"
	// MysqlPVAccessModes is the default access mode for the private volume
	MysqlPVAccessModes = "ReadWriteOnce"
	// MysqlPVHostPath is the default path for the private volume host
	MysqlPVHostPath = "/mnt/data"
	// MysqlResourceRequestStorage is the default resource requested from storage
	MysqlResourceRequestStorage = "20Gi"
	// MysqlResourceRequestStorageTestSize is the default testing size
	MysqlResourceRequestStorageTestSize = "2Gi"

	// Future DBs VV ---------------------------------------------
)
