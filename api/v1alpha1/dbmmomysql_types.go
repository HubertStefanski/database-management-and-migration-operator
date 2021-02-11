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

// DBMMOMYSQLDeployment defines the desired state of the mysqlDeployment
type DBMMOMYSQLDeployment struct {
	DeploymentType *string            `json:"deploymentType,omitempty"`
	EnvFrom        []v1.EnvFromSource `json:"envFrom,omitempty"`
}

// DBMMOMySQLStatus defines the observed state of DBMMOMySQL
type DBMMOMySQLStatus struct {
	Nodes                  []string `json:"nodes,omitempty"`
	Services               []string `json:"services,omitempty"`
	PersistentVolumeClaims []string `json:"persistentVolumeClaims,omitempty"`
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
