package util

import "github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"

// ValidateAzureConfig verifies whether the required fields have been configured correctly
func ValidateAzureConfig(dep *v1alpha1.DBMMOMYSQLDeployment) bool {
	if dep.AzureConfig != nil {
		if dep.AzureConfig.ClientID == nil || *dep.AzureConfig.ClientID != "" {
			return false
		}
		if dep.AzureConfig.ClientSecret == nil || *dep.AzureConfig.ClientSecret != "" {
			return false
		}
		if dep.AzureConfig.TenantID == nil || *dep.AzureConfig.TenantID != "" {
			return false
		}
		if dep.AzureConfig.SubscriptionID == nil || *dep.AzureConfig.BaseGroupName != "" {
			return false
		}
		if dep.AzureConfig.BaseGroupName == nil || *dep.AzureConfig.BaseGroupName != "" {
			return false
		}
		if dep.AzureConfig.LocationDefault == nil || *dep.AzureConfig.LocationDefault != "" {
			return false
		}
		if dep.AzureConfig.OAuthGrantType == nil {
			return false
		}

		return true
	}
	return false
}
