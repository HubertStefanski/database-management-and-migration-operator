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

package v1alpha1

import (
	mysql "github.com/Azure/azure-sdk-for-go/services/preview/mysql/mgmt/2020-07-01-preview/mysqlflexibleservers"
	"github.com/Azure/go-autorest/autorest/azure"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
// Important: Run "make" to regenerate code after modifying this file

// DBMMOMySQLSpec defines the desired state of DBMMOMySQL
type DBMMOMySQLSpec struct {
	// Size indicates the number of pods to be deployed
	Size int32 `json:"size,omitempty"`
	// Deployment defines the desired state of the deployment for this resource
	Deployment *DBMMOMYSQLDeployment `json:"deployment,omitempty"`
}

// OAuthGrantType defines the desired type of OAuthGrant
type OAuthGrantType int

// AzureConfig defines all required fields for Azure
type AzureConfig struct {
	ClientID               *string            `json:"azureClientID,omitempty"`
	ClientSecret           *string            `json:"azureClientSecret,omitempty"`
	TenantID               *string            `json:"azureTenantID,omitempty"`
	SubscriptionID         *string            `json:"azureSubscriptionID,omitempty"`
	BaseGroupName          *string            `json:"azureBaseGroupName,omitempty"`
	LocationDefault        *string            `json:"azureLocationDefault,omitempty"`
	ConfigurationName      *string            `json:"azureConfigurationName,omitempty"`
	OAuthGrantType         *OAuthGrantType    `json:"oauthGrantType,omitempty"`
	AuthorizationServerURL *string            `json:"authorizationServerURL,omitempty"`
	CloudName              *string            `json:"cloudName,omitempty"` //"AzurePublicCloud"
	UseDeviceFlow          *bool              `json:"useDeviceFlow,omitempty"`
	KeepResources          *bool              `json:"keepResources,omitempty"`
	UserAgent              *string            `json:"userAgent,omitempty"`
	Environment            *azure.Environment `json:"environment,omitempty"`
	AzureFwRule            *AzureFwRule       `json:"azureFwRule,omitempty"`
}

//AzureFwRule defines desired state of the Azure firewall rule
type AzureFwRule struct {
	FirewallRuleName *string `json:"firewallRuleName,omitempty"`
	StartIPAddr      *string `json:"startIPAddr,omitempty"`
	EndIPAddr        *string `json:"endIPAddr,omitempty"`
}

// DBMMOMYSQLDeployment defines the desired state of the mysqlDeployment
type DBMMOMYSQLDeployment struct {
	ServerName        *string            `json:"serverName,omitempty"`
	ConfigurationName *string            `json:"configurationName,omitempty"`
	StorageCapacity   *int32             `json:"storageCapacity,omitempty"`
	DeploymentType    *string            `json:"deploymentType,omitempty"`
	EnvFrom           []v1.EnvFromSource `json:"envFrom,omitempty"`
	AzureConfig       *AzureConfig       `json:"azureConfig,omitempty"`
}

//AzureState indicates the state of the Azure server in one line
type AzureState string

const (
	//AzureCreated indicates that the server has been created
	AzureCreated AzureState = "AzureCreated"
	//AzureNotCreated indicates that the server has not been created
	AzureNotCreated AzureState = "NotCreated"
	//AzureError indicates that there was an issue while trying to invoke a connection to Azure
	AzureError AzureState = "Error"
)

//AzureStatus Indicates the currents status of the Azure deployment, including Creation, State and the Created Server
type AzureStatus struct {
	Created bool         `json:"created"`
	State   AzureState   `json:"azureState"`
	Server  mysql.Server `json:"mysqlServer"`
}

// DBMMOMySQLStatus defines the observed state of DBMMOMySQL
type DBMMOMySQLStatus struct {
	Nodes                  []string    `json:"nodes,omitempty"`
	Services               []string    `json:"services,omitempty"`
	PersistentVolumeClaims []string    `json:"persistentVolumeClaims"`
	AzureStatus            AzureStatus `json:"azureStatus"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DBMMOMySQL is the Schema for the dbmmomysqls API
type DBMMOMySQL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DBMMOMySQLSpec   `json:"spec,omitempty"`
	Status DBMMOMySQLStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DBMMOMySQLList contains a list of DBMMOMySQL
type DBMMOMySQLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DBMMOMySQL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DBMMOMySQL{}, &DBMMOMySQLList{})
}
