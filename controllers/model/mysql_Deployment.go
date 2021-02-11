package model

import (
	"github.com/HubertStefanski/database-management-and-migration-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// MysqlDeploymentGetEnvFrom returns the environment variables contained within a secret
func MysqlDeploymentGetEnvFrom(m *v1alpha1.DBMMOMySQL) []v1.EnvFromSource {
	var envFrom []v1.EnvFromSource
	if m.Spec.Deployment != nil && m.Spec.Deployment.EnvFrom != nil {
		for _, v := range m.Spec.Deployment.EnvFrom {
			envFrom = append(envFrom, *v.DeepCopy())
		}
	}
	return envFrom
}
