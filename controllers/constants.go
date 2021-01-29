package controllers

const (
	// Operator/Project Constants ---------------------------------------------

	// OperatorName is the name of the operator
	OperatorName = "dbmmo-operator"

	// MYSQL Constants ---------------------------------------------

	//MysqlServiceName is the default name of the mysql service
	MysqlServiceName = "mysql-service"
	//MysqlServicePort is the default port on which the service will run
	MysqlServicePort = 3306
	//MysqlAppSelector is the label mapped to the `App:` label
	MysqlAppSelector = "mysql"
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

	// TODO: refactor this to use a secret as opposed to envars

	//MysqlSecretEnvVal is the default password with which mysql will be set up
	MysqlSecretEnvVal = "password"
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
	// MysqlPVAccessModes is the default access mode for the private volume
	MysqlPVAccessModes = "ReadWriteOnce"
	// MysqlPVHostPath is the default path for the private volume host
	MysqlPVHostPath = "/mnt/data"
	// MysqlResourceRequestStorage is the default resource requested from storage
	MysqlResourceRequestStorage = "20Gi"

	// Future DBs VV ---------------------------------------------
)
