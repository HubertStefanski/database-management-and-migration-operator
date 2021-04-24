// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureConfig) DeepCopyInto(out *AzureConfig) {
	*out = *in
	if in.ClientID != nil {
		in, out := &in.ClientID, &out.ClientID
		*out = new(string)
		**out = **in
	}
	if in.ClientSecret != nil {
		in, out := &in.ClientSecret, &out.ClientSecret
		*out = new(string)
		**out = **in
	}
	if in.TenantID != nil {
		in, out := &in.TenantID, &out.TenantID
		*out = new(string)
		**out = **in
	}
	if in.SubscriptionID != nil {
		in, out := &in.SubscriptionID, &out.SubscriptionID
		*out = new(string)
		**out = **in
	}
	if in.BaseGroupName != nil {
		in, out := &in.BaseGroupName, &out.BaseGroupName
		*out = new(string)
		**out = **in
	}
	if in.LocationDefault != nil {
		in, out := &in.LocationDefault, &out.LocationDefault
		*out = new(string)
		**out = **in
	}
	if in.ConfigurationName != nil {
		in, out := &in.ConfigurationName, &out.ConfigurationName
		*out = new(string)
		**out = **in
	}
	if in.OAuthGrantType != nil {
		in, out := &in.OAuthGrantType, &out.OAuthGrantType
		*out = new(OAuthGrantType)
		**out = **in
	}
	if in.AuthorizationServerURL != nil {
		in, out := &in.AuthorizationServerURL, &out.AuthorizationServerURL
		*out = new(string)
		**out = **in
	}
	if in.CloudName != nil {
		in, out := &in.CloudName, &out.CloudName
		*out = new(string)
		**out = **in
	}
	if in.UseDeviceFlow != nil {
		in, out := &in.UseDeviceFlow, &out.UseDeviceFlow
		*out = new(bool)
		**out = **in
	}
	if in.KeepResources != nil {
		in, out := &in.KeepResources, &out.KeepResources
		*out = new(bool)
		**out = **in
	}
	if in.UserAgent != nil {
		in, out := &in.UserAgent, &out.UserAgent
		*out = new(string)
		**out = **in
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = new(azure.Environment)
		**out = **in
	}
	if in.AzureFwRule != nil {
		in, out := &in.AzureFwRule, &out.AzureFwRule
		*out = new(AzureFwRule)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureConfig.
func (in *AzureConfig) DeepCopy() *AzureConfig {
	if in == nil {
		return nil
	}
	out := new(AzureConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureFwRule) DeepCopyInto(out *AzureFwRule) {
	*out = *in
	if in.FirewallRuleName != nil {
		in, out := &in.FirewallRuleName, &out.FirewallRuleName
		*out = new(string)
		**out = **in
	}
	if in.StartIPAddr != nil {
		in, out := &in.StartIPAddr, &out.StartIPAddr
		*out = new(string)
		**out = **in
	}
	if in.EndIPAddr != nil {
		in, out := &in.EndIPAddr, &out.EndIPAddr
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureFwRule.
func (in *AzureFwRule) DeepCopy() *AzureFwRule {
	if in == nil {
		return nil
	}
	out := new(AzureFwRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureStatus) DeepCopyInto(out *AzureStatus) {
	*out = *in
	in.ServerInfo.DeepCopyInto(&out.ServerInfo)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureStatus.
func (in *AzureStatus) DeepCopy() *AzureStatus {
	if in == nil {
		return nil
	}
	out := new(AzureStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOIngress) DeepCopyInto(out *DBMMOIngress) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOIngress.
func (in *DBMMOIngress) DeepCopy() *DBMMOIngress {
	if in == nil {
		return nil
	}
	out := new(DBMMOIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOMYSQLDeployment) DeepCopyInto(out *DBMMOMYSQLDeployment) {
	*out = *in
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(DBMMOIngress)
		(*in).DeepCopyInto(*out)
	}
	if in.ServerName != nil {
		in, out := &in.ServerName, &out.ServerName
		*out = new(string)
		**out = **in
	}
	if in.ConfigurationName != nil {
		in, out := &in.ConfigurationName, &out.ConfigurationName
		*out = new(string)
		**out = **in
	}
	if in.StorageCapacity != nil {
		in, out := &in.StorageCapacity, &out.StorageCapacity
		*out = new(string)
		**out = **in
	}
	if in.DeploymentType != nil {
		in, out := &in.DeploymentType, &out.DeploymentType
		*out = new(string)
		**out = **in
	}
	if in.ConfirmMigrate != nil {
		in, out := &in.ConfirmMigrate, &out.ConfirmMigrate
		*out = new(bool)
		**out = **in
	}
	if in.EnvFrom != nil {
		in, out := &in.EnvFrom, &out.EnvFrom
		*out = make([]v1.EnvFromSource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ServerCredentials != nil {
		in, out := &in.ServerCredentials, &out.ServerCredentials
		*out = new(MysqlCredentials)
		(*in).DeepCopyInto(*out)
	}
	if in.AzureConfig != nil {
		in, out := &in.AzureConfig, &out.AzureConfig
		*out = new(AzureConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOMYSQLDeployment.
func (in *DBMMOMYSQLDeployment) DeepCopy() *DBMMOMYSQLDeployment {
	if in == nil {
		return nil
	}
	out := new(DBMMOMYSQLDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOMySQL) DeepCopyInto(out *DBMMOMySQL) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOMySQL.
func (in *DBMMOMySQL) DeepCopy() *DBMMOMySQL {
	if in == nil {
		return nil
	}
	out := new(DBMMOMySQL)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBMMOMySQL) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOMySQLList) DeepCopyInto(out *DBMMOMySQLList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DBMMOMySQL, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOMySQLList.
func (in *DBMMOMySQLList) DeepCopy() *DBMMOMySQLList {
	if in == nil {
		return nil
	}
	out := new(DBMMOMySQLList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DBMMOMySQLList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOMySQLSpec) DeepCopyInto(out *DBMMOMySQLSpec) {
	*out = *in
	if in.Deployment != nil {
		in, out := &in.Deployment, &out.Deployment
		*out = new(DBMMOMYSQLDeployment)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOMySQLSpec.
func (in *DBMMOMySQLSpec) DeepCopy() *DBMMOMySQLSpec {
	if in == nil {
		return nil
	}
	out := new(DBMMOMySQLSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DBMMOMySQLStatus) DeepCopyInto(out *DBMMOMySQLStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Services != nil {
		in, out := &in.Services, &out.Services
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PersistentVolumeClaims != nil {
		in, out := &in.PersistentVolumeClaims, &out.PersistentVolumeClaims
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.AzureStatus.DeepCopyInto(&out.AzureStatus)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DBMMOMySQLStatus.
func (in *DBMMOMySQLStatus) DeepCopy() *DBMMOMySQLStatus {
	if in == nil {
		return nil
	}
	out := new(DBMMOMySQLStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MysqlCredentials) DeepCopyInto(out *MysqlCredentials) {
	*out = *in
	if in.MysqlAdministratorLogin != nil {
		in, out := &in.MysqlAdministratorLogin, &out.MysqlAdministratorLogin
		*out = new(string)
		**out = **in
	}
	if in.MysqlAdministratorLoginPassword != nil {
		in, out := &in.MysqlAdministratorLoginPassword, &out.MysqlAdministratorLoginPassword
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MysqlCredentials.
func (in *MysqlCredentials) DeepCopy() *MysqlCredentials {
	if in == nil {
		return nil
	}
	out := new(MysqlCredentials)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerInfo) DeepCopyInto(out *ServerInfo) {
	*out = *in
	if in.Tags != nil {
		in, out := &in.Tags, &out.Tags
		*out = make(map[string]*string, len(*in))
		for key, val := range *in {
			var outVal *string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = new(string)
				**out = **in
			}
			(*out)[key] = outVal
		}
	}
	if in.Location != nil {
		in, out := &in.Location, &out.Location
		*out = new(string)
		**out = **in
	}
	if in.ID != nil {
		in, out := &in.ID, &out.ID
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerInfo.
func (in *ServerInfo) DeepCopy() *ServerInfo {
	if in == nil {
		return nil
	}
	out := new(ServerInfo)
	in.DeepCopyInto(out)
	return out
}
