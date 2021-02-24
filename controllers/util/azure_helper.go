package util

import (
	"context"
	"database/sql"
	"fmt"
	mysql "github.com/Azure/azure-sdk-for-go/services/preview/mysql/mgmt/2020-07-01-preview/mysqlflexibleservers"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	"github.com/HubertStefanski/database-management-and-migration-operator/controllers/constants"
	_ "github.com/go-sql-driver/mysql"
)

var (
	armAuthorizer autorest.Authorizer
	environment   *azure.Environment
)

const (
	// OAuthGrantTypeServicePrincipal for client credentials flow
	OAuthGrantTypeServicePrincipal v1alpha1.OAuthGrantType = iota
	// OAuthGrantTypeDeviceFlow for device flow
	OAuthGrantTypeDeviceFlow
)

// GetServersClient returns
func getServersClient(m *v1alpha1.DBMMOMySQL) mysql.ServersClient {
	serversClient := mysql.NewServersClient(*m.Spec.Deployment.AzureConfig.SubscriptionID)
	a, _ := GetResourceManagementAuthorizer(m)
	serversClient.Authorizer = a
	//_ = serversClient.AddToUserAgent(*m.Spec.Deployment.AzureConfig.UserAgent)
	return serversClient
}

func ServerExists(ctx context.Context, m *v1alpha1.DBMMOMySQL) (bool, error) {
	server := getServersClient(m)
	if res, err := server.Get(ctx, *m.Spec.Deployment.AzureConfig.BaseGroupName, *m.Spec.Deployment.ServerName); res.Response.StatusCode == 404 {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil

}

func GetServer(ctx context.Context, m *v1alpha1.DBMMOMySQL) (mysql.Server, error) {
	server := getServersClient(m)
	if res, err := server.Get(ctx, *m.Spec.Deployment.AzureConfig.BaseGroupName, *m.Spec.Deployment.ServerName); res.Response.StatusCode == 404 {
		return res, nil
	} else if err != nil {
		return res, err
	}
	return mysql.Server{}, nil

}

// CreateServer creates a new MySQL Server
func CreateServer(ctx context.Context, m *v1alpha1.DBMMOMySQL) (server mysql.Server, err error) {
	serversClient := getServersClient(m)

	// Create the server
	future, err := serversClient.Create(
		ctx,
		*m.Spec.Deployment.AzureConfig.BaseGroupName,
		*m.Spec.Deployment.ServerName,
		mysql.Server{
			Location: to.StringPtr(*m.Spec.Deployment.AzureConfig.LocationDefault),
			Sku: &mysql.Sku{
				Name: to.StringPtr("Standard_D16ds_v4"),
				Tier: "GeneralPurpose",
			},
			ServerProperties: &mysql.ServerProperties{
				AdministratorLogin:         to.StringPtr(constants.MysqlAdminUser),    //TODO replace me with cr field val
				AdministratorLoginPassword: to.StringPtr(constants.MysqlSecretEnvVal), //TODO replace me with cr field val
				Version:                    mysql.FiveFullStopSeven,                   // 5.7
				StorageProfile: &mysql.StorageProfile{
					StorageMB: to.Int32Ptr(524288),
				},
			},
		})

	if err != nil {
		return server, fmt.Errorf("cannot create mysql server: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, serversClient.Client)
	if err != nil {
		return server, fmt.Errorf("cannot get the mysql server create or update future response: %v", err)
	}

	return future.Result(serversClient)
}

// UpdateServerStorageCapacity given the server name and the new storage capacity it updates the server's storage capacity.
func UpdateServerStorageCapacity(ctx context.Context, serverName string, storageCapacity int32, m *v1alpha1.DBMMOMySQL) (server mysql.Server, err error) {
	serversClient := getServersClient(m)

	future, err := serversClient.Update(
		ctx,
		*m.Spec.Deployment.AzureConfig.BaseGroupName,
		serverName,
		mysql.ServerForUpdate{
			ServerPropertiesForUpdate: &mysql.ServerPropertiesForUpdate{
				StorageProfile: &mysql.StorageProfile{
					StorageMB: &storageCapacity,
				},
			},
		},
	)
	if err != nil {
		return server, fmt.Errorf("cannot update mysql server: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, serversClient.Client)
	if err != nil {
		return server, fmt.Errorf("cannot get the mysql server update future response: %v", err)
	}

	return future.Result(serversClient)
}

// DeleteServer deletes the MySQL server.
func DeleteServer(ctx context.Context, serverName string, m *v1alpha1.DBMMOMySQL) (resp autorest.Response, err error) {
	serversClient := getServersClient(m)

	future, err := serversClient.Delete(ctx, *m.Spec.Deployment.AzureConfig.BaseGroupName, serverName)
	if err != nil {
		return resp, fmt.Errorf("cannot delete the mysql server: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, serversClient.Client)
	if err != nil {
		return resp, fmt.Errorf("cannot get the mysql server update future response: %v", err)
	}

	return future.Result(serversClient)
}

// GetFwRulesClient returns the FirewallClient
func getFwRulesClient(m *v1alpha1.DBMMOMySQL) mysql.FirewallRulesClient {
	fwrClient := mysql.NewFirewallRulesClient(*m.Spec.Deployment.AzureConfig.SubscriptionID)
	a, _ := GetResourceManagementAuthorizer(m)
	fwrClient.Authorizer = a
	_ = fwrClient.AddToUserAgent(*m.Spec.Deployment.AzureConfig.UserAgent)
	return fwrClient
}

// CreateOrUpdateFirewallRule given the firewallname and new properties it updates the firewall rule.
func CreateOrUpdateFirewallRule(ctx context.Context, m *v1alpha1.DBMMOMySQL) error {
	fwrClient := getFwRulesClient(m)

	_, err := fwrClient.CreateOrUpdate(
		ctx,
		*m.Spec.Deployment.AzureConfig.BaseGroupName,
		*m.Spec.Deployment.ServerName,
		*m.Spec.Deployment.AzureConfig.AzureFwRule.FirewallRuleName,
		mysql.FirewallRule{
			FirewallRuleProperties: &mysql.FirewallRuleProperties{
				StartIPAddress: m.Spec.Deployment.AzureConfig.AzureFwRule.StartIPAddr,
				EndIPAddress:   m.Spec.Deployment.AzureConfig.AzureFwRule.EndIPAddr,
			},
		},
	)

	return err
}

// GetConfigurationsClient creates and returns the configuration client for the server.
func getConfigurationsClient(m *v1alpha1.DBMMOMySQL) mysql.ConfigurationsClient {
	configClient := mysql.NewConfigurationsClient(*m.Spec.Deployment.AzureConfig.SubscriptionID)
	a, _ := GetResourceManagementAuthorizer(m)
	configClient.Authorizer = a
	_ = configClient.AddToUserAgent(*m.Spec.Deployment.AzureConfig.UserAgent)
	return configClient
}

// GetConfiguration given the server name and configuration name it returns the configuration.
func GetConfiguration(ctx context.Context, serverName, configurationName string, m *v1alpha1.DBMMOMySQL) (mysql.Configuration, error) {
	configClient := getConfigurationsClient(m)

	// Get the configuration.
	configuration, err := configClient.Get(ctx, *m.Spec.Deployment.AzureConfig.BaseGroupName, serverName, configurationName)

	if err != nil {
		return configuration, fmt.Errorf("cannot get the configuration with name %s", configurationName)
	}

	return configuration, err
}

// UpdateConfiguration given the name of the configuation and the configuration object it updates the configuration for the given server.
func UpdateConfiguration(ctx context.Context, configurationName string, configuration mysql.Configuration, m *v1alpha1.DBMMOMySQL) (updatedConfig mysql.Configuration, err error) {
	configClient := getConfigurationsClient(m)

	future, err := configClient.Update(ctx, *m.Spec.Deployment.AzureConfig.BaseGroupName, *m.Spec.Deployment.ServerName, *m.Spec.Deployment.ConfigurationName, configuration)

	if err != nil {
		return updatedConfig, fmt.Errorf("cannot update the configuration with name %s", configurationName)
	}

	err = future.WaitForCompletionRef(ctx, configClient.Client)
	if err != nil {
		return updatedConfig, fmt.Errorf("cannot get the mysql configuration update future response: %v", err)
	}

	return future.Result(configClient)
}

// GetResourceManagementAuthorizer gets an OAuthTokenAuthorizer for Azure Resource Manager
func GetResourceManagementAuthorizer(m *v1alpha1.DBMMOMySQL) (autorest.Authorizer, error) {
	if armAuthorizer != nil {
		return armAuthorizer, nil
	}

	var a autorest.Authorizer
	var err error

	a, err = getAuthorizerForResource(m, getEnvironment().ResourceManagerEndpoint)

	if err == nil {
		// cache
		armAuthorizer = a
	} else {
		// clear cache
		armAuthorizer = nil
	}
	return armAuthorizer, err
}

func getAuthorizerForResource(m *v1alpha1.DBMMOMySQL, resource string) (autorest.Authorizer, error) {
	var a autorest.Authorizer
	var err error

	switch *m.Spec.Deployment.AzureConfig.OAuthGrantType {

	case OAuthGrantTypeServicePrincipal:
		oauthConfig, err := adal.NewOAuthConfig(
			getEnvironment().ActiveDirectoryEndpoint, *m.Spec.Deployment.AzureConfig.TenantID)
		if err != nil {
			return nil, err
		}

		token, err := adal.NewServicePrincipalToken(
			*oauthConfig, *m.Spec.Deployment.AzureConfig.ClientID, *m.Spec.Deployment.AzureConfig.ClientSecret, resource)
		if err != nil {
			return nil, err
		}
		a = autorest.NewBearerAuthorizer(token)

	case OAuthGrantTypeDeviceFlow:
		deviceConfig := auth.NewDeviceFlowConfig(*m.Spec.Deployment.AzureConfig.ClientID, *m.Spec.Deployment.AzureConfig.TenantID)
		deviceConfig.Resource = resource
		a, err = deviceConfig.Authorizer()
		if err != nil {
			return nil, err
		}

	default:
		return a, fmt.Errorf("invalid grant type specified")
	}

	return a, err
}

func getEnvironment() *azure.Environment {
	if environment != nil {
		return environment
	}
	env, err := azure.EnvironmentFromName(constants.CloudName)
	if err != nil {
		// TODO: move to initialization of var
		panic(fmt.Sprintf(
			"invalid cloud name '%s' specified, cannot continue\n", constants.CloudName))
	}
	environment = &env
	return environment
}

func ConnectAndExec(query, user, password, host, database string) error {
	// Initialize connection string.
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)

	// Initialize connection object.
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
