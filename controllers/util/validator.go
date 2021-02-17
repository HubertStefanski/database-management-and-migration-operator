package util

import "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"

// ValidateAzureConfig verifies whether the required fields have been configured correctly
func ValidateAzureConfig(dep *v1alpha1.DBMMOMYSQLDeployment) bool {
	if dep.AzureConfig != nil {
		if dep.AzureConfig.AzureClientID == nil || *dep.AzureConfig.AzureClientID != "" {
			return false
		}
		if dep.AzureConfig.AzureClientSecret == nil || *dep.AzureConfig.AzureClientSecret != "" {
			return false
		}
		if dep.AzureConfig.AzureTenantID == nil || *dep.AzureConfig.AzureTenantID != "" {
			return false
		}
		if dep.AzureConfig.AzureSubscriptionID == nil || *dep.AzureConfig.AzureBaseGroupName != "" {
			return false
		}
		if dep.AzureConfig.AzureBaseGroupName == nil || *dep.AzureConfig.AzureBaseGroupName != "" {
			return false
		}
		if dep.AzureConfig.AzureLocationDefault == nil || *dep.AzureConfig.AzureLocationDefault != "" {
			return false
		}
		return true
	}
	return false
}
